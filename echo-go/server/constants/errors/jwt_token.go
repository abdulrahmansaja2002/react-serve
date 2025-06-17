package errors

import "errors"

var (
	ErrInvalidToken        = errors.New("missing or invalid JWT token")
	ErrFailedToParseClaims = errors.New("failed to parse JWT claims")
)
