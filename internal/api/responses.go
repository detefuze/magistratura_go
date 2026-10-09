package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func BadRequest(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusBadRequest, ErrorResponse{Error: msg})
}

func NotFound(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusNotFound, ErrorResponse{Error: msg})
}

func Internal(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusInternalServerError, ErrorResponse{Error: msg})
}

// UnprocessableEntity — 422 для невалидных данных
func UnprocessableEntity(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusUnprocessableEntity, ErrorResponse{Error: msg})
}

// NoContent — 204 без тела, для успешного DELETE
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
