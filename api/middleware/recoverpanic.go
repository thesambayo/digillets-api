package middleware

import (
	"fmt"
	"net/http"
)

func (middleware *Middleware) RecoverFromPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				responseWriter.Header().Set("Connection", "close")
				// Call the ServerError helper method to return a 500 Internal Server response.
				middleware.httpx.ServerErrorResponse(responseWriter, request, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(responseWriter, request)
	})
}
