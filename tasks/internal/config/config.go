package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

// Config общий конфиг
type Config struct {
	Env         string `env:"ENV" envDefault:"prod"`
	Postgres    Postgres
	HTTPServer  HTTPServer
	AuthGRPC    AuthGRPC
	TokenConfig TokenConfig
}

type HTTPServer struct {
	Host        string `env:"HTTP_HOST" envDefault:"localhost"`
	Port        string `env:"HTTP_PORT" envDefault:"8081"`
	Address     string
	Timeout     time.Duration `env:"HTTP_TIMEOUT" envDefault:"5s"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	User        string        `env:"HTTP_USER" envDefault:"user"`
	Password    string        `env:"HTTP_PASSWORD" envDefault:"password"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"54322"`
	User     string `env:"POSTGRES_USER" envDefault:"user"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	Db       string `env:"POSTGRES_DB" envDefault:"postgres"`
	SslMode  string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	DSN      string
}

type AuthGRPC struct {
	Host    string `env:"AUTH_GRPC_HOST" envDefault:"localhost"`
	Port    string `env:"AUTH_GRPC_PORT" envDefault:"50051"`
	Address string
}

type TokenConfig struct {
	AccessTokenSecretKey string `env:"ACCESS_TOKEN_SECRET"`
	AuthPrefix           string `env:"AUTH_PREFIX" envDefault:"Bearer "`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("loading config from env is failed: %w", err)
	}
	buildHTTPAddress(&cfg.HTTPServer)
	buildDSN(&cfg.Postgres)
	buildGRPCAddress(&cfg.AuthGRPC)

	return cfg, nil
}

func buildHTTPAddress(httpserver *HTTPServer) {
	httpserver.Address = fmt.Sprintf("%s:%s", httpserver.Host, httpserver.Port)
}

func buildGRPCAddress(authGRPC *AuthGRPC) {
	authGRPC.Address = fmt.Sprintf("%s:%s", authGRPC.Host, authGRPC.Port)
}

func buildDSN(p *Postgres) {
	p.DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.Db, p.SslMode)
}
