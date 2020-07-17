package main

import (
	"encoding/json"
	"net/http"
)

type healthCheck struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func healthCheckView(w http.ResponseWriter, r *http.Request) {
	data := healthCheck{
		Status:  "OK",
		Message: "Healthy",
	}
	applyResponseJSON(200, data, w)
}

func applyResponseJSON(code int, data interface{}, w http.ResponseWriter) {
	if code != 200 {
		w.WriteHeader(code)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		se := serverError{}
		se.applyServerError(err, "Server error", w)
		return
	}
	w.Write(jsonData)
	w.Header().Set("Content-Type", "application/json")
}

type serverError struct {
	Err          error  `json:"error"`
	ErrorMessage string `json:"message"`
}

func (serverErr *serverError) applyServerError(err error, errorMessage string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	se := &serverError{
		Err:          err,
		ErrorMessage: errorMessage,
	}
	jsonData, err := json.Marshal(se)
	w.Write(jsonData)
}
