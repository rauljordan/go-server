package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/rauljordan/go-server/internal/mocks"
	"github.com/rauljordan/go-server/middleware"
	"github.com/rauljordan/go-server/models"
	"github.com/rauljordan/go-server/server"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	t.Run("skip when authentication is not needed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		srv := mocks.NewMockServer(ctrl)
		handler := middleware.Authentication(srv)

		req := httptest.NewRequest("GET", "http://not-authorized.com", nil)

		nextCalled := false
		nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			nextCalled = true
		})

		srv.EXPECT().ShouldAuthenticatePath(req.URL.Path).Return(false)

		handlerToTest := handler(nextHandler)
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)

		assert.True(t, nextCalled, "next handler was not called, but it should have been")
	})

	t.Run("no token present within header", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		srv := mocks.NewMockServer(ctrl)
		handler := middleware.Authentication(srv)

		req := httptest.NewRequest("GET", "http://no-auth-token.com", nil)
		rec := httptest.NewRecorder()

		nextCalled := false
		nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			nextCalled = true
		})

		srv.EXPECT().ShouldAuthenticatePath(req.URL.Path).Return(true)

		handlerToTest := handler(nextHandler)
		handlerToTest.ServeHTTP(rec, req)

		assert.False(t, nextCalled, "next handler was called, but it should not have been")
		assert.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode)
		assert.Equal(t, "Unauthorized", rec.Body.String())
	})
	
	t.Run("malformed token within header", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		srv := mocks.NewMockServer(ctrl)
		handler := middleware.Authentication(srv)

		req := httptest.NewRequest("GET", "http://malformed-auth-token.com", nil)
		req.Header.Set("Authorization", "foo")
		rec := httptest.NewRecorder()

		nextCalled := false
		nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			nextCalled = true
		})

		srv.EXPECT().ShouldAuthenticatePath(req.URL.Path).Return(true)

		handlerToTest := handler(nextHandler)
		handlerToTest.ServeHTTP(rec, req)

		assert.False(t, nextCalled, "next handler was called, but it should not have been")
		assert.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode)
		assert.Equal(t, "Unauthorized", rec.Body.String())
	})

	t.Run("success - token without whitespace", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		srv := mocks.NewMockServer(ctrl)
		handler := middleware.Authentication(srv)

		JWTKey:= []byte("secret")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
			UserID: 1,
			StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Second).Unix()},
		})
		tokenString, err := token.SignedString(JWTKey)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("GET", "http://success.com", nil)
		req.Header.Set("Authorization", tokenString)
		rec := httptest.NewRecorder()

		nextCalled := false
		nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			nextCalled = true
		})

		srv.EXPECT().ShouldAuthenticatePath(req.URL.Path).Return(true)
		srv.EXPECT().Config().Return(&server.Config{JWTKey: JWTKey})

		handlerToTest := handler(nextHandler)
		handlerToTest.ServeHTTP(rec, req)

		assert.True(t, nextCalled)
	})

	t.Run("success - token with whitespace", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		srv := mocks.NewMockServer(ctrl)
		handler := middleware.Authentication(srv)

		JWTKey:= []byte("secret")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
			UserID: 1,
			StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Second).Unix()},
		})
		tokenString, err := token.SignedString(JWTKey)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("GET", "http://success.com", nil)
		req.Header.Set("Authorization", fmt.Sprintf("   %s   ", tokenString))
		rec := httptest.NewRecorder()

		nextCalled := false
		nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			nextCalled = true
		})

		srv.EXPECT().ShouldAuthenticatePath(req.URL.Path).Return(true)
		srv.EXPECT().Config().Return(&server.Config{JWTKey: JWTKey})

		handlerToTest := handler(nextHandler)
		handlerToTest.ServeHTTP(rec, req)

		assert.True(t, nextCalled)
	})
}
