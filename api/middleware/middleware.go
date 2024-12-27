package middleware

import (
	"github.com/thesambayo/digillets-api/api/httpx"
	"github.com/thesambayo/digillets-api/internal/config"
	"github.com/thesambayo/digillets-api/internal/data"
)

type Middleware struct {
	config config.Config
	httpx  *httpx.Utils
	models *data.Models
}

func New(cfg config.Config, httpx *httpx.Utils, models *data.Models) *Middleware {
	return &Middleware{
		cfg,
		httpx,
		models,
	}
}
