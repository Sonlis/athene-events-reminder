# athene-events-reminder

Telegram bot to remind about Athene events.

## Running locally

The following environment variables are required:
```sh
export POSTGRES_HOST= # Host on which the postgres DB is running.
export DB_NAME= # Name of the db to connect to, within postgres.
export POSTGRES_PORT= # Port of postgres db to connect to.
export POSTGRES_USER= # Postgres user to use to connect.
export POSTGRES_PASSWORD= # Password of the user.
```
Then:
```sh
go mod tidy
go run cmd/server/server.go
```

## Running tests

Make sure that docker is installed locally.
```sh
make test
```
