package conv

import (
	"errors"
	"fmt"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

var (
	ErrValueIsNull        = errors.New("value is nil")
	ErrValueIsUnknown     = errors.New("value is unknown")
	ErrValueIsEmpty       = errors.New("value is empty")
	ErrValueTypeUnhandled = errors.New("value type is unhandled, this usually means this package is out of date with the upstream provider framework")
)

func IsValueIsNullError(err error) bool {
	return util.MatchError(err, ErrValueIsNull)
}

func IsValueIsUnknownError(err error) bool {
	return util.MatchError(err, ErrValueIsUnknown)
}

func IsValueIsEmptyError(err error) bool {
	return util.MatchError(err, ErrValueIsEmpty)
}

func ValueTypeUnhandledError(scope string, av attr.Value) error {
	return fmt.Errorf("%w: scope=%q; type=%T", ErrValueTypeUnhandled, scope, av)
}

func IsValueTypeUnhandledError(err error) bool {
	return util.MatchError(err, ErrValueTypeUnhandled)
}
