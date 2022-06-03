package validation_test

import (
	"math/big"
	"testing"

	"github.com/dcarbone/terraform-plugin-framework-utils/validation"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type comparisonTest struct {
	name        string
	op          validation.CompareOp
	tgt         interface{}
	act         attr.Value
	meta        []interface{}
	expectPanic bool
	expectError bool
}

func (ct comparisonTest) do(t *testing.T) {
	if ct.expectPanic {
		defer func() {
			if v := recover(); v != nil {
				t.Logf("Panic seen: %v", v)
			} else {
				t.Log("panic expected")
				t.Fail()
			}
		}()
	}
	err := validation.CompareAttrValues(ct.act, ct.op, ct.tgt, ct.meta...)
	if err != nil {
		t.Logf("comparison error seen: %v", err)
		if !ct.expectError {
			t.Fail()
		}
	} else if ct.expectError {
		t.Log("expected comparison error")
		t.Fail()
	}
}

func TestComparison_Bool(t *testing.T) {
	theTests := []comparisonTest{
		{
			name: "bool_eq_ok",
			op:   validation.Equal,
			tgt:  true,
			act:  types.Bool{Value: true},
		},
		{
			name:        "bool_eq_nok",
			op:          validation.Equal,
			tgt:         true,
			act:         types.Bool{Value: false},
			expectError: true,
		},
		{
			name: "bool_neq_ok",
			op:   validation.NotEqual,
			tgt:  true,
			act:  types.Bool{Value: false},
		},
		{
			name:        "bool_neq_nok",
			op:          validation.NotEqual,
			tgt:         true,
			act:         types.Bool{Value: true},
			expectError: true,
		},
	}

	for _, ct := range theTests {
		t.Run(ct.name, func(t *testing.T) {
			ct.do(t)
		})
	}
}

func TestComparison_Float(t *testing.T) {
	theTests := []comparisonTest{
		{
			name: "float_eq_ok",
			op:   validation.Equal,
			tgt:  float64(0.1),
			act:  types.Float64{Value: 0.1},
		},
		{
			name:        "float_eq_nok",
			op:          validation.Equal,
			tgt:         float64(1.0),
			act:         types.Float64{Value: 0.1},
			expectError: true,
		},
		{
			name: "float_neq_ok",
			op:   validation.NotEqual,
			tgt:  1.0,
			act:  types.Float64{Value: 0.1},
		},
		{
			name:        "float_neq_nok",
			op:          validation.NotEqual,
			tgt:         float64(1.0),
			act:         types.Float64{Value: 1.0},
			expectError: true,
		},
		{
			name: "float_lt_ok",
			op:   validation.LessThan,
			tgt:  float64(1.0),
			act:  types.Float64{Value: 0.9},
		},
		{
			name:        "float_lt_nok",
			op:          validation.LessThan,
			tgt:         float64(1.0),
			act:         types.Float64{Value: 1.0},
			expectError: true,
		},
		{
			name: "float_lte_lt_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  float64(1.0),
			act:  types.Float64{Value: 0.9},
		},
		{
			name:        "float_lte_lt_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         float64(1.0),
			act:         types.Float64{Value: 1.1},
			expectError: true,
		},
		{
			name: "float_lte_eq_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  float64(1.0),
			act:  types.Float64{Value: 1.0},
		},
		{
			name:        "float_lte_eq_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         float64(1.0),
			act:         types.Float64{Value: 1.1},
			expectError: true,
		},
		{
			name: "float_gt_ok",
			op:   validation.GreaterThan,
			tgt:  float64(1.0),
			act:  types.Float64{Value: 1.1},
		},
		{
			name:        "float_gt_nok",
			op:          validation.GreaterThan,
			tgt:         float64(1.0),
			act:         types.Float64{Value: 1.0},
			expectError: true,
		},
		{
			name: "float_gte_gt_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  float64(1.0),
			act:  types.Float64{Value: 1.1},
		},
		{
			name:        "float_gte_gt_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         float64(1.0),
			act:         types.Float64{Value: 0.9},
			expectError: true,
		},
		{
			name: "float_gte_eq_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  float64(1.0),
			act:  types.Float64{Value: 1.0},
		},
		{
			name:        "float_gte_eq_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         float64(1.0),
			act:         types.Float64{Value: 0.9},
			expectError: true,
		},
	}

	for _, ct := range theTests {
		t.Run(ct.name, func(t *testing.T) {
			ct.do(t)
		})
	}
}

func TestComparison_Int(t *testing.T) {
	theTests := []comparisonTest{
		{
			name: "int_eq_ok",
			op:   validation.Equal,
			tgt:  1,
			act:  types.Int64{Value: 1},
		},
		{
			name:        "int_eq_nok",
			op:          validation.Equal,
			tgt:         1,
			act:         types.Int64{Value: 0},
			expectError: true,
		},
		{
			name: "int_neq_ok",
			op:   validation.NotEqual,
			tgt:  1,
			act:  types.Int64{Value: 0},
		},
		{
			name:        "int_neq_nok",
			op:          validation.NotEqual,
			tgt:         1,
			act:         types.Int64{Value: 1},
			expectError: true,
		},
		{
			name: "int_lt_ok",
			op:   validation.LessThan,
			tgt:  1,
			act:  types.Int64{Value: 0},
		},
		{
			name:        "int_lt_nok",
			op:          validation.LessThan,
			tgt:         1,
			act:         types.Int64{Value: 1},
			expectError: true,
		},
		{
			name: "int_lte_lt_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  1,
			act:  types.Int64{Value: 0},
		},
		{
			name:        "int_lte_lt_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         1,
			act:         types.Int64{Value: 2},
			expectError: true,
		},
		{
			name: "int_lte_eq_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  1,
			act:  types.Int64{Value: 1},
		},
		{
			name:        "int_lte_eq_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         1,
			act:         types.Int64{Value: 2},
			expectError: true,
		},
		{
			name: "int_gt_ok",
			op:   validation.GreaterThan,
			tgt:  1,
			act:  types.Int64{Value: 2},
		},
		{
			name:        "int_gt_nok",
			op:          validation.GreaterThan,
			tgt:         1,
			act:         types.Int64{Value: 1},
			expectError: true,
		},
		{
			name: "int_gte_gt_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  1,
			act:  types.Int64{Value: 2},
		},
		{
			name:        "int_gte_gt_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         1,
			act:         types.Int64{Value: 0},
			expectError: true,
		},
		{
			name: "int_gte_eq_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  1,
			act:  types.Int64{Value: 1},
		},
		{
			name:        "int_gte_eq_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         1,
			act:         types.Int64{Value: 0},
			expectError: true,
		},
	}

	for _, ct := range theTests {
		t.Run(ct.name, func(t *testing.T) {
			ct.do(t)
		})
	}
}

func TestComparison_BigFloat(t *testing.T) {
	theTests := []comparisonTest{
		{
			name: "bigfloat_eq_ok",
			op:   validation.Equal,
			tgt:  big.NewFloat(0.1),
			act:  types.Number{Value: big.NewFloat(0.1)},
		},
		{
			name:        "bigfloat_eq_nok",
			op:          validation.Equal,
			tgt:         big.NewFloat(1.0),
			act:         types.Number{Value: big.NewFloat(0.1)},
			expectError: true,
		},
		{
			name: "bigfloat_neq_ok",
			op:   validation.NotEqual,
			tgt:  1.0,
			act:  types.Number{Value: big.NewFloat(0.1)},
		},
		{
			name:        "bigfloat_neq_nok",
			op:          validation.NotEqual,
			tgt:         big.NewFloat(1.0),
			act:         types.Number{Value: big.NewFloat(1.0)},
			expectError: true,
		},
		{
			name: "bigfloat_lt_ok",
			op:   validation.LessThan,
			tgt:  big.NewFloat(1.0),
			act:  types.Number{Value: big.NewFloat(0.9)},
		},
		{
			name:        "bigfloat_lt_nok",
			op:          validation.LessThan,
			tgt:         big.NewFloat(1.0),
			act:         types.Number{Value: big.NewFloat(1.0)},
			expectError: true,
		},
		{
			name: "bigfloat_lte_lt_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  big.NewFloat(1.0),
			act:  types.Number{Value: big.NewFloat(0.9)},
		},
		{
			name:        "bigfloat_lte_lt_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         big.NewFloat(1.0),
			act:         types.Number{Value: big.NewFloat(1.1)},
			expectError: true,
		},
		{
			name: "bigfloat_lte_eq_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  big.NewFloat(1.0),
			act:  types.Number{Value: big.NewFloat(1.0)},
		},
		{
			name:        "bigfloat_lte_eq_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         big.NewFloat(1.0),
			act:         types.Number{Value: big.NewFloat(1.1)},
			expectError: true,
		},
		{
			name: "bigfloat_gt_ok",
			op:   validation.GreaterThan,
			tgt:  big.NewFloat(1.0),
			act:  types.Number{Value: big.NewFloat(1.1)},
		},
		{
			name:        "bigfloat_gt_nok",
			op:          validation.GreaterThan,
			tgt:         big.NewFloat(1.0),
			act:         types.Number{Value: big.NewFloat(1.0)},
			expectError: true,
		},
		{
			name: "bigfloat_gte_gt_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  big.NewFloat(1.0),
			act:  types.Number{Value: big.NewFloat(1.1)},
		},
		{
			name:        "bigfloat_gte_gt_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         big.NewFloat(1.0),
			act:         types.Number{Value: big.NewFloat(0.9)},
			expectError: true,
		},
		{
			name: "bigfloat_gte_eq_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  big.NewFloat(1.0),
			act:  types.Number{Value: big.NewFloat(1.0)},
		},
		{
			name:        "bigfloat_gte_eq_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         big.NewFloat(1.0),
			act:         types.Number{Value: big.NewFloat(0.9)},
			expectError: true,
		},
	}

	for _, ct := range theTests {
		t.Run(ct.name, func(t *testing.T) {
			ct.do(t)
		})
	}
}

func TestComparison_String(t *testing.T) {
	theTests := []comparisonTest{
		{
			name: "string_eq_ok",
			op:   validation.Equal,
			act:  types.String{Value: "hi"},
			tgt:  "hi",
		},
		{
			name:        "string_eq_nok",
			op:          validation.Equal,
			act:         types.String{Value: "hi"},
			tgt:         "ih",
			expectError: true,
		},
		{
			name: "string_neq_ok",
			op:   validation.NotEqual,
			act:  types.String{Value: "hi"},
			tgt:  "ih",
		},
		{
			name:        "string_neq_nok",
			op:          validation.NotEqual,
			act:         types.String{Value: "hi"},
			tgt:         "hi",
			expectError: true,
		},
	}

	for _, ct := range theTests {
		t.Run(ct.name, func(t *testing.T) {
			ct.do(t)
		})
	}
}

func TestComparison_Strings(t *testing.T) {
	const (
		one   = "one"
		two   = "two"
		oNe   = "oNe"
		twO   = "twO"
		three = "three"
	)

	var (
		targetOneTwo = []string{one, two}
	)

	theTests := []comparisonTest{
		// list []string sensitive eq
		{
			name: "strings_list_eq_sensitive_ok",
			op:   validation.Equal,
			act:  types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: one}, types.String{Value: two}}},
			tgt:  targetOneTwo,
		},
		{
			name:        "strings_list_eq_sensitive_nok_casing",
			op:          validation.Equal,
			act:         types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: oNe}, types.String{Value: twO}}},
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "strings_list_eq_sensitive_nok_order",
			op:          validation.Equal,
			act:         types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: two}, types.String{Value: one}}},
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "strings_list_eq_sensitive_nok_len",
			op:          validation.Equal,
			act:         types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: one}, types.String{Value: two}, types.String{Value: three}}},
			tgt:         targetOneTwo,
			expectError: true,
		},

		// list []string insensitive eq
		{
			name: "strings_list_eq_insensitive_ok",
			op:   validation.Equal,
			act:  types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: one}, types.String{Value: two}}},
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name: "strings_list_eq_insensitive_ok_casing",
			op:   validation.Equal,
			act:  types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: oNe}, types.String{Value: twO}}},
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name:        "strings_list_eq_insensitive_nok_order",
			op:          validation.Equal,
			act:         types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: two}, types.String{Value: one}}},
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},
		{
			name:        "strings_list_eq_insensitive_nok_len",
			op:          validation.Equal,
			act:         types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: one}, types.String{Value: two}, types.String{Value: three}}},
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},

		// list []string sensitive neq
		{
			name: "strings_list_neq_sensitive_ok_order",
			op:   validation.NotEqual,
			act:  types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: two}, types.String{Value: one}}},
			tgt:  targetOneTwo,
		},
		{
			name: "strings_list_neq_sensitive_ok_casing",
			op:   validation.NotEqual,
			act:  types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: oNe}, types.String{Value: twO}}},
			tgt:  targetOneTwo,
		},
		{
			name: "strings_list_neq_sensitive_ok_len",
			op:   validation.NotEqual,
			act:  types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: one}, types.String{Value: two}, types.String{Value: three}}},
			tgt:  targetOneTwo,
		},
		{
			name:        "strings_list_neq_sensitive_nok",
			op:          validation.NotEqual,
			act:         types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: one}, types.String{Value: two}}},
			tgt:         targetOneTwo,
			expectError: true,
		},

		// list []string insensitive neq
		{
			name: "strings_list_neq_insensitive_ok_order",
			op:   validation.NotEqual,
			act:  types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: two}, types.String{Value: one}}},
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name: "strings_list_neq_insensitive_ok_len",
			op:   validation.NotEqual,
			act:  types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: one}, types.String{Value: two}, types.String{Value: three}}},
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name:        "strings_list_neq_insensitive_nok_same",
			op:          validation.NotEqual,
			act:         types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: one}, types.String{Value: two}}},
			tgt:         targetOneTwo,
			meta:        []interface{}{true},
			expectError: true,
		},
		{
			name:        "strings_list_neq_insensitive_nok_casing",
			op:          validation.NotEqual,
			act:         types.List{ElemType: types.StringType, Elems: []attr.Value{types.String{Value: oNe}, types.String{Value: twO}}},
			tgt:         targetOneTwo,
			meta:        []interface{}{true},
			expectError: true,
		},

		// string sensitive oneof
		{
			name: "strings_oneof_sensitive_ok_first",
			op:   validation.OneOf,
			act:  types.String{Value: one},
			tgt:  targetOneTwo,
		},
		{
			name: "strings_oneof_sensitive_ok_last",
			op:   validation.OneOf,
			act:  types.String{Value: two},
			tgt:  targetOneTwo,
		},
		{
			name:        "strings_oneof_sensitive_nok_extra",
			op:          validation.OneOf,
			act:         types.String{Value: three},
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "strings_oneof_sensitive_nok_casing",
			op:          validation.OneOf,
			act:         types.String{Value: oNe},
			tgt:         targetOneTwo,
			expectError: true,
		},

		// string insensitive oneof
		{
			name: "strings_oneof_insensitive_ok_first",
			op:   validation.OneOf,
			act:  types.String{Value: one},
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name: "strings_oneof_insensitive_ok_first_casing",
			op:   validation.OneOf,
			act:  types.String{Value: oNe},
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name: "strings_oneof_insensitive_ok_last",
			op:   validation.OneOf,
			act:  types.String{Value: two},
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name: "strings_oneof_insensitive_ok_last_casing",
			op:   validation.OneOf,
			act:  types.String{Value: twO},
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name:        "strings_oneof_insensitive_nok",
			op:          validation.OneOf,
			act:         types.String{Value: three},
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},

		// string sensitive notoneof
		{
			name: "strings_notoneof_sensitive_ok_extra",
			op:   validation.NotOneOf,
			act:  types.String{Value: three},
			tgt:  targetOneTwo,
		},
		{
			name: "strings_notoneof_sensitive_ok_casing",
			op:   validation.NotOneOf,
			act:  types.String{Value: oNe},
			tgt:  targetOneTwo,
		},
		{
			name:        "strings_notoneof_sensitive_nok_first",
			op:          validation.NotOneOf,
			act:         types.String{Value: one},
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "strings_notoneof_sensitive_nok_last",
			op:          validation.NotOneOf,
			act:         types.String{Value: two},
			tgt:         targetOneTwo,
			expectError: true,
		},

		// string insensitive notoneof
		{
			name: "strings_notoneof_insensitive_ok",
			op:   validation.NotOneOf,
			act:  types.String{Value: three},
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name:        "strings_notoneof_insensitive_nok_first",
			op:          validation.NotOneOf,
			act:         types.String{Value: one},
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},
		{
			name:        "strings_notoneof_insensitive_nok_first_casing",
			op:          validation.NotOneOf,
			act:         types.String{Value: oNe},
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},
		{
			name:        "strings_notoneof_insensitive_nok_last",
			op:          validation.NotOneOf,
			act:         types.String{Value: two},
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},
		{
			name:        "strings_notoneof_insensitive_nok_last_casing",
			op:          validation.NotOneOf,
			act:         types.String{Value: twO},
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},
	}

	for _, ct := range theTests {
		t.Run(ct.name, func(t *testing.T) {
			ct.do(t)
		})
	}
}
