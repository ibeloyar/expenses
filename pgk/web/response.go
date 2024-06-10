package web

import (
	"encoding/json"
	"net/http"
)

// WriteOK send on client response with status code 200 and data
// if data == nil, will send response without body
func WriteOK(w http.ResponseWriter, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// WriteCreated send on client response with status code 201 and data
// if data == nil, will send response without body
func WriteCreated(w http.ResponseWriter, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// WriteNoContent send on client response with status code 204 and data
// if data == nil, will send response without body
// used for delete response
func WriteNoContent(w http.ResponseWriter, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// WriteUnauthorized send on client response with status code 401 and error text
func WriteUnauthorized(w http.ResponseWriter, e error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	err := json.NewEncoder(w).Encode(&WebError{
		Code:    http.StatusUnauthorized,
		Message: e.Error(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// WriteForbidden send on client response with status code 403 and error text
func WriteForbidden(w http.ResponseWriter, e error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)

	err := json.NewEncoder(w).Encode(&WebError{
		Code:    http.StatusForbidden,
		Message: e.Error(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// WriteBadRequest send on client response with status code 400 and error text
func WriteBadRequest(w http.ResponseWriter, e error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	err := json.NewEncoder(w).Encode(&WebError{
		Code:    http.StatusBadRequest,
		Message: e.Error(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// WriteNotFound send on client response with status code 404 and error text
func WriteNotFound(w http.ResponseWriter, e error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	err := json.NewEncoder(w).Encode(&WebError{
		Code:    http.StatusNotFound,
		Message: e.Error(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// WriteServerError send on client response with status code 500 and text "internal server error"
func WriteServerError(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	err := json.NewEncoder(w).Encode(&WebError{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
