package conv

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TFTypeValueType uint8

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
// 		1. Unknown
//  	2. Null
//  	3. Valued
//
// This function tests every attr.Value implementation for each of these states, and additionally performs an
// "emptiness" check.
//
// There are bespoke errors for Unknown, Null, and Empty.  See errors.go for details on how to use these error types.
//
// This function will only ever return one error, going from greatest to least significance:
//
// 		unknown > null > empty > non-empty
//
// A 'nil' response from this function means the attribute's value was defined to a non-"empty" value at runtime. See
// function body for a particular type if you're interested in what "empty" means.
func TestAttributeValueState(v attr.Value) error {
	var (
		undefined bool
		null      bool
		empty     bool
	)

	// first, check for unknown and null
	switch v.(type) {
	case types.Bool, *types.Bool:
		tv := ValueToBoolType(v)
		undefined = tv.Unknown
		null = tv.Null
	// bool values cannot be "empty"

	case types.Float64, *types.Float64:
		tv := ValueToFloat64Type(v)
		undefined = tv.Unknown
		null = tv.Null
	// float values cannot be "empty"

	case types.Int64, *types.Int64:
		tv := ValueToInt64Type(v)
		undefined = tv.Unknown
		null = tv.Null
	// int values cannot be "empty"

	case types.List, *types.List:
		tv := ValueToListType(v)
		undefined = tv.Unknown
		null = tv.Null
		empty = AttributeValueLength(v) == 0

	case types.Map, *types.Map:
		tv := ValueToMapType(v)
		undefined = tv.Unknown
		null = tv.Null
		empty = AttributeValueLength(v) == 0

	case types.Number, *types.Number:
		tv := ValueToNumberType(v)
		undefined = tv.Unknown
		null = tv.Null

	case types.Object, *types.Object:
		tv := ValueToObjectType(v)
		undefined = tv.Unknown
		null = tv.Null
	// todo: implement object "emptiness" check

	case types.Set, *types.Set:
		tv := ValueToSetType(v)
		undefined = tv.Unknown
		null = tv.Null
		empty = AttributeValueLength(v) > 0

	case types.String, *types.String:
		tv := ValueToStringType(v)
		undefined = tv.Unknown
		null = tv.Null
		empty = StringValueToString(v) == ""

	default:
		panic(fmt.Sprintf("no way to test for valued state for types %T", v))
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
func BoolValueToString(v attr.Value) string {
	return strconv.FormatBool(ValueToBoolType(v).Value)
}

// Float64ValueToString accepts an instance of either types.Float64 or *types.Float64, attempting to convert the value to
// a string.
func Float64ValueToString(v attr.Value) string {
	return strconv.FormatFloat(ValueToFloat64Type(v).Value, 'g', int(FloatPrecision), 64)
}

// Int64ValueToString accepts an instance of either types.Int64 or *types.Int64, attempting to convert the value to a string.
func Int64ValueToString(v attr.Value) string {
	return strconv.FormatInt(ValueToInt64Type(v).Value, 10)
}

// NumberValueToString accepts an instance of either types.Number or *types.Number, attempting to convert the value to
// a string.
func NumberValueToString(v attr.Value) string {
	return ValueToNumberType(v).Value.String()
}

// StringValueToString accepts an instance of either types.String or *types.String, returning the raw string value
func StringValueToString(v attr.Value) string {
	return ValueToStringType(v).Value
}

// StringValueToBytes accepts an instance of either types.String or *types.String, returning the raw string value cast
// to a byte slice
func StringValueToBytes(v attr.Value) []byte {
	return []byte(StringValueToString(v))
}

// AttributeValueToString will attempt to execute the appropriate AttributeStringerFunc from the ones registered.
func AttributeValueToString(v attr.Value) string {
	switch v.(type) {
	case types.Bool, *types.Bool:
		return BoolValueToString(v)

	case types.Float64, *types.Float64:
		return Float64ValueToString(v)

	case types.Int64, *types.Int64:
		return Int64ValueToString(v)

	case types.Number, *types.Number:
		return NumberValueToString(v)

	case types.String, *types.String:
		return StringValueToString(v)

	default:
		panic(fmt.Sprintf("no stringer func registered for type %T", v))
	}
}

// LengthOfListValue returns the number of elements in the List attribute.  This will return 0 if the attribute was not set,
// set to null, or defined as an empty list.
func LengthOfListValue(v attr.Value) int {
	return len(ValueToListType(v).Elems)
}

// LengthOfMapValue returns the number of elements in the Map attribute.  This will return 0 if the attribute was not set,
// set to null, or defined as an empty map.
func LengthOfMapValue(v attr.Value) int {
	return len(ValueToMapType(v).Elems)
}

// LengthOfSetValue returns the number of elements in the Set attribute.  This will return 0 if the attribute was not set,
// set to null, or defined as an empty set.
func LengthOfSetValue(v attr.Value) int {
	return len(ValueToSetType(v).Elems)
}

// LengthOfStringValue returns the number of bytes in the String attribute.  This will return 0 if the attribute was not set,
// set to 0, or defined as an empty string.
func LengthOfStringValue(v attr.Value) int {
	return len(ValueToStringType(v).Value)
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
	return ValueToBoolType(v).Value
}

// BoolValueToBoolPtr accepts either a types.Bool or *types.Bool, extracting the raw bool value within and returning
// a pointer to a copy of that value
//
// If the Value is unknown or null, a nil is returned.
func BoolValueToBoolPtr(v attr.Value) *bool {
	vt := ValueToBoolType(v)
	if vt.Unknown || vt.Null {
		return nil
	}
	vPtr := new(bool)
	*vPtr = vt.Value
	return vPtr
}

// NumberValueToBigFloat accepts either a types.Number or *types.Number, returning the raw *big.Float value.  This may
// be nil if the value was not set.
func NumberValueToBigFloat(v attr.Value) *big.Float {
	return ValueToNumberType(v).Value
}

// NumberValueToInt64 accepts either a types.Number or *types.Number, returning an int64 representation of the
// *big.Float value within.  It will return [0, big.Exact] of the value was not set.
func NumberValueToInt64(v attr.Value) (int64, big.Accuracy) {
	vt := ValueToNumberType(v)
	if vt.Value == nil {
		return 0, big.Exact
	}
	return vt.Value.Int64()
}

// NumberValueToFloat64 accepts either a types.Number or *types.Number, returning a float64 representation of the
// *big.Float value within.  It will return [0.0, big.Exact] of the value was not set.
func NumberValueToFloat64(v attr.Value) (float64, big.Accuracy) {
	vt := ValueToNumberType(v)
	if vt.Value == nil {
		return 0.0, big.Exact
	}
	return vt.Value.Float64()
}

// Int64ValueToInt64 accepts either a types.Int64 or *types.Int64, returning the raw int64 value within
func Int64ValueToInt64(v attr.Value) int64 {
	return ValueToInt64Type(v).Value
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
	if vt.Unknown || vt.Null {
		return nil
	}
	vPtr := new(int)
	*vPtr = int(vt.Value)
	return vPtr
}

// Float64ValueToFloat64 accepts either a types.Float64 or *types.Float64, returning the raw float64 value within
func Float64ValueToFloat64(v attr.Value) float64 {
	return ValueToFloat64Type(v).Value
}

// Float64ValueToFloat32 accepts either a types.Float64 or *types.Float64, returning a float32 representation of the
// raw float64 value
func Float64ValueToFloat32(v attr.Value) float32 {
	return float32(Float64ValueToFloat64(v))
}

// StringValueToFloat64 accepts either a types.String or *types.string, attempting to parse the value as a float64
func StringValueToFloat64(v attr.Value) (float64, error) {
	return strconv.ParseFloat(ValueToStringType(v).Value, 64)
}

// StringValueToInt64 accepts either a types.String or *types.String, attempting to parse the value as an int64.
func StringValueToInt64(v attr.Value) (int, error) {
	return strconv.Atoi(ValueToStringType(v).Value)
}

// StringValueToStringPtr accepts an instance of either types.String or *types.String, returning a pointer to a copy
// of the raw string value
//
// If the Value is unknown or null, a nil is returned.
func StringValueToStringPtr(v attr.Value) *string {
	vt := ValueToStringType(v)
	if vt.Unknown || vt.Null {
		return nil
	}
	vPtr := new(string)
	*vPtr = vt.Value
	return vPtr
}

// Int64ListToInts accepts an instance of either types.List or *types.List where ElemType MUST be types.Int64Type,
// returning a slice of ints of the value of each element.
func Int64ListToInts(v attr.Value) []int {
	vt := ValueToListType(v)
	out := make([]int, len(vt.Elems))
	for i, ve := range vt.Elems {
		out[i] = Int64ValueToInt(ve.(types.Int64))
	}
	return out
}

// Int64SetToInts accepts an instance of either types.Set or *types.set where ElemType MUST be types.Int64Type
// returning a slice of ints of the value of each element
func Int64SetToInts(v attr.Value) []int {
	vt := ValueToSetType(v)
	out := make([]int, len(vt.Elems))
	for i, ve := range vt.Elems {
		out[i] = Int64ValueToInt(ve.(types.Int64))
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
		panic(fmt.Sprintf("unable to determine float64 value of type %T", v))
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
		panic(fmt.Sprintf("unable to determine int64 value of type %T", v))
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
		bf, _, err := big.ParseFloat(StringValueToString(v), 10, FloatPrecision, big.ToZero)
		return bf, err

	default:
		panic(fmt.Sprintf("unable to parse type %T into *big.Float", v))
	}
}

// BoolToBoolValue takes a bool and wraps it up as a types.Bool
func BoolToBoolValue(b bool) types.Bool {
	return types.Bool{Value: b}
}

// Int64ToInt64Value takes an int64 and wraps it up as a types.Int64
func Int64ToInt64Value(i int64) types.Int64 {
	return types.Int64{Value: i}
}

// Int64ToNumberValue takes an int64 and wraps it up as a types.Number
func Int64ToNumberValue(i int64) types.Number {
	return types.Number{Value: new(big.Float).SetInt64(i)}
}

// IntToInt64Value takes an int and wraps it up as a types.Int64
func IntToInt64Value(i int) types.Int64 {
	return Int64ToInt64Value(int64(i))
}

// IntPtrToInt64Value takes an *int and wraps it up as a types.Int64
//
// If the go value is nil, Null will be true on the outgoing attr.Value type
func IntPtrToInt64Value(i *int) types.Int64 {
	if i == nil {
		return types.Int64{Null: true}
	}
	return types.Int64{Value: int64(*i)}
}

// IntToNumberValue takes an int and wraps it up as a types.Number
func IntToNumberValue(i int) types.Number {
	return Int64ToNumberValue(int64(i))
}

// Float64ToFloat64Value takes a float64 and wraps it up as a types.Float64
func Float64ToFloat64Value(f float64) types.Float64 {
	return types.Float64{Value: f}
}

// Float64ToNumberValue takes a float64 and wraps it up as a types.Number
func Float64ToNumberValue(f float64) types.Number {
	return types.Number{Value: big.NewFloat(f)}
}

// Float32ToFloat64Value takes a float32 and wraps it up as a types.Float64
func Float32ToFloat64Value(f float32) types.Float64 {
	return Float64ToFloat64Value(float64(f))
}

// Float32ToNumberValue takes a float32 and wraps it up as a types.Number
func Float32ToNumberValue(f float32) types.Number {
	return Float64ToNumberValue(float64(f))
}

// StringToStringValue takes a string and wraps it up as a types.String
func StringToStringValue(s string) types.String {
	return types.String{Value: s}
}

// BytesToStringValue takes a byte slice and wraps it as a types.String.  If the provided slice is `nil`, then the
// resulting String type will be marked as "null".
func BytesToStringValue(b []byte) types.String {
	if b == nil {
		return types.String{Null: true}
	}
	return StringToStringValue(string(b))
}

// StringPtrToStringValue takes a *string and wraps it up as a types.String
// If the go value is nil, Null will be true on the outgoing attr.Value type
func StringPtrToStringValue(s *string) types.String {
	if s == nil {
		return types.String{Null: true}
	}
	return types.String{Value: *s}
}

// IntsToInt64List takes a slice of ints and creates a typed types.List with a ElemType of types.Int64Type and each
// value of Elems being an instance of types.Int64
//
// If nullOnEmpty parameter is `true`, sets the Null field on the returned types.List to `true`.  This can be used to
// avoid Terraform state inconsistencies under certain circumstances.
func IntsToInt64List(in []int, nullOnEmpty bool) types.List {
	inLen := len(in)
	lt := types.List{
		Elems:    make([]attr.Value, len(in)),
		ElemType: types.Int64Type,
	}
	if nullOnEmpty && inLen == 0 {
		lt.Null = true
	} else {
		for i, n := range in {
			lt.Elems[i] = IntToInt64Value(n)
		}
	}
	return lt
}

// IntsToInt64Set takes a slice of ints and creates a typed types.Set with an ElemType of types.Int64Type and each
// value of Elems being an instance of types.Int64
//
// If nullOnEmpty parameter is `true`, sets the Null field on the returned types.Set to `true`.  This can be used to
// avoid Terraform state inconsistencies under certain circumstances.
func IntsToInt64Set(in []int, nullOnEmpty bool) types.Set {
	inLen := len(in)
	st := types.Set{
		Elems:    make([]attr.Value, inLen),
		ElemType: types.Int64Type,
	}
	if nullOnEmpty && inLen == 0 {
		st.Null = true
	} else {
		for i, n := range in {
			st.Elems[i] = IntToInt64Value(n)
		}
	}
	return st
}
