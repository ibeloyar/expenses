// Package web provides a set of useful functions and types
// for building a simple web application
package web

// WebError will send in any response with status codes 4xx and 5xx
type WebError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
