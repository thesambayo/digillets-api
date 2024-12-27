package middleware

import (
	"net/http"

	"github.com/thesambayo/digillet-api/api/contexts"
)

func (middleware *Middleware) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		user := contexts.ContextGetUser(req)
		if user.IsAnonymous() {
			middleware.httpx.AuthenticationRequiredResponse(resWriter, req)
			return
		}

		next.ServeHTTP(resWriter, req)
	})
}
