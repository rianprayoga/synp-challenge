package handler

import (
	"errors"
	"net/http"
	appError "orders-app/internal/error"
	"orders-app/internal/model"
)

func (app *HttpHandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	var req model.OrderRequest
	err := app.readJSON(r, &req)
	if err != nil {
		processError(app, w, err)
		return
	}

	res, err := app.DB.AddOrder(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJson(w, http.StatusCreated, res)
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
