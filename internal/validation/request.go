// Package validation ...
package validation

import orgerror "orgstructure/internal/errors"

// ValidateStr ...
func ValidateStr(s string, maxLen int) error {

	if len(s) <= maxLen && len(s) >= 1 {
		return nil
	}
	return orgerror.ErrValidationStrLen
}
