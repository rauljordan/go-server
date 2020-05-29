package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/rauljordan/go-server/internal/mocks"
	"github.com/rauljordan/go-server/models"
	"github.com/rauljordan/go-server/server"
)

func TestSignup_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	email := "someone@mail.com"
	password := "123456"
	jwtKey := []byte("secret")
	userID := uint64(1)

	reqBody := &authRequest{
		Email:    email,
		Password: password,
	}
	enc, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	mockDB := mocks.NewMockDatabase(ctrl)
	mockServer := mocks.NewMockServer(ctrl)

	// Expectations for mock server calls.
	mockServer.EXPECT().Config().Return(&server.Config{
		JWTKey: jwtKey,
	})
	mockServer.EXPECT().Database().Return(mockDB)

	// Expectations for mock DB calls.
	mockDB.EXPECT().CreateUser(
		gomock.Any(), email, gomock.Any(),
	).Return(userID, nil)

	// Create a request to pass to our handler.
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(enc))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Signup(mockServer))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status, got %v want %v", status, http.StatusOK)
	}

	res := &authResponse{}
	if err := json.NewDecoder(rr.Body).Decode(res); err != nil {
		t.Fatal(err)
	}

	// We check we are able to parse the JWT properly and check
	// its claims contain user ID number 1.
	checkParsedKey := func(*jwt.Token) (interface{}, error) {
		return jwtKey, nil
	}
	claims := &models.Claims{}
	if _, err := jwt.ParseWithClaims(string(res.Token), claims, checkParsedKey); err != nil {
		t.Fatal(err)
	}
	if claims.UserID != 1 {
		t.Errorf("Expected userID 1, received %d", claims.UserID)
	}
}
