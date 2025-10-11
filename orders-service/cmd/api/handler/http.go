package handler

import (
	"context"
	"errors"
	"net/http"
	appError "orders-app/internal/error"
	"orders-app/internal/model"
	"rpc"
	"strings"
	"time"

	"github.com/google/uuid"
)

const clientTimeout = time.Second * 20

func (app *HttpHandler) Hello(w http.ResponseWriter, r *http.Request) {

	app.writeJson(w, http.StatusOK, model.OrderRequest{
		ItemId: "123",
		Qty:    2,
	})
}

func (app *HttpHandler) AddOrder(w http.ResponseWriter, r *http.Request) {

	var req model.OrderRequest
	err := app.readJSON(r, &req)
	if err != nil {
		processError(app, w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), clientTimeout)
	defer cancel()

	itemStock, err := app.InventoryService.CheckStock(ctx, &rpc.CheckStockRequest{
		ItemId: req.ItemId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "item not found") {
			app.errorJSON(w, appError.ErrItemNotFound, http.StatusBadRequest)
			return
		}

		app.errorJSON(w, err)
		return
	}

	if itemStock.Stock < int32(req.Qty) {
		app.errorJSON(w, appError.ErrOutOfStock, http.StatusBadRequest)
		return
	}

	orderId := uuid.New().String()
	_, err = app.InventoryService.ReserveStock(ctx, &rpc.ReserveStockRequest{
		ItemId:  req.ItemId,
		Version: itemStock.Version,
		OrderId: orderId,
	})
	if err != nil {
		_, err := app.DB.AddOrder(orderId, false, req)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		app.writeJson(w, http.StatusAccepted, model.OrderResponse{
			ID:     orderId,
			Status: "REJECTED",
		})
	}

	_, err = app.DB.AddOrder(orderId, true, req)
	if err != nil {
		_, err = app.InventoryService.ReleaseStock(ctx, &rpc.ReleaseStockRequest{
			ItemId:  req.ItemId,
			OrderId: orderId,
			Qty:     int32(req.Qty),
		})

		if err != nil {
			app.errorJSON(w, appError.ErrInternalServer)
			return
		}
		app.writeJson(w, http.StatusAccepted, model.OrderResponse{
			ID:     orderId,
			Status: "REJECTED",
		})
	}

	app.writeJson(w, http.StatusCreated, model.OrderResponse{
		ID:     orderId,
		Status: "CONFIRMED",
	})

}

func processError(app *HttpHandler, w http.ResponseWriter, err error) {
	ok, httpError := isHttpError(err)
	if ok {
		app.errorJSON(w, httpError.Err, httpError.StatusCode)
		return
	}

	app.errorJSON(w, appError.ErrInternalServer)
	return
}

func isHttpError(err error) (bool, *appError.HttpError) {
	var httpError *appError.HttpError
	if errors.As(err, &httpError) {
		return true, httpError
	}

	return false, nil
}
