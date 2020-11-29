package account

import (
	"context"
	"net/http"

	"github.com/Alkemic/webrss/repository"
)

var userCtxKey = struct{}{}

func SetUser(r *http.Request, user repository.User) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, userCtxKey, user)
	(*r) = *(r.WithContext(ctx))
}

func GetUser(r *http.Request) repository.User {
	userRaw := r.Context().Value(userCtxKey)
	val := userRaw.(repository.User)
	return val
}
