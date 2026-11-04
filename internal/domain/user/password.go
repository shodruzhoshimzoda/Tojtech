package user

import (
	"fmt"
	"strings"
)

const (
	PasswordMinLength = 8
	PasswordMaxLength = 16
)

// PasswordPolicyError - the password does not follow the rules; Problems says what exactly is wrong.
type PasswordPolicyError struct {
	Problems []string
}

func (e *PasswordPolicyError) Error() string {
	return "password must " + joinWithAnd(e.Problems)
}

func ValidatePassword(password string) error {
	var (
		length                       int
		hasLower, hasUpper, hasDigit bool
		hasForbidden                 bool
	)

	for _, r := range password { // range goes over characters, not bytes
		length++
		switch {
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= '0' && r <= '9':
			hasDigit = true
		default:
			hasForbidden = true
		}
	}

	var problems []string
	if length < PasswordMinLength || length > PasswordMaxLength {
		problems = append(problems, fmt.Sprintf("be %d-%d characters long", PasswordMinLength, PasswordMaxLength))
	}
	if hasForbidden {
		problems = append(problems, "contain only Latin letters and digits")
	}
	if !hasLower {
		problems = append(problems, "contain a lowercase letter")
	}
	if !hasUpper {
		problems = append(problems, "contain an uppercase letter")
	}
	if !hasDigit {
		problems = append(problems, "contain a digit")
	}

	if len(problems) > 0 {
		return &PasswordPolicyError{Problems: problems}
	}
	return nil
}

func joinWithAnd(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	default:
		return strings.Join(items[:len(items)-1], ", ") + " and " + items[len(items)-1]
	}
}