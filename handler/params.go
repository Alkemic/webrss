package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Alkemic/go-route"
)

func requestIntParam(req *http.Request, key string) (int64, error) {
	rawValue, ok := route.GetParam(req, key)
	if !ok {
		return -1, errors.New("missing parameter")
	}
	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return -1, fmt.Errorf("cannot convert param '%s' to int: %w", key, err)
	}
	return int64(value), nil
}

func routeIntParam(key string, req *http.Request) (int, bool, error) {
	query := req.URL.Query()
	rawValue := query.Get(key)
	if rawValue == "" {
		return 0, false, fmt.Errorf("missing '%s' param", key)
	}
	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, true, fmt.Errorf("'%s' is not int: %w", key, err)
	}
	return value, true, nil
}
