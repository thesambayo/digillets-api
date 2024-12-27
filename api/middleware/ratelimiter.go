package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IP-based rate limiter
// this pattern for rate-limiting will only work if your API application is running on a single-machine.
// If your infrastructure is distributed, with your application running on multiple servers behind a load balancer,
// then you’ll need to use an alternative approach 👇🏼
// 1. If you’re using HAProxy or Nginx as a load balancer or reverse proxy,
// both of these have built-in functionality for rate limiting that it would probably be sensible to use.
// 2. Alternatively, you could use a fast database like Redis to maintain a request count for clients,
// running on a server which all your application servers can communicate with.
func (middleware *Middleware) RateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	// declare a mutex and a map to hold the clients IP addresses and rate limiters
	var (
		mutex   sync.Mutex
		clients = make(map[string]*client)
	)
	// Launch a background goroutine which removes old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock the mutex to prevent any rate limiter checks from happening while the cleanup is taking place.
			mutex.Lock()
			// Loop through all clients. If they haven't been seen within the last three
			// minutes, delete the corresponding entry from the map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			// Importantly, unlock the mutex when the cleanup is complete.
			mutex.Unlock()
		}
	}()

	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		if middleware.config.Limiter.Enabled {
			// extract the client's IP address from the request.
			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				middleware.httpx.ServerErrorResponse(resWriter, req, err)
				return
			}
			// Lock the mutex to prevent this code from being executed concurrently.
			mutex.Lock()

			// Check to see if the IP address already exists in the map. If it doesn't, then
			// Create and add a new client struct to the map
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(middleware.config.Limiter.Rps), middleware.config.Limiter.Burst),
				}
			}
			// Update the last seen time for the client.
			clients[ip].lastSeen = time.Now()

			// Call the Allow() method on the rate limiter for the current IP address. If
			// the request isn't allowed, unlock the mutex and send a 429 Too Many Requests
			// response, just like before.
			if !clients[ip].limiter.Allow() {
				mutex.Unlock()
				middleware.httpx.RateLimitExceededResponse(resWriter, req)
				return
			}

			// Very importantly, unlock the mutex before calling the next handler in the chain.
			// Notice that we DON'T use defer to unlock the mutex, as that would mean
			// that the mutex isn't unlocked until all the handlers downstream of this
			// middleware have also returned.
			mutex.Unlock()
		}

		next.ServeHTTP(resWriter, req)
	})
}
