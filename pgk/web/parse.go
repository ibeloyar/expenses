package web

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

var (
	ErrIDMustBeenPosInt = errors.New("id must been positive integer")
)

// ParseQueryParams return map containing query params with keys from names slice and error
func ParseQueryParams(r *http.Request, names ...string) (map[string]string, error) {
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(names))

	for _, name := range names {
		if len(queryParams[name]) != 0 {
			result[name] = queryParams[name][0]
		}
	}

	return result, nil
}

// ParseIDFromURL return URL path id parameter (integer) with name from second function argument
// Example: "/user/{userID}" => id, err := ParseIDFromURL(r, "userID")
func ParseIDFromURL(r *http.Request, name string) (int, error) {
	idFromURL := r.PathValue(name)

	param, err := strconv.Atoi(idFromURL)
	if err != nil || param < 1 {
		return 0, ErrIDMustBeenPosInt
	}

	return param, nil
}
