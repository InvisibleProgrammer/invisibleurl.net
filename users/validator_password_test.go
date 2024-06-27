package users

import (
	"testing"
)

func TestPassword_valid(t *testing.T) {
	password := "12Debrecen99!"

	err := validatePassword(password, password)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestPassword_confirm_doesnot_match(t *testing.T) {
	password := "12Debrecen99!"
	confirmPassword := "12Debrecen99"

	err := validatePassword(password, confirmPassword)
	if err == nil {
		t.Fatal("validation should fail if password and confirm password does not match")
	}
}
