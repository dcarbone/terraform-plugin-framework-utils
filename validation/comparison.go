package validation

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/dcarbone/terraform-plugin-framework-utils/conv"
	"github.com/dcarbone/terraform-plugin-framework-utils/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type CompareOp string

const (
	Equal                CompareOp = "=="
	LessThan             CompareOp = "<"
	LessThanOrEqualTo    CompareOp = "<="
	GreaterThan          CompareOp = ">"
	GreaterThanOrEqualTo CompareOp = ">="
	NotEqual             CompareOp = "<>"
	OneOf                CompareOp = "|"
	NotOneOf             CompareOp = "^|"
)

func (op CompareOp) String() string {
	return string(op)
}

func (op CompareOp) Name() string {
	switch op {
	case Equal:
		return "equal"
	case LessThan:
		return "less_than"
	case LessThanOrEqualTo:
		return "less_than_or_equal_to"
	case GreaterThan:
		return "greater_than"
	case GreaterThanOrEqualTo:
		return "greater_than_or_equal_to"
	case NotEqual:
		return "not_equal"
	case OneOf:
		return "one_of"
	case NotOneOf:
		return "not_one_of"

	default:
		return "UNKNOWN"
	}
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
type ComparisonFunc func(av attr.Value, op CompareOp, target interface{}, meta ...interface{}) error

var (
	comparisonFuncsMu sync.Mutex
	comparisonFuncs   map[string]ComparisonFunc
)

func compareBool(av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	actBool := conv.BoolValueToBool(av)
	expBool, err := util.TryCoerceToBool(target)
	if err != nil {
		return UnexpectedComparisonTargetTypeError("compare_bool", true, target, err)
	}
	switch op {
	case Equal:
		if actBool == expBool {
			return nil
		}
	case NotEqual:
		if actBool != expBool {
			return nil
		}

	default:
		return NoComparisonFuncRegisteredError(av, op)
	}

	return ComparisonFailedError(op, expBool, actBool)
}

func compareFloat64(av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	actF64, _, err := conv.AttributeValueToFloat64(av)
	if err != nil {
		return TypeConversionFailedError(err)
	}
	expF64, err := util.TryCoerceToFloat64(target)
	if err != nil {
		return UnexpectedComparisonTargetTypeError("compare_float64", float64(0), target, err)
	}

	switch op {
	case Equal:
		if actF64 == expF64 {
			return nil
		}
	case NotEqual:
		if actF64 != expF64 {
			return nil
		}
	case GreaterThan:
		if actF64 > expF64 {
			return nil
		}
	case GreaterThanOrEqualTo:
		if actF64 >= expF64 {
			return nil
		}
	case LessThan:
		if actF64 < expF64 {
			return nil
		}
	case LessThanOrEqualTo:
		if actF64 <= expF64 {
			return nil
		}

	default:
		return NoComparisonFuncRegisteredError(av, op)
	}

	return ComparisonFailedError(op, expF64, actF64)
}

func compareInt64(av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	actI64, _, err := conv.AttributeValueToInt64(av)
	if err != nil {
		return TypeConversionFailedError(err)
	}
	tgtI64, err := util.TryCoerceToInt64(target)
	if err != nil {
		return UnexpectedComparisonTargetTypeError("compare_int64", int64(0), target, err)
	}

	switch op {
	case Equal:
		if actI64 == tgtI64 {
			return nil
		}
	case NotEqual:
		if actI64 != tgtI64 {
			return nil
		}
	case GreaterThan:
		if actI64 > tgtI64 {
			return nil
		}
	case GreaterThanOrEqualTo:
		if actI64 >= tgtI64 {
			return nil
		}
	case LessThan:
		if actI64 < tgtI64 {
			return nil
		}
	case LessThanOrEqualTo:
		if actI64 <= tgtI64 {
			return nil
		}

	default:
		return NoComparisonFuncRegisteredError(av, op)
	}

	return ComparisonFailedError(op, tgtI64, actI64)
}

func compareInt(av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	return compareInt64(av, op, int64(target.(int)))
}

func compareBigFloat(av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	actualBF := conv.NumberValueToBigFloat(av)
	expectedBF, err := util.TryCoerceToBigFloat(target)
	if err != nil {
		return UnexpectedComparisonTargetTypeError("compare_bigfloat", (*big.Float)(nil), target, nil)
	}

	cmp := actualBF.Cmp(expectedBF)

	switch op {
	case Equal:
		if cmp == 0 {
			return nil
		}
	case NotEqual:
		if cmp == 0 {
			exp, _ := expectedBF.Float64()
			act, _ := actualBF.Float64()
			return ComparisonFailedError(op, exp, act)
		}
	case GreaterThan:
		if cmp == 1 {
			return nil
		}
	case GreaterThanOrEqualTo:
		if cmp == 0 || cmp == 1 {
			return nil
		}
	case LessThan:
		if cmp == -1 {
			return nil
		}
	case LessThanOrEqualTo:
		if cmp == -1 || cmp == 0 {
			return nil
		}

	default:
		return NoComparisonFuncRegisteredError(av, op)
	}

	exp, _ := expectedBF.Float64()
	act, _ := actualBF.Float64()
	return ComparisonFailedError(op, exp, act)
}

func compareString(av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	actStr := conv.StringValueToString(av)
	tgtStr, ok := target.(string)
	if !ok {
		return UnexpectedComparisonTargetTypeError("compare_string", "", target, nil)
	}
	switch op {
	case Equal:
		if actStr == tgtStr {
			return nil
		}
	case NotEqual:
		if actStr != tgtStr {
			return nil
		}

	default:
		return NoComparisonFuncRegisteredError(av, op)
	}

	return ComparisonFailedError(op, tgtStr, actStr)
}

func compareStrings(av attr.Value, op CompareOp, target interface{}, meta ...interface{}) error {
	var (
		expStrs         []string
		caseInsensitive bool
	)

	if len(meta) > 0 {
		if b, ok := meta[0].(bool); ok {
			caseInsensitive = b
		}
	}

	actStr := conv.AttributeValueToString(av)
	if caseInsensitive {
		actStr = strings.ToLower(actStr)
	}

	if targ, ok := target.([]interface{}); !ok {
		return UnexpectedComparisonTargetTypeError("compare_strings", make([]string, 0), target, nil)
	} else {
		expStrs = make([]string, len(targ))
		for i, v := range targ {
			if s, ok := v.(string); ok {
				if caseInsensitive {
					expStrs[i] = s
				} else {
					expStrs[i] = strings.ToLower(s)
				}
			} else {
				return UnexpectedComparisonTargetTypeError("compare_strings", make([]string, 0), v, nil)
			}
		}
	}

	switch op {
	case OneOf:
		for _, v := range expStrs {
			if actStr == v {
				return nil
			}
		}

	case NotOneOf:
		for _, v := range expStrs {
			if actStr == v {
				return ComparisonFailedError(op, actStr, expStrs)
			}
		}
		return nil

	default:
		return NoComparisonFuncRegisteredError(av, op)
	}

	return ComparisonFailedError(op, actStr, expStrs)
}

// DefaultComparisonFuncs returns the complete list of default comparison functions
func DefaultComparisonFuncs() map[string]ComparisonFunc {
	return map[string]ComparisonFunc{
		util.KeyFN(false):             compareBool,
		util.KeyFN(0.0):               compareFloat64,
		util.KeyFN(int64(0)):          compareInt64,
		util.KeyFN(0):                 compareInt,
		util.KeyFN((*big.Float)(nil)): compareBigFloat,
		util.KeyFN(""):                compareString,
		util.KeyFN(make([]string, 0)): compareStrings,
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
		return fmt.Errorf("%w for operation %q with target type %T", ErrNoComparisonFuncRegistered, op, target)
	}
}

func addComparisonFailedDiagnostic(op CompareOp, target interface{}, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse, err error) {
	switch op {
	case Equal:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Attribute value does not match expected",
			fmt.Sprintf("Attribute value must equal %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case NotEqual:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Attribute value is not allowed",
			fmt.Sprintf("Attribute value must not equal %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case LessThan:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value is above threshold",
			fmt.Sprintf("Attribute value must be less than %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case LessThanOrEqualTo:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value is above threshold",
			fmt.Sprintf("Attribute value must be less than or equal to %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case GreaterThan:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value is below threshold",
			fmt.Sprintf("Attribute value must be greater than %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case GreaterThanOrEqualTo:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value is below threshold",
			fmt.Sprintf("Attribute value must be greater than or equal to %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case OneOf:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value is not within allowed list",
			fmt.Sprintf("Attribute value must be one of %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)

	case NotOneOf:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Value is not within allowed list",
			fmt.Sprintf("Attribute value must not be one of %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)

	default:
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Unknown comparison operation",
			fmt.Sprintf("Specified unknown comparison operation: %s", op),
		)
	}
}
