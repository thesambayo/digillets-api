package contexts

import (
	"context"
	"net/http"

	"github.com/thesambayo/digillets-api/internal/data/users"
)

type contextKey string

const userContextKey = contextKey("USER")

// ContextSetUser() helper add the user information to the request context.
func ContextSetUser(req *http.Request, user *users.User) *http.Request {
	ctx := context.WithValue(req.Context(), userContextKey, user)
	return req.WithContext(ctx)
}

func ContextGetUser(req *http.Request) *users.User {
	user, ok := req.Context().Value(userContextKey).(*users.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
