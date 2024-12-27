package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
	"github.com/thesambayo/digillet-api/api/contexts"
	"github.com/thesambayo/digillet-api/internal/constants"
	"github.com/thesambayo/digillet-api/internal/data/users"
)

func (middleware *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		// Add the "Vary: Authorization" header to the response. This indicates to any
		// caches that the response may vary based on the value of the Authorization
		// header in the request.
		resWriter.Header().Add("Vary", "Authorization")

		// Retrieve the value of the Authorization header from the request. This will
		// return the empty string "" if there is no such header found.
		authorizationHeader := req.Header.Get("Authorization")

		// If there is no Authorization header found, use the ContextSetUser() helper
		// that we just made to add the AnonymousUser to the request context.
		// Then we call the next handler in the chain and
		// return without executing any of the code below.
		if authorizationHeader == "" {
			req = contexts.ContextSetUser(req, users.AnonymousUser)
			next.ServeHTTP(resWriter, req)
			return
		}

		// Otherwise, we expect the value of the Authorization header to be in the format
		// "Bearer <token>". We try to split this into its constituent parts, and if the
		// header isn't in the expected format we return a 401 Unauthorized response
		// using the invalidAuthenticationTokenResponse() helper (which we will create
		// in a moment).
		// headerParts := strings.Split(authorizationHeader, " ")
		headerParts := strings.Fields(authorizationHeader)
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			middleware.httpx.InvalidAuthenticationTokenResponse(resWriter, req)
			return
		}

		// Extract the actual authentication token from the header parts.
		token := headerParts[1]

		// Parse the JWT and extract the claims. This will return an error if the JWT
		// contents doesn't match the signature (i.e. the token has been
		// tampered with) or the algorithm isn't valid.
		claims, err := jwt.HMACCheck([]byte(token), []byte(middleware.config.Jwt.Secret))
		if err != nil {
			middleware.httpx.InvalidAuthenticationTokenResponse(resWriter, req)
			return
		}
		// Check if the JWT is still valid at this moment in time.
		if !claims.Valid(time.Now()) {
			middleware.httpx.InvalidAuthenticationTokenResponse(resWriter, req)
			return
		}
		// Check that the issuer is our application.
		// if claims.Issuer != "operations.oryoltd.org" {
		// 	middleware.httpx.InvalidAuthenticationTokenResponse(resWriter, req)
		// 	return
		// }
		// // Check that our application is in the expected audiences for the JWT.
		// if !claims.AcceptAudience("operations.oryoltd.org") {
		// 	middleware.httpx.InvalidAuthenticationTokenResponse(resWriter, req)
		// 	return
		// }

		// At this point, we know that the JWT is all OK and we can trust the data in
		// it. We extract the user ID from the claims subject
		UserPublicID := claims.Subject
		// Lookup the user record from the database.
		user, err := middleware.models.Users.GetByPublicId(UserPublicID)

		if err != nil {
			switch {
			case errors.Is(err, constants.ErrRecordNotFound):
				middleware.httpx.InvalidAuthenticationTokenResponse(resWriter, req)
			default:
				middleware.httpx.ServerErrorResponse(resWriter, req, err)
			}
			return
		}

		req = contexts.ContextSetUser(req, user)
		// Call the next handler in the chain.
		next.ServeHTTP(resWriter, req)
	})
}
