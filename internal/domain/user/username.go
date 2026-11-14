package user

import "errors"

const (

	// The username lenght
	UsernameMinLength = 3
	UsernameMaxLength = 20
)

var (
	ErrUsernameRequired = errors.New("username is required")
	ErrUsernameTooShort = errors.New("username must be at least 3 characters long")
	ErrUsernameTooLong  = errors.New("username must be at most 20 characters long")
	ErrUsernameChars    = errors.New("username may contain only Latin letters, digits and underscore")
	ErrUsernameTaken    = errors.New("username is already taken")
)

// ValidateUsername checks a username that was already trimmed.
func ValidateUsername(username string) error {
	if username == "" {
		return ErrUsernameRequired
	}

	for _, r := range username { // range goes over characters, not bytes
		isLetter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		isDigit := r >= '0' && r <= '9'
		if !isLetter && !isDigit && r != '_' {
			return ErrUsernameChars
		}
	}

	// here every character is ASCII, so len() equals the number of characters
	if len(username) < UsernameMinLength {
		return ErrUsernameTooShort
	}
	if len(username) > UsernameMaxLength {
		return ErrUsernameTooLong
	}

	return nil
}

// IsUsernameError reports whether err is a "bad username" error (HTTP 400).
// ErrUsernameTaken is not included: it is a different case (HTTP 409).
func IsUsernameError(err error) bool {
	return errors.Is(err, ErrUsernameRequired) ||
		errors.Is(err, ErrUsernameTooShort) ||
		errors.Is(err, ErrUsernameTooLong) ||
		errors.Is(err, ErrUsernameChars)
}