package web

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

var (
	ErrIDMustBeenPosInt    = errors.New("id must been positive integer")
	ErrPageMustBeenPosInt  = errors.New("page must been positive integer greater than zero")
	ErrLimitMustBeenPosInt = errors.New("limit must been positive integer greater than zero")
)

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

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

// ParseQueryPagination return *Pagination struct from query params
// and takes a pointer to the request object and initialization value for pagination (can be nil).
//
// *Pagination contains limit and page fields (int)
func ParseQueryPagination(r *http.Request, init *Pagination) (*Pagination, error) {
	p := &Pagination{}

	if init != nil {
		p = init
	}

	qp, err := ParseQueryParams(r, "limit", "page")
	if err != nil {
		return nil, err
	}

	if qp["page"] != "" {
		qpPage, err := strconv.Atoi(qp["page"])
		if err != nil {
			return nil, ErrPageMustBeenPosInt
		}

		p.Page = qpPage
	}
	if p.Page < 1 {
		return nil, ErrLimitMustBeenPosInt
	}

	if qp["limit"] != "" {
		qpLimit, err := strconv.Atoi(qp["limit"])
		if err != nil {
			return nil, ErrLimitMustBeenPosInt
		}

		p.Limit = qpLimit
	}
	if p.Limit < 1 {
		return nil, ErrLimitMustBeenPosInt
	}

	return p, nil
}

// ParseSearchString return search string from query params
func ParseSearchString(r *http.Request) (string, error) {
	qp, err := ParseQueryParams(r, "search")
	if err != nil {
		return "", err
	}

	return qp["search"], nil
}
