package conv

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ValueToBoolType ensures we have a types.Bool literal
func ValueToBoolType(v attr.Value) types.Bool {
	if vb, ok := v.(types.Bool); ok {
		return vb
	} else if vb, ok := v.(*types.Bool); ok {
		return *vb
	} else {
		panic(fmt.Sprintf("cannot pass type %T to conv.ValueToBoolType", v))
	}
}

// ValueToFloat64Type ensures we have a types.Float64 literal
func ValueToFloat64Type(v attr.Value) types.Float64 {
	if vb, ok := v.(types.Float64); ok {
		return vb
	} else if vb, ok := v.(*types.Float64); ok {
		return *vb
	} else {
		panic(fmt.Sprintf("cannot pass type %T to conv.ValueToFloat64Type", v))
	}
}

// ValueToInt64Type ensures we have a types.Int64 literal
func ValueToInt64Type(v attr.Value) types.Int64 {
	if vb, ok := v.(types.Int64); ok {
		return vb
	} else if vb, ok := v.(*types.Int64); ok {
		return *vb
	} else {
		panic(fmt.Sprintf("cannot pass type %T to conv.ValueToInt64Type", v))
	}
}

// ValueToListType ensures we have a types.List literal
func ValueToListType(v attr.Value) types.List {
	if vb, ok := v.(types.List); ok {
		return vb
	} else if vb, ok := v.(*types.List); ok {
		return *vb
	} else {
		panic(fmt.Sprintf("cannot pass type %T to conv.ValueToListType", v))
	}
}

// ValueToMapType ensures we have a types.Map literal
func ValueToMapType(v attr.Value) types.Map {
	if vb, ok := v.(types.Map); ok {
		return vb
	} else if vb, ok := v.(*types.Map); ok {
		return *vb
	} else {
		panic(fmt.Sprintf("cannot pass type %T to conv.ValueToMapType", v))
	}
}

// ValueToNumberType ensures we have a types.Number literal
func ValueToNumberType(v attr.Value) types.Number {
	if vb, ok := v.(types.Number); ok {
		return vb
	} else if vb, ok := v.(*types.Number); ok {
		return *vb
	} else {
		panic(fmt.Sprintf("cannot pass type %T to conv.ValueToNumberType", v))
	}
}

// ValueToObjectType ensures we have a types.Object literal
func ValueToObjectType(v attr.Value) types.Object {
	if vb, ok := v.(types.Object); ok {
		return vb
	} else if vb, ok := v.(*types.Object); ok {
		return *vb
	} else {
		panic(fmt.Sprintf("cannot pass type %T to conv.ValueToObjectType", v))
	}
}

// ValueToSetType ensures we have a types.Set literal
func ValueToSetType(v attr.Value) types.Set {
	if vb, ok := v.(types.Set); ok {
		return vb
	} else if vb, ok := v.(*types.Set); ok {
		return *vb
	} else {
		panic(fmt.Sprintf("cannot pass type %T to conv.ValueToSetType", v))
	}
}

// ValueToStringType ensures we have a types.String literal
func ValueToStringType(v attr.Value) types.String {
	if vb, ok := v.(types.String); ok {
		return vb
	} else if vb, ok := v.(*types.String); ok {
		return *vb
	} else {
		panic(fmt.Sprintf("cannot pass type %T to conv.StringValueToString", v))
	}
}

// TestAttributeValueState - Determine the state of the attribute value
//
// An Attribute Value can have one of 3 main states:
//  1. Unknown
//  2. Null
//  3. Valued
//
// This function tests every attr.Value implementation for each of these states, and additionally performs an
// "emptiness" check.
//
// There are bespoke errors for Unknown, Null, and Empty.  See errors.go for details on how to use these error types.
//
// This function will only ever return one error, going from greatest to least significance:
//
//	unknown > null > empty > non-empty
//
// A 'nil' response from this function means the attribute's value was defined to a non-"empty" value at runtime. See
// function body for a particular type if you're interested in what "empty" means.
func TestAttributeValueState(av attr.Value) error {
	var (
		empty bool

		undefined = av.IsUnknown()
		null      = av.IsNull()
	)

	switch av.(type) {
	case types.List, *types.List:
		tv := ValueToListType(av)
		undefined = tv.IsUnknown()
		null = tv.IsNull()
		empty = AttributeValueLength(av) == 0

	case types.Map, *types.Map:
		tv := ValueToMapType(av)
		undefined = tv.IsUnknown()
		null = tv.IsNull()
		empty = AttributeValueLength(av) == 0

	case types.Set, *types.Set:
		tv := ValueToSetType(av)
		undefined = tv.IsUnknown()
		null = tv.IsNull()
		empty = AttributeValueLength(av) > 0

	case types.String, *types.String:
		tv := ValueToStringType(av)
		undefined = tv.IsUnknown()
		null = tv.IsNull()
		empty = AttributeValueToString(av) == ""
	}

	if undefined {
		return ErrValueIsUnknown
	} else if null {
		return ErrValueIsNull
	} else if empty {
		return ErrValueIsEmpty
	}

	return nil
}

// BoolValueToString accepts an instance of either types.Bool or *types.Bool, attempting to convert the value to a string.
// DEPRECATED: use AttributeValueToString in all cases
func BoolValueToString(v attr.Value) string {
	return AttributeValueToString(v)
}

// Float64ValueToString accepts an instance of either types.Float64 or *types.Float64, attempting to convert the value to
// a string.
// DEPRECATED: use AttributeValueToString in all cases
func Float64ValueToString(v attr.Value) string {
	return AttributeValueToString(v)
}

// Int64ValueToString accepts an instance of either types.Int64 or *types.Int64, attempting to convert the value to a string.
// DEPRECATED: use AttributeValueToString in all cases
func Int64ValueToString(v attr.Value) string {
	return AttributeValueToString(v)
}

// NumberValueToString accepts an instance of either types.Number or *types.Number, attempting to convert the value to
// a string.
// DEPRECATED: use AttributeValueToString in all cases
func NumberValueToString(v attr.Value) string {
	return AttributeValueToString(v)
}

// StringValueToString accepts an instance of either types.String or *types.String, returning the raw string value
// DEPRECATED: use AttributeValueToString in all cases
func StringValueToString(v attr.Value) string {
	return AttributeValueToString(v)
}

// StringValueToBytes accepts an instance of either types.String or *types.String, returning the raw string value cast
// to a byte slice
func StringValueToBytes(v attr.Value) []byte {
	return []byte(AttributeValueToString(v))
}

// AttributeValueToString will attempt to execute the appropriate AttributeStringerFunc from the ones registered.
func AttributeValueToString(v attr.Value) string {
	if s, ok := v.(types.String); ok {
		return s.ValueString()
	}
	return v.String()
}

// AttributeValueToStrings attempts to convert the provided attr.Value into a slice of strings.
func AttributeValueToStrings(av attr.Value) []string {
	switch av.(type) {
	case types.List, *types.List:
		return StringListToStrings(av)

	case types.Set, *types.Set:
		return StringSetToStrings(av)
	default:
		out := make([]string, 0)
		out = append(out, AttributeValueToString(av))
		return out
	}
}

// LengthOfListValue returns the number of elements in the List attribute.  This will return 0 if the attribute was not set,
// set to null, or defined as an empty list.
func LengthOfListValue(v attr.Value) int {
	return len(ValueToListType(v).Elements())
}

// LengthOfMapValue returns the number of elements in the Map attribute.  This will return 0 if the attribute was not set,
// set to null, or defined as an empty map.
func LengthOfMapValue(v attr.Value) int {
	return len(ValueToMapType(v).Elements())
}

// LengthOfSetValue returns the number of elements in the Set attribute.  This will return 0 if the attribute was not set,
// set to null, or defined as an empty set.
func LengthOfSetValue(v attr.Value) int {
	return len(ValueToSetType(v).Elements())
}

// LengthOfStringValue returns the number of bytes in the String attribute.  This will return 0 if the attribute was not set,
// set to 0, or defined as an empty string.
func LengthOfStringValue(v attr.Value) int {
	return len(ValueToStringType(v).ValueString())
}

// AttributeValueLength attempts to determine the "length" of an attribute value, for types where that value has
// significance.
func AttributeValueLength(v attr.Value) int {
	switch v.(type) {
	case types.List, *types.List:
		return LengthOfListValue(v)

	case types.Map, *types.Map:
		return LengthOfMapValue(v)

	case types.Set, *types.Set:
		return LengthOfSetValue(v)

	case types.String, *types.String:
		return LengthOfStringValue(v)

	default:
		panic(fmt.Sprintf("unable to determine length of attribute value of type %T", v))
	}
}

// BoolValueToBool accepts either a types.Bool or *types.Bool and extracts the raw bool value within
func BoolValueToBool(v attr.Value) bool {
	return ValueToBoolType(v).ValueBool()
}

// BoolValueToBoolPtr accepts either a types.Bool or *types.Bool, extracting the raw bool value within and returning
// a pointer to a copy of that value
//
// If the Value is unknown or null, a nil is returned.
func BoolValueToBoolPtr(v attr.Value) *bool {
	vt := ValueToBoolType(v)
	if vt.IsUnknown() || vt.IsNull() {
		return nil
	}
	vPtr := new(bool)
	*vPtr = vt.ValueBool()
	return vPtr
}

// NumberValueToBigFloat accepts either a types.Number or *types.Number, returning the raw *big.Float value.  This may
// be nil if the value was not set.
func NumberValueToBigFloat(v attr.Value) *big.Float {
	return ValueToNumberType(v).ValueBigFloat()
}

// NumberValueToInt64 accepts either a types.Number or *types.Number, returning an int64 representation of the
// *big.Float value within.  It will return [0, big.Exact] of the value was not set.
func NumberValueToInt64(v attr.Value) (int64, big.Accuracy) {
	vt := ValueToNumberType(v)
	if vt.IsNull() || vt.IsUnknown() {
		return 0, big.Exact
	}
	return vt.ValueBigFloat().Int64()
}

// NumberValueToInt accepts either a types.Number or *types.Number, returning an int representation of the *big.Float
// value within.  It will return [0, big.Exact] if the value was not set
func NumberValueToInt(v attr.Value) (int, big.Accuracy) {
	iv, acc := NumberValueToInt64(v)
	return int(iv), acc
}

// NumberValueToFloat64 accepts either a types.Number or *types.Number, returning a float64 representation of the
// *big.Float value within.  It will return [0.0, big.Exact] of the value was not set.
func NumberValueToFloat64(v attr.Value) (float64, big.Accuracy) {
	vt := ValueToNumberType(v)
	if vt.IsUnknown() || vt.IsNull() {
		return 0.0, big.Exact
	}
	return vt.ValueBigFloat().Float64()
}

// Int64ValueToInt64 accepts either a types.Int64 or *types.Int64, returning the raw int64 value within
func Int64ValueToInt64(v attr.Value) int64 {
	return ValueToInt64Type(v).ValueInt64()
}

// Int64ValueToInt accepts either a types.Int64 or *types.Int64, returning an int representation of the value within
func Int64ValueToInt(v attr.Value) int {
	return int(Int64ValueToInt64(v))
}

// Int64ValueToIntPtr accepts either a types.Int64 or *types.Int64, returning a pointer to a copy of the int
// representation of the value within
//
// If the Value is unknown or null, a nil is returned.
func Int64ValueToIntPtr(v attr.Value) *int {
	vt := ValueToInt64Type(v)
	if vt.IsUnknown() || vt.IsNull() {
		return nil
	}
	vPtr := new(int)
	*vPtr = int(vt.ValueInt64())
	return vPtr
}

// Float64ValueToFloat64 accepts either a types.Float64 or *types.Float64, returning the raw float64 value within
func Float64ValueToFloat64(v attr.Value) float64 {
	return ValueToFloat64Type(v).ValueFloat64()
}

// Float64ValueToFloat32 accepts either a types.Float64 or *types.Float64, returning a float32 representation of the
// raw float64 value
func Float64ValueToFloat32(v attr.Value) float32 {
	return float32(Float64ValueToFloat64(v))
}

// StringValueToFloat64 accepts either a types.String or *types.string, attempting to parse the value as a float64
func StringValueToFloat64(v attr.Value) (float64, error) {
	return strconv.ParseFloat(ValueToStringType(v).ValueString(), 64)
}

// StringValueToInt64 accepts either a types.String or *types.String, attempting to parse the value as an int64.
func StringValueToInt64(v attr.Value) (int, error) {
	return strconv.Atoi(ValueToStringType(v).ValueString())
}

// StringValueToStringPtr accepts an instance of either types.String or *types.String, returning a pointer to a copy
// of the raw string value
//
// If the Value is unknown or null, a nil is returned.
func StringValueToStringPtr(v attr.Value) *string {
	vt := ValueToStringType(v)
	if vt.IsUnknown() || vt.IsNull() {
		return nil
	}
	vPtr := new(string)
	*vPtr = vt.ValueString()
	return vPtr
}

// StringListToStrings accepts an instance of either types.List or *types.List where ElementType MUST be types.StringType,
// returning a slice of strings of the value of each element
func StringListToStrings(v attr.Value) []string {
	vt := ValueToListType(v)
	out := make([]string, len(vt.Elements()))
	for i, ve := range vt.Elements() {
		out[i] = AttributeValueToString(ve)
	}
	return out
}

// StringSetToStrings accepts an instance of either types.Set or *types.Set where ElementType MUST be types.StringType,
// returning a slice of strings of the value of each element
func StringSetToStrings(v attr.Value) []string {
	vt := ValueToSetType(v)
	out := make([]string, len(vt.Elements()))
	for i, ve := range vt.Elements() {
		out[i] = AttributeValueToString(ve)
	}
	return out
}

// Int64ListToInts accepts an instance of either types.List or *types.List where ElementType MUST be types.Int64Type,
// returning a slice of ints of the value of each element.
func Int64ListToInts(v attr.Value) []int {
	vt := ValueToListType(v)
	out := make([]int, len(vt.Elements()))
	for i, ve := range vt.Elements() {
		out[i] = Int64ValueToInt(ve)
	}
	return out
}

// Int64SetToInts accepts an instance of either types.Set or *types.set where ElementType MUST be types.Int64Type
// returning a slice of ints of the value of each element
func Int64SetToInts(v attr.Value) []int {
	vt := ValueToSetType(v)
	out := make([]int, len(vt.Elements()))
	for i, ve := range vt.Elements() {
		out[i] = Int64ValueToInt(ve)
	}
	return out
}

// NumberListToInts accepts either an instance of types.List or *types.List where ElementType MUST be types.NumberType
// returning a slice of ints of the value of each element
func NumberListToInts(v attr.Value) []int {
	vt := ValueToListType(v)
	out := make([]int, len(vt.Elements()))
	for i, ve := range vt.Elements() {
		iv, _ := NumberValueToInt(ve)
		out[i] = iv
	}
	return out
}

// NumberSetToInts accepts either an instance of types.Set or *types.Set where ElementType MUST be types.NumberType
// returning a slice of ints of the value of each element
func NumberSetToInts(v attr.Value) []int {
	vt := ValueToSetType(v)
	out := make([]int, len(vt.Elements()))
	for i, ve := range vt.Elements() {
		iv, _ := NumberValueToInt(ve)
		out[i] = iv
	}
	return out
}

// AttributeValueToFloat64 accepts either a literal or pointer to a concrete attr.Value implementation, attempting to
// to return a float64 representation of its value.
func AttributeValueToFloat64(v attr.Value) (float64, big.Accuracy, error) {
	switch v.(type) {
	case types.Float64, *types.Float64:
		return Float64ValueToFloat64(v), big.Exact, nil

	case types.Int64, *types.Int64:
		return float64(Int64ValueToInt64(v)), big.Exact, nil

	case types.Number, *types.Number:
		f, a := NumberValueToFloat64(v)
		return f, a, nil

	case types.String, *types.String:
		f, err := StringValueToFloat64(v)
		return f, big.Exact, err

	default:
		return 0, 0, ValueTypeUnhandledError("attr_to_float64", v)
	}
}

// AttributeValueToInt64 accepts either a literal or pointer to a concrete attr.Value implementation, attempting to
// return an int64 representation of its value.
func AttributeValueToInt64(v attr.Value) (int64, big.Accuracy, error) {
	switch v.(type) {
	case types.Float64, *types.Float64:
		f := Float64ValueToFloat64(v)
		i := int64(f)
		if f > float64(i) {
			return i, big.Below, nil
		} else {
			return i, big.Exact, nil
		}

	case types.Int64, *types.Int64:
		return Int64ValueToInt64(v), big.Exact, nil

	case types.Number, *types.Number:
		i, a := NumberValueToInt64(v)
		return i, a, nil

	case types.String, *types.String:
		i, err := StringValueToInt64(v)
		return int64(i), big.Exact, err

	default:
		return 0, 0, ValueTypeUnhandledError("attr_to_int64", v)
	}
}

// AttributeValueToBigFloat accepts either a literal or pointer to a concrete attr.Value implementation, attempting to
// return a *big.Float instance of its value.
func AttributeValueToBigFloat(v attr.Value) (*big.Float, error) {
	switch v.(type) {
	case types.Float64, *types.Float64:
		return big.NewFloat(Float64ValueToFloat64(v)), nil

	case types.Int64, *types.Int64:
		return big.NewFloat(0).SetInt64(Int64ValueToInt64(v)), nil

	case types.Number, *types.Number:
		return NumberValueToBigFloat(v), nil

	case types.String, *types.String:
		bf, _, err := big.ParseFloat(AttributeValueToString(v), 10, FloatPrecision, big.ToZero)
		return bf, err

	default:
		return nil, ValueTypeUnhandledError("attr_to_bigfloat", v)
	}
}

// BoolToBoolValue takes a bool and wraps it up as a types.Bool
// DEPRECATED: use types.BoolValue() directly
func BoolToBoolValue(b bool) types.Bool {
	return types.BoolValue(b)
}

// BoolPtrToBoolValue accepts a bool pointer, returning a types.Bool with the dereferenced value.
//
// If the provided pointer is nil, the returned Bool type will be set as Null.
func BoolPtrToBoolValue(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

// Int64ToInt64Value takes an int64 and wraps it up as a types.Int64
// DEPRECATED: use types.Int64Value() directly
func Int64ToInt64Value(i int64) types.Int64 {
	return types.Int64Value(i)
}

// Int64ToNumberValue takes an int64 and wraps it up as a types.Number
func Int64ToNumberValue(i int64) types.Number {
	return types.NumberValue(new(big.Float).SetInt64(i))
}

// IntToInt64Value takes an int and wraps it up as a types.Int64
func IntToInt64Value(i int) types.Int64 {
	return types.Int64Value(int64(i))
}

// IntPtrToInt64Value takes an *int and wraps it up as a types.Int64
//
// If the go value is nil, Null will be true on the outgoing attr.Value type
func IntPtrToInt64Value(i *int) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*i))
}

// IntToNumberValue takes an int and wraps it up as a types.Number
func IntToNumberValue(i int) types.Number {
	return Int64ToNumberValue(int64(i))
}

// Float64ToFloat64Value takes a float64 and wraps it up as a types.Float64
// DEPRECATED: use types.Float64Value() directly
func Float64ToFloat64Value(f float64) types.Float64 {
	return types.Float64Value(f)
}

// Float64ToNumberValue takes a float64 and wraps it up as a types.Number
func Float64ToNumberValue(f float64) types.Number {
	return types.NumberValue(big.NewFloat(f))
}

// Float32ToFloat64Value takes a float32 and wraps it up as a types.Float64
func Float32ToFloat64Value(f float32) types.Float64 {
	return types.Float64Value(float64(f))
}

// Float32ToNumberValue takes a float32 and wraps it up as a types.Number
func Float32ToNumberValue(f float32) types.Number {
	return Float64ToNumberValue(float64(f))
}

// StringToStringValue takes a string and wraps it up as a types.String
// DEPRECATED: use types.StringValue() directly
func StringToStringValue(s string) types.String {
	return types.StringValue(s)
}

// BytesToStringValue takes a byte slice and wraps it as a types.String.  If the provided slice is `nil`, then the
// resulting String type will be marked as "null".
func BytesToStringValue(b []byte) types.String {
	if b == nil {
		return types.StringNull()
	}
	return types.StringValue(string(b))
}

// StringPtrToStringValue takes a *string and wraps it up as a types.String
//
// If the go value is nil, Null will be true on the outgoing attr.Value type
func StringPtrToStringValue(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}

// StringsToStringList takes a slice of strings and creates a typed types.List with an ElementType of types.String
// and each value of Elements being an instance of types.String
//
// If nullOnEmpty parameter is `true`, the returned types.List will be set to Null.  This can be used to
// avoid Terraform state inconsistencies under certain circumstances.
func StringsToStringList(in []string, nullOnEmpty bool) types.List {
	inLen := len(in)

	if nullOnEmpty && inLen == 0 {
		return types.ListNull(types.StringType)
	}

	elems := make([]attr.Value, inLen)
	for i, n := range in {
		elems[i] = types.StringValue(n)
	}

	return types.ListValueMust(types.StringType, elems)
}

// StringsToStringSet takes a slice of strings and creates a typed types.Set with an ElementType of types.String
// and each value of Elements being an instance of types.String
//
// If nullOnEmpty parameter is `true`, the returned types.Set will be set to Null.  This can be used to
// avoid Terraform state inconsistencies under certain circumstances.
func StringsToStringSet(in []string, nullOnEmpty bool) types.Set {
	inLen := len(in)

	if nullOnEmpty && inLen == 0 {
		return types.SetNull(types.StringType)
	}

	elems := make([]attr.Value, inLen)
	for i, n := range in {
		elems[i] = types.StringValue(n)
	}

	return types.SetValueMust(types.StringType, elems)
}

// IntsToInt64List takes a slice of ints and creates a typed types.List with a ElementType of types.Int64Type and each
// value of Elements being an instance of types.Int64
//
// If nullOnEmpty parameter is `true`, the returned types.List will be set to Null.  This can be used to
// avoid Terraform state inconsistencies under certain circumstances.
func IntsToInt64List(in []int, nullOnEmpty bool) types.List {
	inLen := len(in)

	if nullOnEmpty && inLen == 0 {
		return types.ListNull(types.Int64Type)
	}

	elems := make([]attr.Value, inLen)
	for i, n := range in {
		elems[i] = IntToInt64Value(n)
	}

	return types.ListValueMust(types.Int64Type, elems)
}

// IntsToInt64Set takes a slice of ints and creates a typed types.Set with an ElementType of types.Int64Type and each
// value of Elements being an instance of types.Int64
//
// If nullOnEmpty parameter is `true`, the returned types.Set will be set to Null.  This can be used to
// avoid Terraform state inconsistencies under certain circumstances.
func IntsToInt64Set(in []int, nullOnEmpty bool) types.Set {
	inLen := len(in)

	if nullOnEmpty && inLen == 0 {
		return types.SetNull(types.Int64Type)
	}

	elems := make([]attr.Value, inLen)
	for i, n := range in {
		elems[i] = IntToInt64Value(n)
	}

	return types.SetValueMust(types.Int64Type, elems)
}
