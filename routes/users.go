package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rauljordan/go-server/models"
	"github.com/rauljordan/go-server/server"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenExpiryLength = 20 * time.Minute // JWT tokens should expire after 20 mins.
	hashCost          = 8                // Standard hash cost for bcrypt password generation.
)

var (
	// Define the number of processed login attempts metric.
	loginsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "webserver_login_processed_total",
		Help: "The total number of processed login events",
	})
	// Define the number of successful login attempts metric.
	loginsSucceeded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "webserver_login_succeeded_total",
		Help: "The total number of succeeded login events",
	})
)

// Request type for authentication endpoints signup/login.
type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response type for authentication endpoints signup/login.
type authResponse struct {
	Token           []byte `json:"token"`
	TokenExpiration uint64 `json:"token_expiration"`
}

// Signup a user and return a JWT token + expiration timestamp.
func Signup(srv server.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal("No server dependency available")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &authRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), hashCost)
		userID, err := srv.Database().CreateUser(r.Context(), req.Email, hashedPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tokenString, expirationTimestamp, err := tokenStringForUser(srv.Config().JWTKey, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res := &authResponse{
			Token:           []byte(tokenString),
			TokenExpiration: expirationTimestamp,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

// Login a user and return a JWT token + expiration timestamp.
func Login(srv server.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal("No server dependency available")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		loginsProcessed.Inc()
		req := &authRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := srv.Database().User(r.Context(), req.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(req.Password)); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		tokenString, expirationTimestamp, err := tokenStringForUser(srv.Config().JWTKey, user.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res := &authResponse{
			Token:           []byte(tokenString),
			TokenExpiration: expirationTimestamp,
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		loginsSucceeded.Inc()
	}
}

// Generate a JWT token strict given a jwt key and a user id.
func tokenStringForUser(jwtKey []byte, userID uint64) (string, uint64, error) {
	expirationTime := time.Now().Add(tokenExpiryLength)
	claims := &models.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", 0, errors.Wrap(err, "could not sign token")
	}
	return tokenString, uint64(expirationTime.Unix()), nil
}
