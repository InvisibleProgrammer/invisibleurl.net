package users

import "testing"

func TestEmail_valid(t *testing.T) {
	email := "test@test.com"

	err := validateEmail(email)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmail_valid_accept_comment(t *testing.T) {
	email := "test+1@test.com"

	err := validateEmail(email)

	if err == nil {
		return
	}

	if validationError, ok := err.(*ValidationError); ok && len(validationError.ValidationErrors) > 0 {
		t.Fatalf("invalid email: %s; validation errors: %v", email, validationError.ValidationErrors)
	} else {
		t.Fatalf("unkown error: %v", err)
	}
}

func TestEmail_no_at(t *testing.T) {
	email := "testtest.com"

	err := validateEmail(email)

	if err == nil {
		t.Fatal("email should contain a @ character")
	}
}

func TestEmail_no_dot_in_domain(t *testing.T) {
	email := "test@testcom"

	err := validateEmail(email)

	if err == nil {
		t.Fatal("email should contain a @ character")
	}
}
