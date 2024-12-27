package routes

import (
	"expvar"
	"net/http"

	"github.com/thesambayo/digillets-api/api/httpx"
	"github.com/thesambayo/digillets-api/api/middleware"
	"github.com/thesambayo/digillets-api/internal/config"
	"github.com/thesambayo/digillets-api/internal/data"
)

type Routes struct {
	config config.Config
	httpx  *httpx.Utils
	models *data.Models
}

func Handlers(cfg config.Config, models *data.Models, httpx *httpx.Utils) http.Handler {
	router := http.NewServeMux()
	middleware := middleware.New(cfg, httpx, models)

	routes := &Routes{
		config: cfg,
		models: models,
		httpx:  httpx,
	}

	// healthcheck
	router.HandleFunc("GET /{$}", routes.HealthcheckHandler)
	// Register a new GET /debug/vars endpoint pointing to the expvar handler.
	router.HandleFunc("GET /debug/vars", expvar.Handler().ServeHTTP)

	// USERS
	router.HandleFunc("POST /v1/users/register", routes.CreateUser)
	router.HandleFunc("POST /v1/users/login", routes.AuthenticateUser)
	router.HandleFunc(
		"GET /v1/users/profile",
		middleware.RequireAuthenticatedUser(routes.GetUserProfile),
	)

	// middlewares usages around servemux
	return middleware.Metrics( // first
		middleware.RecoverFromPanic(
			middleware.EnableCORS(
				middleware.RateLimit(
					middleware.Authenticate( // last
						router,
					),
				),
			),
		),
	)
}
