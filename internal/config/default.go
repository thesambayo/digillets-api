package config

// digilletsadmin password'passwordllets';
// postgres://{rolename}:{passsword}@{dbURL}/{dbname}
// postgres://digilletsadmin:passwordllets@localhost/digillets
// go run ./cmd/ -cors-trusted-origins='http://localhost:5500 https://oryo.com'
// psql --host=localhost --dbname=digillets --username=digilletsadmin
func DefaultConfig() Config {
	return Config{
		Port: 5500,
		Env:  "development",
		Limiter: Limiter{
			Rps:     2,
			Burst:   4,
			Enabled: true,
		},
		Cors: Cors{
			TrustedOrigins: []string{"http://localhost:5500"},
		},
		Jwt: Jwt{
			Secret: "pei3einoh0Beem6uM6Ungohn2heiv5lah1ael4joopie5JaigeikoozaoTew2Eh6",
		},
		DB: DB{
			Dsn:          "postgres://digilletsadmin:passwordllets@localhost/digillets?sslmode=disable",
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxIdleTime:  "15m",
		},
	}
}
