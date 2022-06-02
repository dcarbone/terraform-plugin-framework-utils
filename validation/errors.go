package validation

import (
	"errors"
	"fmt"

	"github.com/dcarbone/terraform-plugin-framework-utils/internal/util"
)

var (
	ErrNoComparisonFuncRegistered     = errors.New("no comparison func registered")
	ErrTypeConversionFailed           = errors.New("type conversion failed")
	ErrComparisonFailed               = errors.New("comparison failed")
	ErrUnexpectedComparisonTargetType = errors.New("unexpected comparison target value type")
)

func NoComparisonFuncRegisteredError(t interface{}, op CompareOp) error {
	return fmt.Errorf("%w: type=%T; op=%q", ErrNoComparisonFuncRegistered, t, op.Name())
}

func IsNoComparisonFuncRegisteredError(err error) bool {
	return util.MatchError(err, ErrNoComparisonFuncRegistered)
}

func TypeConversionFailedError(err error) error {
	return fmt.Errorf("%w: %v", ErrTypeConversionFailed, err)
}

func IsTypeConversionFailedError(err error) bool {
	return util.MatchError(err, ErrTypeConversionFailed)
}

func ComparisonFailedError(op CompareOp, target, actual interface{}) error {
	return fmt.Errorf("%[1]w: op=%[2]q; target=%[3]T(%[3]v); actual=%[4]T(%[4]v)", ErrComparisonFailed, op.Name(), target, actual)
}

func IsComparisonFailedError(err error) bool {
	return util.MatchError(err, ErrComparisonFailed)
}

func UnexpectedComparisonTargetTypeError(scope string, expected, actual interface{}, err error) error {
	if err != nil {
		return fmt.Errorf("%w: scope=%q; actual=%T; err=%v", ErrUnexpectedComparisonTargetType, scope, actual, err)
	}
	if expected == nil {
		return fmt.Errorf("%w: scope=%q; type=%T", ErrUnexpectedComparisonTargetType, scope, actual)
	}
	return fmt.Errorf("%w: scope=%q; expected=%T; actual=%T", ErrUnexpectedComparisonTargetType, scope, expected, actual)
}

func IsUnexpectedAttributeValueTypeError(err error) bool {
	return util.MatchError(err, ErrUnexpectedComparisonTargetType)
}
