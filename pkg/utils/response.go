package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (res *Response) Write(p []byte) (int, error) {
	return len(p), nil
}

func JSON(w http.ResponseWriter, statusCode int, data *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, &Response{
			Message: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
