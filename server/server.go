package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rauljordan/go-server/internal/db"
)

// Define the routes in the server which do not
// require authentication to resolve.
var (
	noAuthPaths = []string{"login", "signup", "metrics", "debug"}
)

// Server defines a struct which can retrieve dependencies
// such as a database, configs, or more. It can also determine whether
// or not a url path requires JWT authentication.
type Server interface {
	Config() *Config
	Database() db.Database
	ShouldAuthenticatePath(string) bool
}

// Config struct for the server, including a database endpoint
// url, a secret JWT key for authentication, and a port value.
type Config struct {
	DatabaseUrl string
	JWTKey      []byte
	Port        int
}

// Broker allows the ability to launch a new API server
// using a passed in configuration.
type Broker struct {
	cfg    *Config
	router *mux.Router
	db     db.Database
}

// New instantiates a server from a configuration,
// initializing its dependencies such as access to a
// database.
func New(ctx context.Context, cfg *Config) (*Broker, error) {
	if cfg.DatabaseUrl == "" {
		return nil, errors.New("must provide a DATABASE_URL env value")
	}
	if string(cfg.JWTKey) == "" {
		return nil, errors.New("must provide a JWT_KEY env value")
	}
	bkr := &Broker{
		cfg: cfg,
	}
	serverDB, err := db.StartDB(ctx, cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("could not initialize database: %v", err)
	}
	bkr.db = serverDB
	return bkr, nil
}

// Start initializes a server by using a route-binding function
// and registering an http listener on a port.
func (bkr *Broker) Start(binder func(s Server, r *mux.Router)) {
	bkr.router = mux.NewRouter()
	binder(bkr, bkr.router)
	port := fmt.Sprintf(":%d", bkr.Config().Port)
	log.Printf("Starting API server on port %s", port)
	if err := http.ListenAndServe(port, bkr.router); errors.Is(err, http.ErrServerClosed) {
		log.Println("Server has shut down")
	} else {
		log.Fatal("Server has shut down unexpectedly")
	}
}

// Config allows retrieval of the server's
// configuration options
func (bkr *Broker) Config() *Config {
	return bkr.cfg
}

// Database allows retrieval of the server's
// database dependency.
func (bkr *Broker) Database() db.Database {
	return bkr.db
}

// ShouldAuthenticatePath determines if a URL
// path should require JWT authentication to
// succeed in its respective http handler.
func (bkr *Broker) ShouldAuthenticatePath(path string) bool {
	for _, p := range noAuthPaths {
		if strings.Contains(path, p) {
			return false
		}
	}
	return true
}
