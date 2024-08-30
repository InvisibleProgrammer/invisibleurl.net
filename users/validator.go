package users

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	ValidationErrors []string
	Err              error
}

func (validationError *ValidationError) Error() string {
	return validationError.Err.Error()
}

func validatePassword(password string) error {
	validationErrors := make([]string, 0)

	if len(password) < 8 {
		validationErrors = append(validationErrors, "password minimum length is 8 characters!")
	}

	re := regexp.MustCompile(`[a-zA-Z]`)
	if !re.MatchString(password) {
		validationErrors = append(validationErrors, "password should contain at least one letter")
	}

	re = regexp.MustCompile(`\d`)
	if !re.MatchString(password) {
		validationErrors = append(validationErrors, "password should contain at least one numeric character")
	}

	re = regexp.MustCompile(`[^\w]`)
	if !re.MatchString(password) {
		validationErrors = append(validationErrors, "password should contain at least one non-alphanumeric character")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return &ValidationError{
		validationErrors,
		fmt.Errorf("validation failed"),
	}
}

func validateConfirmPassword(password, confirmPassword string) error {
	validationErrors := make([]string, 0)
	if password != confirmPassword {
		validationErrors = append(validationErrors, "password and confirm password does not match!")
	}

	if len(password) < 8 {
		validationErrors = append(validationErrors, "password minimum length is 8 characters!")
	}

	re := regexp.MustCompile(`[a-zA-Z]`)
	if !re.MatchString(password) {
		validationErrors = append(validationErrors, "password should contain at least one letter")
	}

	re = regexp.MustCompile(`\d`)
	if !re.MatchString(password) {
		validationErrors = append(validationErrors, "password should contain at least one numeric character")
	}

	re = regexp.MustCompile(`[^\w]`)
	if !re.MatchString(password) {
		validationErrors = append(validationErrors, "password should contain at least one non-alphanumeric character")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return &ValidationError{
		validationErrors,
		fmt.Errorf("validation failed"),
	}
}

func validateEmail(emailAddress string) error {
	validator := validator.New()

	if err := validator.Var(emailAddress, "required,email"); err != nil {
		errors := make([]string, 1)
		errors = append(errors, "invalid email: %s", emailAddress)

		return &ValidationError{
			errors,
			err,
		}
	}

	return nil
}
