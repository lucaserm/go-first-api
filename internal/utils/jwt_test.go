package utils

import "testing"

func TestCreateJwt(t *testing.T) {
	token, err := CreateJWT("123")
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}
}
