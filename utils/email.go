package utils

import (
	"errors"
	"regexp"
)

// VerifyEmailFormat выполняет проверку email
func VerifyEmailFormat(email string) error {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !emailRegexp.MatchString(email) {
		return errors.New("Invalid email format")
	}

	return nil
}
