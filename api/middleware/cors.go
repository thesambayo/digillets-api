package middleware

import "net/http"

func (middleware *Middleware) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Origin" header.
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		// Get the value of the request's Origin header.
		origin := r.Header.Get("Origin")

		// Only run this if there's an Origin request header present AND at least one trusted origin is configured.
		if origin != "" && len(middleware.config.Cors.TrustedOrigins) != 0 {
			// Loop through the list of trusted origins, checking to see if the request
			// origin exactly matches one of them.
			for i := range middleware.config.Cors.TrustedOrigins {
				if origin == middleware.config.Cors.TrustedOrigins[i] {
					// If there is a match, then set a "Access-Control-Allow-Origin"
					// response header with the request origin as the value.
					w.Header().Set("Access-Control-Allow-Origin", origin)
					// Check if the request has the HTTP method OPTIONS and contains the
					// "Access-Control-Request-Method" header. If it does, then we treat
					// it as a preflight request.
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// Set the necessary preflight response headers, as discussed previously.
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
						// Write the headers along with a 200 OK status and return from
						// the middleware with no further action.
						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
