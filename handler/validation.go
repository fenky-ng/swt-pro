package handler

import (
	"context"
	"regexp"
	"strings"

	"github.com/fenky-ng/swt-pro/generated"
)

func validateRegistration(request generated.RegistrationRequest) []string {
	var errorMessages []string

	// validate phone number
	errorMessages = append(errorMessages, validatePhoneNumber(request.PhoneNumber)...)

	// validate full name
	errorMessages = append(errorMessages, validateFullName(request.FullName)...)

	// validate password
	if !validatePassword(request.Password) {
		errorMessages = append(errorMessages, "Passwords must be minimum 6 characters and maximum 64 characters, containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters")
	}

	return errorMessages
}

func validatePhoneNumber(input string) []string {
	var errorMessages []string

	if len(input) < 10 || len(input) > 13 {
		errorMessages = append(errorMessages, "Phone numbers must be at minimum 10 characters and maximum 13 characters")
	}
	if !strings.HasPrefix(input, "+62") {
		errorMessages = append(errorMessages, "Phone numbers must start with the Indonesia country code “+62”")
	}

	return errorMessages
}

func validateFullName(input string) []string {
	var errorMessages []string

	if len(input) < 3 || len(input) > 60 {
		errorMessages = append(errorMessages, "Full name must be at minimum 3 characters and maximum 60 characters")
	}

	return errorMessages
}

func validatePassword(input string) bool {
	validations := []string{
		".{6,64}",   // length of password
		"[A-Z]",     // at least capital character
		"[0-9]",     // at least one numeric character
		"[^\\d\\w]", // at least one non alpha-numeric character
	}
	for _, validation := range validations {
		match, _ := regexp.MatchString(validation, input)
		if !match {
			return false
		}
	}

	return true
}

func checkNewPhoneNumber(
	ctx context.Context,
	s *Server,
	phoneNumber string,
) (bool, error) {
	var (
		res bool
		err error
	)

	user, err := s.Repository.GetUserByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return res, err
	}

	if user.ID == 0 {
		res = true
	}

	return res, err
}
