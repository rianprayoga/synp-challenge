package handler

import (
	"encoding/json"
	"errors"
	appError "inventories-app/internal/error"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type JSONResponse struct {
	Message string `json:"message"`
}

func (app *HttpHandler) writeJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {

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

func (app *HttpHandler) readJSON(r *http.Request, data interface{}) error {

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("invalid request, must only containt single JSON")
	}

	err = validate.Struct(data)
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	ve := validationErrors[0] // get the 1st error

	return &appError.HttpError{
		Err:        errors.New(ve.Field() + " does not match the requierment"),
		StatusCode: http.StatusBadRequest,
	}
}

func (app *HttpHandler) errorJSON(w http.ResponseWriter, err error, status ...int) error {

	statusCode := http.StatusInternalServerError
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JSONResponse{
		Message: err.Error(),
	}

	return app.writeJson(w, statusCode, payload)
}
