package validation

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/conv"
	"github.com/dcarbone/terraform-plugin-framework-utils/v3/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		return string(op)
	}
}

// ComparisonFunc executes a specific comparison of an attribute value to the targeted value.  You are guaranteed that
// the target type will be the type or one of the types the function was registered with.  If you register a single func
// with more than type target type, you must perform type assertion / conversion yourself.
//
// The returned error is expected to be testable for type:
//
//	nil 					- Comparison succeeded
//	ErrComparisonFailed 	- Must be returned when any comparison operation fails
//	ErrTypeConversionFailed - Must be returned if the function performs an internal type conversion before comparison that errored
//	any other error			- Treated as unhandled error
//
// To see the default list of functions, see DefaultComparisonFuncs.
//
// To register a new function or overwrite an existing function, see SetComparisonFunc
type ComparisonFunc func(ctx context.Context, av attr.Value, op CompareOp, target interface{}, meta ...interface{}) error

var (
	comparisonFuncsMu sync.Mutex
	comparisonFuncs   map[string]ComparisonFunc
)

func compareBool(ctx context.Context, av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	actBool := conv.BoolValueToBool(av)
	expBool, err := util.TryCoerceToBool(target)
	if err != nil {
		return UnexpectedComparisonTargetTypeError("compare_bool", target, op, true, err)
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
		return NoComparisonFuncRegisteredError(op, av)
	}

	return ComparisonFailedError(actBool, op, expBool)
}

func compareFloat64(_ context.Context, av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	actF64, _, err := conv.AttributeValueToFloat64(av)
	if err != nil {
		return TypeConversionFailedError(err)
	}
	expF64, err := util.TryCoerceToFloat64(target)
	if err != nil {
		return UnexpectedComparisonTargetTypeError("compare_float64", target, op, float64(0), err)
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
		return NoComparisonFuncRegisteredError(op, av)
	}

	return ComparisonFailedError(actF64, op, expF64)
}

func compareInt64(_ context.Context, av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	actI64, _, err := conv.AttributeValueToInt64(av)
	if err != nil {
		return TypeConversionFailedError(err)
	}
	tgtI64, err := util.TryCoerceToInt64(target)
	if err != nil {
		return UnexpectedComparisonTargetTypeError("compare_int64", target, op, int64(0), err)
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
		return NoComparisonFuncRegisteredError(op, av)
	}

	return ComparisonFailedError(actI64, op, tgtI64)
}

func compareInt(ctx context.Context, av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	return compareInt64(ctx, av, op, int64(target.(int)))
}

func compareBigFloat(_ context.Context, av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	actualBF := conv.NumberValueToBigFloat(av)
	expectedBF, err := util.TryCoerceToBigFloat(target)
	if err != nil {
		return UnexpectedComparisonTargetTypeError("compare_bigfloat", target, op, (*big.Float)(nil), nil)
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
			return ComparisonFailedError(act, op, exp)
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
		return NoComparisonFuncRegisteredError(op, av)
	}

	exp, _ := expectedBF.Float64()
	act, _ := actualBF.Float64()
	return ComparisonFailedError(act, op, exp)
}

func compareString(_ context.Context, av attr.Value, op CompareOp, target interface{}, meta ...interface{}) error {
	var caseInsensitive bool
	if len(meta) > 0 {
		if b, ok := meta[0].(bool); ok {
			caseInsensitive = b
		}
	}
	actStr := conv.AttributeValueToString(av)
	tgtStr, ok := target.(string)
	if !ok {
		return UnexpectedComparisonTargetTypeError("compare_string", target, op, "", nil)
	}
	if caseInsensitive {
		actStr = strings.ToLower(actStr)
		tgtStr = strings.ToLower(tgtStr)
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
		return NoComparisonFuncRegisteredError(op, av)
	}

	return ComparisonFailedError(actStr, op, tgtStr)
}

func compareStringToStrings(av types.String, op CompareOp, targets []string, caseInsensitive bool) error {
	var actStr string
	if caseInsensitive {
		actStr = strings.ToLower(av.ValueString())
		for i, v := range targets {
			targets[i] = strings.ToLower(v)
		}
	} else {
		actStr = av.ValueString()
	}

	switch op {
	case OneOf:
		for _, v := range targets {
			if actStr == v {
				return nil
			}
		}

	case NotOneOf:
		for _, v := range targets {
			if actStr == v {
				return ComparisonFailedError(targets, op, actStr)
			}
		}
		return nil

	default:
		return NoComparisonFuncRegisteredError(op, targets)
	}

	return ComparisonFailedError(av.ValueString(), op, targets)
}

func compareStringsToStrings(actuals []string, op CompareOp, targets []string, caseInsensitive bool) error {
	if caseInsensitive {
		for i, v := range targets {
			targets[i] = strings.ToLower(v)
		}
		for i, v := range actuals {
			actuals[i] = strings.ToLower(v)
		}
	}

	actualsLen := len(actuals)
	targetsLen := len(targets)

	switch op {
	case Equal:
		if actualsLen == targetsLen {
			for i, v := range actuals {
				if targets[i] != v {
					return ComparisonFailedError(actuals[i], op, targets[i])
				}
			}
			return nil
		}

	case NotEqual:
		if actualsLen != targetsLen {
			return nil
		}

		for i, v := range actuals {
			if targets[i] != v {
				return nil
			}
		}
	default:
		return NoComparisonFuncRegisteredError(op, make([]string, 0))
	}

	return ComparisonFailedError(actuals, op, targets)
}

func compareListToStrings(ctx context.Context, av types.List, op CompareOp, targets []string, caseInsensitive bool) error {
	if av.ElementType(ctx) != types.StringType {
		return UnexpectedComparisonActualTypeError("compare_list_strings", av.ElementType(ctx), op, types.StringType, nil)
	}
	return compareStringsToStrings(conv.StringListToStrings(av), op, targets, caseInsensitive)
}

func compareSetToStrings(ctx context.Context, av types.Set, op CompareOp, targets []string, caseInsensitive bool) error {
	if av.ElementType(ctx) != types.StringType {
		return UnexpectedComparisonActualTypeError("compare_set_strings", av.ElementType(ctx), op, types.StringType, nil)
	}
	return compareStringsToStrings(conv.StringSetToStrings(av), op, targets, caseInsensitive)
}

func compareStrings(ctx context.Context, av attr.Value, op CompareOp, target interface{}, meta ...interface{}) error {
	caseInsensitive := false
	if len(meta) > 0 {
		if b, ok := meta[0].(bool); ok {
			caseInsensitive = b
		}
	}

	tgtStrs, ok := target.([]string)
	if !ok {
		return UnexpectedComparisonTargetTypeError("compare_strings", target, op, make([]string, 0), nil)
	}

	switch av.(type) {
	case types.String, *types.String:
		return compareStringToStrings(conv.ValueToStringType(av), op, tgtStrs, caseInsensitive)

	case types.List, *types.List:
		return compareListToStrings(ctx, conv.ValueToListType(av), op, tgtStrs, caseInsensitive)
	case types.Set, *types.Set:
		return compareSetToStrings(ctx, conv.ValueToSetType(av), op, tgtStrs, caseInsensitive)

	default:
		return UnexpectedComparisonActualTypeError("compare_strings", av, op, types.StringType, nil)
	}
}

func compareInt64ToInts(_ context.Context, av types.Int64, op CompareOp, targets []int, _ ...interface{}) error {
	asInt := int(av.ValueInt64())
	switch op {
	case OneOf:
		for _, v := range targets {
			if asInt == v {
				return nil
			}
		}

	case NotOneOf:
		for _, v := range targets {
			if asInt == v {
				return ComparisonFailedError(targets, op, asInt)
			}
		}
		return nil

	default:
		return NoComparisonFuncRegisteredError(op, targets)
	}

	return ComparisonFailedError(av.ValueInt64(), op, targets)
}

func compareNumberToInts(_ context.Context, av types.Number, op CompareOp, targets []int, _ ...interface{}) error {
	if av.IsNull() {
		return ComparisonFailedError(nil, op, targets)
	}
	asInt64, _ := av.ValueBigFloat().Int64()
	asInt := int(asInt64)
	switch op {
	case OneOf:
		for _, v := range targets {
			if asInt == v {
				return nil
			}
		}

	case NotOneOf:
		for _, v := range targets {
			if asInt == v {
				return ComparisonFailedError(targets, op, asInt)
			}
		}
		return nil

	default:
		return NoComparisonFuncRegisteredError(op, targets)
	}

	v, _ := av.ValueBigFloat().Float64()
	return ComparisonFailedError(v, op, targets)
}

func compareIntsToInts(actuals []int, op CompareOp, targets []int) error {
	actualsLen := len(actuals)
	targetsLen := len(targets)

	switch op {
	case Equal:
		if actualsLen == targetsLen {
			for i, v := range actuals {
				if targets[i] != v {
					return ComparisonFailedError(actuals[i], op, targets[i])
				}
			}
			return nil
		}

	case NotEqual:
		if actualsLen != targetsLen {
			return nil
		}

		for i, v := range actuals {
			if targets[i] != v {
				return nil
			}
		}
	default:
		return NoComparisonFuncRegisteredError(op, make([]int, 0))
	}

	return ComparisonFailedError(actuals, op, targets)
}

func compareListToInts(ctx context.Context, av types.List, op CompareOp, targets []int, _ ...interface{}) error {
	elemType := av.ElementType(ctx)
	switch elemType {
	case types.Int64Type:
		return compareIntsToInts(conv.Int64ListToInts(av), op, targets)
	case types.NumberType:
		return compareIntsToInts(conv.NumberListToInts(av), op, targets)

	default:
		return UnexpectedComparisonActualTypeError("compare_ints", elemType, op, types.Int64Type, nil)
	}
}

func compareSetToInts(ctx context.Context, av types.Set, op CompareOp, targets []int, _ ...interface{}) error {
	elemType := av.ElementType(ctx)
	switch elemType {
	case types.Int64Type:
		return compareIntsToInts(conv.Int64SetToInts(av), op, targets)
	case types.NumberType:
		return compareIntsToInts(conv.NumberSetToInts(av), op, targets)

	default:
		return UnexpectedComparisonActualTypeError("compare_ints", elemType, op, types.Int64Type, nil)
	}
}

func compareInts(ctx context.Context, av attr.Value, op CompareOp, target interface{}, _ ...interface{}) error {
	tgtInts, ok := target.([]int)
	if !ok {
		return UnexpectedComparisonTargetTypeError("compare_ints", target, op, make([]int, 0), nil)
	}

	switch av.(type) {
	case types.Int64, *types.Int64:
		return compareInt64ToInts(ctx, conv.ValueToInt64Type(av), op, tgtInts)
	case types.Number, *types.Number:
		return compareNumberToInts(ctx, conv.ValueToNumberType(av), op, tgtInts)

	case types.List, *types.List:
		return compareListToInts(ctx, conv.ValueToListType(av), op, tgtInts)
	case types.Set, *types.Set:
		return compareSetToInts(ctx, conv.ValueToSetType(av), op, tgtInts)

	default:
		return UnexpectedComparisonActualTypeError("compare_ints", av, op, types.Int64{}, nil)
	}
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
		util.KeyFN(make([]int, 0)):    compareInts,
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
func CompareAttrValues(ctx context.Context, av attr.Value, op CompareOp, target interface{}, meta ...interface{}) error {
	if fn, ok := GetComparisonFunc(target); ok {
		return fn(ctx, av, op, target, meta...)
	} else {
		return fmt.Errorf("%w for operation %q with target type %T", ErrNoComparisonFuncRegistered, op, target)
	}
}

func addComparisonFailedDiagnostic(op CompareOp, target interface{}, srcReq interface{}, srcResp interface{}, err error) {
	var (
		req  GenericRequest
		resp *GenericResponse
		terr error
	)

	if req, resp, terr = toGenericTypes(srcReq, srcResp); terr != nil {
		panic(terr.Error())
	}

	switch op {
	case Equal:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Attribute value does not match expected",
			fmt.Sprintf("Attribute value must equal %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case NotEqual:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Attribute value is not allowed",
			fmt.Sprintf("Attribute value must not equal %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case LessThan:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value is above threshold",
			fmt.Sprintf("Attribute value must be less than %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case LessThanOrEqualTo:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value is above threshold",
			fmt.Sprintf("Attribute value must be less than or equal to %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case GreaterThan:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value is below threshold",
			fmt.Sprintf("Attribute value must be greater than %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case GreaterThanOrEqualTo:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value is below threshold",
			fmt.Sprintf("Attribute value must be greater than or equal to %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)
	case OneOf:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value is not within allowed list",
			fmt.Sprintf("Attribute value must be one of %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)

	case NotOneOf:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value is not within allowed list",
			fmt.Sprintf("Attribute value must not be one of %s; err=%v", util.GetPrintableTypeWithValue(target), err),
		)

	default:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Unknown comparison operation",
			fmt.Sprintf("Specified unknown comparison operation: %s", op),
		)
	}
}
