package main

import (
	b64 "encoding/base64"
	"errors"
	appError "inventories-app/internal/error"
	"inventories-app/internal/model"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (app *application) GetItems(w http.ResponseWriter, r *http.Request) {

	var pageSize int
	size := r.URL.Query().Get("size")
	if size == "" {
		pageSize = 10
	} else {
		tmp, err := strconv.Atoi(size)
		if err != nil {
			log.Println(err)
			app.errorJSON(w, appError.ErrUnexpected)
			return
		}
		pageSize = tmp
	}

	cursor := r.URL.Query().Get("cursor")

	items, err := getItems(app, pageSize, cursor)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, appError.ErrUnexpected)
		return
	}

	var resp model.PageResponse[*model.Item]

	if len(items) > pageSize {
		resp.Data = items[:pageSize]
	} else {
		resp.Data = items
	}

	if len(items) > pageSize {
		last := resp.Data[len(resp.Data)-1]
		cursor := b64.StdEncoding.EncodeToString([]byte(last.CreatedAt.Format(time.RFC3339)))
		resp.NextCursor = cursor
	}

	_ = app.writeJson(w, http.StatusOK, resp)
}

func (app *application) GetItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	item, err := app.DB.GetItem(id)
	if err != nil {
		ok, httpError := isHttpError(err)
		if ok {
			app.errorJSON(w, httpError.Err, httpError.StatusCode)
			return
		}

		app.errorJSON(w, appError.ErrUnexpected)
		return
	}

	app.writeJson(w, http.StatusOK, item)
}

func (app *application) AddItem(w http.ResponseWriter, r *http.Request) {
	var req model.CreateItem
	err := app.readJSON(r, &req)
	if err != nil {
		ok, httpError := isHttpError(err)
		if ok {
			app.errorJSON(w, httpError.Err, httpError.StatusCode)
			return
		}

		app.errorJSON(w, err)
		return
	}

	res, err := app.DB.AddItem(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJson(w, http.StatusCreated, res)

}
func (app *application) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := app.DB.GetItem(id)
	if err != nil {
		ok, httpError := isHttpError(err)
		if ok {
			app.errorJSON(w, httpError.Err, httpError.StatusCode)
			return
		}

		app.errorJSON(w, appError.ErrUnexpected)
		return
	}

	if err := app.DB.DeleteItem(id); err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJson(
		w,
		http.StatusOK,
		JSONResponse{
			Message: "item deleted",
		},
	)

}

func isHttpError(err error) (bool, *appError.HttpError) {
	var httpError *appError.HttpError
	if errors.As(err, &httpError) {
		return true, httpError
	}

	return false, nil
}

func getItems(app *application, size int, cursor string) ([]*model.Item, error) {
	if cursor == "" {
		items, err := app.DB.GetItems(size + 1)
		if err != nil {
			return nil, err
		}
		return items, nil
	}

	nextCursor, err := b64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	c, err := time.Parse(time.RFC3339, string(nextCursor))
	if err != nil {
		return nil, err
	}

	items, err := app.DB.GetItemsWithCursor(size+1, c)
	if err != nil {
		return nil, err
	}
	return items, nil
}
