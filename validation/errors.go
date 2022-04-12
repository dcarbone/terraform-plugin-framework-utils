package validation

import "errors"

var (
	ErrNoComparisonFuncRegistered = errors.New("no comparison func registered")
	ErrTypeConversionFailed       = errors.New("type conversion failed")
	ErrComparisonFailed           = errors.New("comparison failed")
)
