// Package web provides a set of useful functions and types for building a simple web application.
//
// The package is written for version Go 1.22.0.
// Does not contain third-party libraries other than packages included in the Go language.
//
// The package allows:
// Writing responses to http requests in json format, with the most popular error codes.
// Extract parameters from the query string.
// For convenience, functions have been added to retrieve the most popular query parameters as page and limit, search
// and path param id.
// Correctly handle panic in the service (PanicRecover).
// Handle service errors with logging (slog.Logger).
package web
