package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Alkemic/go-route"
)

func GetIntParam(req *http.Request, key string) (int64, error) {
	rawValue, ok := route.GetParam(req, key)
	if !ok {
		return -1, errors.New("missing parameter")
	}
	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return -1, fmt.Errorf("cannot convert param '%s' to int: ", key, err)
	}
	return int64(value), nil
}
