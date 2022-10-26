package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	validUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	validFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func validateStringLen(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(value string) error {
	if err := validateStringLen(value, 3, 100); err != nil {
		return err
	}

	if valid := validUsername(value); !valid {
		return fmt.Errorf("must contain only lowercase letters, numbers and underscores")
	}

	return nil
}

func ValidatePassword(value string) error {
	if err := validateStringLen(value, 6, 100); err != nil {
		return err
	}

	return nil
}

func ValidateEmail(value string) error {
	if err := validateStringLen(value, 3, 300); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}

	return nil
}

func ValidateFullName(value string) error {
	if err := validateStringLen(value, 3, 100); err != nil {
		return err
	}

	if valid := validFullname(value); !valid {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}
