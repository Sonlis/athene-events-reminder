package database

import (
	"context"
	"fmt"
	pgx "github.com/jackc/pgx/v5"
	"github.com/sethvargo/go-envconfig"
	"os"
)

type Config struct {
	Host            string        `env:"POSTGRES_HOST,required"`
	Port            int           `env:"POSTGRES_PORT,required"`
	Name            string        `env:"DB_NAME,required"`
	User            string        `env:"POSTGRES_USER,required"`
	Password        string        `env:"POSTGRES_PASSWORD,required"`
}

type DB struct {
	*pgx.Conn
}

func (c *Config) ConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.User, string(c.Password), c.Host, c.Port, c.Name)
}

func Connect(config *Config) (*DB, error) {
	connectionString := config.ConnString()
	conn, err := pgx.Connect(context.Background(), connectionString)
	return &DB{conn}, err
}

// InitTestConfig is a helper function to init a db connection during tests.
func InitTestConfig() (Config, error) {
	var config Config
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5447")
	os.Setenv("POSTGRES_USER", "test_user")
	os.Setenv("POSTGRES_PASSWORD", "test_pass")
	os.Setenv("DB_NAME", "test")
	err := envconfig.Process(context.Background(), &config)
	return config, err
}
