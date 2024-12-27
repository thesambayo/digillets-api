package main

import (
	"os"
	"sync"

	_ "github.com/lib/pq"

	"github.com/thesambayo/digillet-api/api/httpx"
	"github.com/thesambayo/digillet-api/internal/config"
	"github.com/thesambayo/digillet-api/internal/data"
	"github.com/thesambayo/digillet-api/internal/jsonlog"
)

type application struct {
	config config.Config
	logger *jsonlog.Logger
	models *data.Models
	httpx  *httpx.Utils
	wg     *sync.WaitGroup
}

func main() {
	cfg := config.GetConfig()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		httpx:  httpx.New(logger),
		models: data.New(db),
		wg:     &sync.WaitGroup{},
	}

	err = app.serve()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
}

// functional option pattern
// 	app := NewApplication(
// 	WithConfig(cfg),
// 	WithLogger(jsonlog.New(os.Stdout, jsonlog.LevelInfo)),
// 	// WithModels(models),
// )
// type Option func(*application)

// func WithLogger(logger *jsonlog.Logger) Option {
// 	return func(a *application) {
// 		a.logger = logger
// 	}
// }

// func WithConfig(cfg config.Config) Option {
// 	return func(a *application) {
// 		a.config = cfg
// 	}
// }

// func NewApplication(opts ...Option) *application {
// 	app := &application{
// 		// make app defaults
// 		// config: config.DefaultConfig(),
// 	}
// 	for _, opt := range opts {
// 		opt(app)
// 	}
// 	return app
// }
