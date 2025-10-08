package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONResponse struct {
	Message string `json:"message"`
}

func (app *application) writeJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {

	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("invalid request, must only containt single JSON")
	}

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {

	statusCode := http.StatusInternalServerError
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JSONResponse{
		Message: err.Error(),
	}

	return app.writeJson(w, statusCode, payload)
}
