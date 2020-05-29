package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/rauljordan/go-server/models"
	"github.com/rauljordan/go-server/server"
)

// Authentication middleware validating JWT tokens
// included in http request headers.
func Authentication(srv server.Server) func(h http.Handler) http.Handler {
	if srv == nil {
		log.Fatal("Nil dependency was passed to authentication middleware")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !srv.ShouldAuthenticatePath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			token := strings.TrimSpace(r.Header.Get("Authorization"))
			checkParsedKey := func(*jwt.Token) (interface{}, error) {
				return srv.Config().JWTKey, nil
			}
			if _, err := jwt.ParseWithClaims(token, &models.Claims{}, checkParsedKey); err != nil {
				w.WriteHeader(401)
				w.Write([]byte("Unauthorized"))
				return
			}
			// TODO: Use claims.
			next.ServeHTTP(w, r)
		})
	}
}
