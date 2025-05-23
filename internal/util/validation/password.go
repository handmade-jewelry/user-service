package validation

import (
	"unicode"
)

const passwordLength = 18

func ValidatePassword(p string) bool {
	if len(p) != passwordLength {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range p {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
