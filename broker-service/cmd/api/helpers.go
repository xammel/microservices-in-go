package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJson(writer http.ResponseWriter, request *http.Request, data any) error {
	maxBytes := 1048576 // 1 MB
	request.Body = http.MaxBytesReader(writer, request.Body, int64(maxBytes))
	
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
 }

 func (app *Config) writeJson(writer http.ResponseWriter, status int, data any, headers ...http.Header) error {
	output, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			writer.Header()[key] = value
		}
	
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, err = writer.Write(output)
	if err != nil {
		return err
	}

	return nil
 }

 func (app *Config) errorJson(writer http.ResponseWriter, err error, status ...int) error {
	
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJson(writer, statusCode, payload)
 }