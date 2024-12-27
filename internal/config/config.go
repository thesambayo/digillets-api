package config

import (
	"flag"
	"strings"
)

type Cors struct {
	TrustedOrigins []string
}

type Jwt struct {
	Secret string
}

// limiter struct contains fields for the requests-per-second and burst values,
// and a boolean field which we can use to enable/disable rate limiting altogether.s
type Limiter struct {
	// Rps denotes requests-per-second
	Rps     float64
	Burst   int
	Enabled bool
}

type DB struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

// Config holds shared configuration settings.
type Config struct {
	Port    int
	Env     string
	Jwt     Jwt
	Cors    Cors
	Limiter Limiter
	DB      DB
}

// GetConfig creates and returns a new Config.
func GetConfig() Config {
	cfg := Config{}

	// command line flags to read the setting values into the config struct.
	// using defaultConfig to set defaults
	flag.IntVar(&cfg.Port, "port", DefaultConfig().Port, "API server port")
	flag.StringVar(&cfg.Env, "env", DefaultConfig().Env, "Environment (development|staging|production)")

	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", DefaultConfig().Limiter.Rps, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", DefaultConfig().Limiter.Burst, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", DefaultConfig().Limiter.Enabled, "Enable rate limiter")

	flag.StringVar(&cfg.DB.Dsn, "db-dsn", DefaultConfig().DB.Dsn, "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", DefaultConfig().DB.MaxOpenConns, "PostgresSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", DefaultConfig().DB.MaxIdleConns, "PostgresSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", DefaultConfig().DB.MaxIdleTime, "PostgresSQL max connection idle time")

	cfg.Cors.TrustedOrigins = DefaultConfig().Cors.TrustedOrigins
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.Cors.TrustedOrigins = strings.Fields(val)
		return nil
	})

	flag.StringVar(&cfg.Jwt.Secret, "jwt-secret", DefaultConfig().Jwt.Secret, "JWT secret")
	flag.Parse()
	return cfg
}
