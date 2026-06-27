package utils

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if hash == "" {
		t.Error("expected hash to be not empty")
	}

	if hash == "password" {
		t.Error("expected hash to be different from password")
	}
}

func TestComparePassword(t *testing.T) {
	pass := "password"
	hash, err := HashPassword(pass)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if !ComparePassword(hash, pass) {
		t.Error("expected password to match hash")
	}

	if ComparePassword(hash, "not"+pass) {
		t.Error("expected password to not match hash")
	}
}
