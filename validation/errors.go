package validation

import (
	"errors"
	"fmt"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/internal/util"
)

var (
	ErrNoComparisonFuncRegistered     = errors.New("no comparison func registered")
	ErrTypeConversionFailed           = errors.New("type conversion failed")
	ErrComparisonFailed               = errors.New("comparison failed")
	ErrUnexpectedComparisonTargetType = errors.New("unexpected comparison \"target\" value type")
	ErrUnexpectedComparisonActualType = errors.New("unexpected comparison \"actual\" value type")
)

func NoComparisonFuncRegisteredError(op CompareOp, t interface{}) error {
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

func ComparisonFailedError(actual interface{}, op CompareOp, target interface{}) error {
	return fmt.Errorf("%[1]w: op=%[2]q; target=%[3]T(%[3]v); actual=%[4]T(%[4]v)", ErrComparisonFailed, op.Name(), target, actual)
}

func IsComparisonFailedError(err error) bool {
	return util.MatchError(err, ErrComparisonFailed)
}

func UnexpectedComparisonTargetTypeError(scope string, actualTarget interface{}, op CompareOp, expectedTarget interface{}, err error) error {
	if err != nil {
		return fmt.Errorf("%w: scope=%q; op=%q; actual=%T; expected=%T, err=%v", ErrUnexpectedComparisonTargetType, scope, op.Name(), actualTarget, expectedTarget, err)
	}
	return fmt.Errorf("%w: scope=%q; op=%q; actual=%T; expected=%T", ErrUnexpectedComparisonTargetType, scope, op.Name(), actualTarget, expectedTarget)
}

func IsUnexpectedAttributeValueTypeError(err error) bool {
	return util.MatchError(err, ErrUnexpectedComparisonTargetType)
}

func UnexpectedComparisonActualTypeError(scope string, actualActual interface{}, op CompareOp, expectedActual interface{}, err error) error {
	// TODO: enable more detailed type printing, particularly for attr.Value representations
	if err != nil {
		return fmt.Errorf("%w: scope=%q; op=%q; actual=%T; expected=%T; err=%v", ErrUnexpectedComparisonActualType, scope, op.Name(), actualActual, expectedActual, err)
	}
	return fmt.Errorf("%w: scope=%q; op=%q; expected=%T; actual=%T", ErrUnexpectedComparisonActualType, scope, op.Name(), expectedActual, actualActual)
}

func IsUnexpectedComparisonActualType(err error) bool {
	return util.MatchError(err, ErrUnexpectedComparisonActualType)
}
