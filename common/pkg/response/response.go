// Package response provides helpers for writing consistent JSON HTTP responses.
package response

import (
	"encoding/json"
	"net/http"
)

// Success writes a 200-OK JSON body: {"data": v}.
func Success(w http.ResponseWriter, v interface{}) {
	write(w, http.StatusOK, map[string]interface{}{"data": v})
}

// Created writes a 201-Created JSON body.
func Created(w http.ResponseWriter, v interface{}) {
	write(w, http.StatusCreated, map[string]interface{}{"data": v})
}

// WriteError writes an error JSON body.
func WriteError(w http.ResponseWriter, code int, msg string) {
	write(w, code, map[string]string{"error": msg})
}

func write(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
