package db

import (
	"context"
	"fmt"
	pgx "github.com/jackc/pgx/v5"
	"github.com/sethvargo/go-envconfig"
	"os"
	"time"
)

type Config struct {
	Host            string        `env:"DB_HOST"`
	Port            int           `env:"DB_PORT"`
	Name            string        `env:"DB_NAME"`
	User            string        `env:"DB_USER"`
	Password        string        `env:"DB_PASSWORD"`
	ConnTimeout     time.Duration `env:"DB_CONN_TIMEOUT,default=30s"`
	ConnAttempts    int           `env:"DB_CONN_ATTEMPTS,default=10"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS,default=2"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS,default=20"`
	MaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME,default=30m"`
	MaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME,default=0"`
}

func (c *Config) ConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.User, string(c.Password), c.Host, c.Port, c.Name)
}

func Connect(config *Config) (*pgx.Conn, error) {
	connectionString := config.ConnString()
	conn, err := pgx.Connect(context.Background(), connectionString)
	return conn, err
}

// InitTestConfig is a helper function to init a db connection during tests.
func InitTestConfig() (Config, error) {
	var config Config
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5447")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")
	os.Setenv("DB_NAME", "test")
	err := envconfig.Process(context.Background(), &config)
	return config, err
}
