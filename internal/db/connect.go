package db

import (
	"context"
	"fmt"
	pgx "github.com/jackc/pgx/v5"
	"time"
)

type Config struct {
	Host            string        `env:"DB_HOST,default=localhost"`
	Port            int           `env:"DB_PORT,default=5432"`
	Name            string        `env:"DB_NAME,required"`
	User            string        `env:"DB_USER,required"`
	Password        string        `env:"DB_PASSWORD,required"`
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
