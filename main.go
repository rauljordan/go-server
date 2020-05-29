package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rauljordan/go-server/middleware"
	"github.com/rauljordan/go-server/routes"
	"github.com/rauljordan/go-server/server"
)

func main() {
	// Initialize the server using environment variables.
	srv, err := server.New(context.Background(), &server.Config{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		JWTKey:      []byte(os.Getenv("JWT_KEY")),
		Port:        8080,
	})
	if err != nil {
		log.Fatalf("Could not initialize server: %v", err)
	}
	srv.Start(BindRoutes)
}

// BindRoutes adds defined routes to the multiplexer and
// adds required middlewares for the server.
func BindRoutes(srv server.Server, r *mux.Router) {
	// Middleware.
	r.Use(middleware.Authentication(srv))

	// Profiling.
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	// Metrics.
	r.Handle("/metrics", promhttp.Handler())

	// Users routes.
	r.HandleFunc("/signup", routes.Signup(srv)).Methods(http.MethodPost)
	r.HandleFunc("/login", routes.Login(srv)).Methods(http.MethodPost)
}
