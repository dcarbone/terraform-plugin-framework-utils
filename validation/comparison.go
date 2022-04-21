package validation

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/dcarbone/terraform-plugin-framework-utils/conv"
	"github.com/dcarbone/terraform-plugin-framework-utils/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CompareOp string

const (
	Equal                CompareOp = "=="
	LessThan             CompareOp = "<"
	LessThanOrEqualTo    CompareOp = "<="
	GreaterThan          CompareOp = ">"
	GreaterThanOrEqualTo CompareOp = ">="
	NotEqual             CompareOp = "!="
)

func (op CompareOp) String() string {
	return string(op)
}

// ComparisonFunc executes a specific comparison of an attribute value to the targeted value.  You are guaranteed that
// the target type will be the type or one of the types the function was registered with.  If you register a single func
// with more than type target type, you must perform type assertion / conversion yourself.
//
// The returned error is expected to be testable for type:
//
//		nil 					- Comparison succeeded
//		ErrComparisonFailed 	- Must be returned when any comparison operation fails
//		ErrTypeConversionFailed - Must be returned if the function performs an internal type conversion before comparison that errored
//		any other error			- Treated as unhandled error
//
// To see the default list of functions, see DefaultComparisonFuncs.
//
// To register a new function or overwrite an existing function, see SetComparisonFunc
type ComparisonFunc func(av attr.Value, op CompareOp, target interface{}) error

var (
	comparisonFuncsMu sync.Mutex
	comparisonFuncs   map[string]ComparisonFunc
)

func compareBool(av attr.Value, op CompareOp, target interface{}) error {
	switch op {
	case Equal:
		if !av.Equal(types.Bool{Value: target.(bool)}) {
			return ErrComparisonFailed
		}
	case NotEqual:
		if av.Equal(types.Bool{Value: target.(bool)}) {
			return ErrComparisonFailed
		}

	default:
		return ErrNoComparisonFuncRegistered
	}
	return nil
}

func compareFloat64(av attr.Value, op CompareOp, target interface{}) error {
	asF64, _, err := conv.AttributeValueToFloat64(av)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrTypeConversionFailed, err)
	}
	switch op {
	case Equal:
		if asF64 != target.(float64) {
			return ErrComparisonFailed
		}
	case NotEqual:
		if asF64 == target.(float64) {
			return ErrComparisonFailed
		}
	case GreaterThan:
		if asF64 <= target.(float64) {
			return ErrComparisonFailed
		}
	case GreaterThanOrEqualTo:
		if asF64 < target.(float64) {
			return ErrComparisonFailed
		}
	case LessThan:
		if asF64 >= target.(float64) {
			return ErrComparisonFailed
		}
	case LessThanOrEqualTo:
		if asF64 > target.(float64) {
			return ErrComparisonFailed
		}

	default:
		panic(fmt.Sprintf("unknown comparison operator: %q", op))
	}

	return nil
}

func compareInt64(av attr.Value, op CompareOp, target interface{}) error {
	asI64, _, err := conv.AttributeValueToInt64(av)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrTypeConversionFailed, err)
	}
	switch op {
	case Equal:
		if asI64 != target.(int64) {
			return ErrComparisonFailed
		}
	case NotEqual:
		if asI64 == target.(int64) {
			return ErrComparisonFailed
		}
	case GreaterThan:
		if asI64 <= target.(int64) {
			return ErrComparisonFailed
		}
	case GreaterThanOrEqualTo:
		if asI64 < target.(int64) {
			return ErrComparisonFailed
		}
	case LessThan:
		if asI64 >= target.(int64) {
			return ErrComparisonFailed
		}
	case LessThanOrEqualTo:
		if asI64 > target.(int64) {
			return ErrComparisonFailed
		}

	default:
		panic(fmt.Sprintf("unknown comparison operator: %q", op))
	}

	return nil
}

func compareInt(av attr.Value, op CompareOp, target interface{}) error {
	return compareInt64(av, op, int64(target.(int)))
}

func compareBigFloat(av attr.Value, op CompareOp, target interface{}) error {
	asBF, err := conv.AttributeValueToBigFloat(av)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrTypeConversionFailed, err)
	}
	cmp := asBF.Cmp(target.(*big.Float))

	switch op {
	case Equal:
		if cmp != 0 {
			return ErrComparisonFailed
		}
	case NotEqual:
		if cmp == 0 {
			return ErrComparisonFailed
		}
	case GreaterThan:
		if cmp != 1 {
			return ErrComparisonFailed
		}
	case GreaterThanOrEqualTo:
		if cmp == -1 {
			return ErrComparisonFailed
		}
	case LessThan:
		if cmp != 1 {
			return ErrComparisonFailed
		}
	case LessThanOrEqualTo:
		if cmp == 1 {
			return ErrComparisonFailed
		}

	default:
		return ErrNoComparisonFuncRegistered
	}

	return nil
}

func compareString(av attr.Value, op CompareOp, target interface{}) error {
	asStr := conv.AttributeValueToString(av)
	switch op {
	case Equal:
		if asStr != target.(string) {
			return ErrComparisonFailed
		}
	case NotEqual:
		if asStr == target.(string) {
			return ErrComparisonFailed
		}

	default:
		return ErrNoComparisonFuncRegistered
	}

	return nil
}

// DefaultComparisonFuncs returns the complete list of default comparison functions
func DefaultComparisonFuncs() map[string]ComparisonFunc {
	return map[string]ComparisonFunc{
		util.KeyFN(false):          compareBool,
		util.KeyFN(0.0):            compareFloat64,
		util.KeyFN(int64(0)):       compareInt64,
		util.KeyFN(0):              compareInt,
		util.KeyFN(new(big.Float)): compareBigFloat,
		util.KeyFN(""):             compareString,
	}
}

// SetComparisonFunc sets a comparison function to use for comparing attribute values to values of the specified type
func SetComparisonFunc(targetType interface{}, fn ComparisonFunc) {
	comparisonFuncsMu.Lock()
	defer comparisonFuncsMu.Unlock()
	comparisonFuncs[util.KeyFN(targetType)] = fn
}

// GetComparisonFunc attempts to return a previously registered comparison function for a specified op : type
// combination
func GetComparisonFunc(targetType interface{}) (ComparisonFunc, bool) {
	comparisonFuncsMu.Lock()
	defer comparisonFuncsMu.Unlock()
	if fn, ok := comparisonFuncs[util.KeyFN(targetType)]; ok {
		return fn, true
	}
	return nil, false
}

func init() {
	comparisonFuncs = DefaultComparisonFuncs()
}

// CompareAttrValues attempts to execute a comparison between the provided attribute value and the targeted value.
//
// If there is no comparison function registered for the target type, an ErrNoComparisonFuncRegistered
// is returned.
//
// If a function is registered and the comparison fails, an ErrComparisonFailed error will be returned
func CompareAttrValues(av attr.Value, op CompareOp, target interface{}) error {
	if fn, ok := GetComparisonFunc(target); ok {
		return fn(av, op, target)
	} else {
		return fmt.Errorf("%w for operation %q with expected type %T", ErrNoComparisonFuncRegistered, op, target)
	}
}

func addComparisonFailedDiagnostic(op CompareOp, expected interface{}, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	switch op {
	case Equal:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value does not match expected",
			fmt.Sprintf("Value must not be less than %q", conv.GoNumberToString(expected)),
		)
	case NotEqual:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			fmt.Sprintf("Value must not equal %q", conv.GoNumberToString(expected)),
			fmt.Sprintf("Value must not equal %q", conv.GoNumberToString(expected)),
		)
	case LessThan:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value is above threshold",
			fmt.Sprintf("Value must be less than %q", conv.GoNumberToString(expected)),
		)
	case LessThanOrEqualTo:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value is above threshold",
			fmt.Sprintf("Value must be less than or equal to %q", conv.GoNumberToString(expected)),
		)
	case GreaterThan:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value is below threshold",
			fmt.Sprintf("Value must be greater than %q", conv.GoNumberToString(expected)),
		)
	case GreaterThanOrEqualTo:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value is below threshold",
			fmt.Sprintf("Value must be greater than or equal to %q", conv.GoNumberToString(expected)),
		)

	default:
		panic(fmt.Sprintf("no diagnostic message handler defined for op %q for attribute %q", op, req.AttributePath.String()))
	}
}
