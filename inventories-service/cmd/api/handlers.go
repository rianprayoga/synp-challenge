package main

import (
	b64 "encoding/base64"
	"fmt"
	"inventories-app/internal/model"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "inventoris app running")
}

func (app *application) GetItems(w http.ResponseWriter, r *http.Request) {

	var pageSize int
	size := r.URL.Query().Get("size")
	if size == "" {
		pageSize = 10
	} else {
		tmp, err := strconv.Atoi(size)
		if err != nil {
			log.Println(err)
			app.errorJSON(w, fmt.Errorf("something went wrong"))
			return
		}
		pageSize = tmp
	}

	cursor := r.URL.Query().Get("cursor")

	items, err := getItems(app, pageSize, cursor)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, fmt.Errorf("something went wrong"))
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

func getItems(app *application, size int, cursor string) ([]*model.Item, error) {
	if cursor == "" {
		items, err := app.DB.GetItems(size + 1)
		if err != nil {
			return nil, err
		}
		return items, nil
	}

	nextCursor, err := b64.StdEncoding.DecodeString(cursor)
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
