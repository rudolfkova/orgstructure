// Package validation ...
package validation

import (
	"errors"
)

var (
	// ErrValidationStrLen ...
	ErrValidationStrLen = errors.New("invalid length of string")
)

// ValidateStr ...
func ValidateStr(s string, maxLen int) error {

	if len(s) <= maxLen && len(s) >= 1 {
		return nil
	}
	return ErrValidationStrLen
}
