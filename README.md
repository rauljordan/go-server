# Go-Server

A full example of a pure go API server with JWT authentication, postgres, prometheus, and pprof for instrumentation. This project includes examples on how to test a real Go webserver ready for production. Long story short: I was fed-up with _too many_ "hello world" API servers out there on the Internet and tutorial videos and decided to create one that was actually useful, included auth, and included postgres db access!

## Installing

```
go get github.com/rauljordan/go-server
```

## Initializing Postgres + Database Migrations

1. Install docker for your operating system
2. Run `docker-compose -f docker-compose.dev.yml up -d`
3. Run `export POSTGRESQL_URL='postgres://postgres:password@localhost:5432/go-server?sslmode=disable'`
4. Run `make migrate`

You may need to install `golang-migrate` [here](https://github.com/golang-migrate/migrate) if the last step fails for you.

## Building and Running the Server

1. `make server`
2. `./go-server`

```
2020/05/28 23:30:06 Established db connection
2020/05/28 23:30:06 Starting API server on port ::8080
```

Try creating a new user!

```
curl -d '{"email": "someone@email.com", "password": "123456"}' http://localhost:8080/signup
```

You'll get your JWT token and expiration timestamp:
```
{"token":"ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SjFjMlZ5WDJsa0lqb3lMQ0psZUhBaU9qRTFPVEEzTWpnek5EZDkuWGNQdXE0RFpYYU5Ia3kxRkc1cHNKTVlIN3Z0cUZYMWZzZk9Fal80SFJBMA==","token_expiration":1590728347}
```

## Running Tests

```
go test -v ./...
```

## Instrumenting With PProf

Prom metrics are available at `http://localhost:8080/metrics`, and pprof is also available at `http://localhost:8080/debug/pprof` by default. You can view a flame graph of the heap by doing:

```
go tool pprof -http localhost:7070 http://localhost:8080/debug/pprof/heap
```

## Run With Grafana + Prometheus

This will spin up a production config docker-compose with the server included, so make sure you're not currently running the server on the side before doing this.

```
Run `docker-compose -f docker-compose.prod.yml up -d`
```

1. Navigate to http://localhost:3000
2. User: admin, password: admin
3. Change the data source of grafana to `http://prometheus:9090`
4. Create a dashboard, add a new panel with any metrics from the Go server you wish to see
