package shared

import (
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Error writes an error response
func Error(w http.ResponseWriter, err *AppError) {
	JSON(w, err.Status, err)
}

// ErrorFromErr writes an error response from a standard error
func ErrorFromErr(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*AppError); ok {
		Error(w, appErr)
		return
	}
	JSON(w, http.StatusInternalServerError, map[string]string{
		"error":   "internal_error",
		"message": err.Error(),
	})
}

// Success writes a success response with status 200
func Success(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, data)
}

// Created writes a success response with status 201
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, data)
}

// DecodeJSON decodes JSON from request body
func DecodeJSON(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}


