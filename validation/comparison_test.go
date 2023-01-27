package validation_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/validation"
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
	err := validation.CompareAttrValues(context.Background(), ct.act, ct.op, ct.tgt, ct.meta...)
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
			name: "eq_ok",
			op:   validation.Equal,
			tgt:  true,
			act:  types.BoolValue(true),
		},
		{
			name:        "eq_nok",
			op:          validation.Equal,
			tgt:         true,
			act:         types.BoolValue(false),
			expectError: true,
		},
		{
			name: "neq_ok",
			op:   validation.NotEqual,
			tgt:  true,
			act:  types.BoolValue(false),
		},
		{
			name:        "neq_nok",
			op:          validation.NotEqual,
			tgt:         true,
			act:         types.BoolValue(true),
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
			name: "eq_ok",
			op:   validation.Equal,
			tgt:  float64(0.1),
			act:  types.Float64Value(0.1),
		},
		{
			name:        "eq_nok",
			op:          validation.Equal,
			tgt:         float64(1.0),
			act:         types.Float64Value(0.1),
			expectError: true,
		},
		{
			name: "neq_ok",
			op:   validation.NotEqual,
			tgt:  1.0,
			act:  types.Float64Value(0.1),
		},
		{
			name:        "neq_nok",
			op:          validation.NotEqual,
			tgt:         float64(1.0),
			act:         types.Float64Value(1.0),
			expectError: true,
		},
		{
			name: "lt_ok",
			op:   validation.LessThan,
			tgt:  float64(1.0),
			act:  types.Float64Value(0.9),
		},
		{
			name:        "lt_nok",
			op:          validation.LessThan,
			tgt:         float64(1.0),
			act:         types.Float64Value(1.0),
			expectError: true,
		},
		{
			name: "lte_lt_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  float64(1.0),
			act:  types.Float64Value(0.9),
		},
		{
			name:        "lte_lt_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         float64(1.0),
			act:         types.Float64Value(1.1),
			expectError: true,
		},
		{
			name: "lte_eq_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  float64(1.0),
			act:  types.Float64Value(1.0),
		},
		{
			name:        "lte_eq_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         float64(1.0),
			act:         types.Float64Value(1.1),
			expectError: true,
		},
		{
			name: "gt_ok",
			op:   validation.GreaterThan,
			tgt:  float64(1.0),
			act:  types.Float64Value(1.1),
		},
		{
			name:        "gt_nok",
			op:          validation.GreaterThan,
			tgt:         float64(1.0),
			act:         types.Float64Value(1.0),
			expectError: true,
		},
		{
			name: "gte_gt_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  float64(1.0),
			act:  types.Float64Value(1.1),
		},
		{
			name:        "gte_gt_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         float64(1.0),
			act:         types.Float64Value(0.9),
			expectError: true,
		},
		{
			name: "gte_eq_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  float64(1.0),
			act:  types.Float64Value(1.0),
		},
		{
			name:        "gte_eq_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         float64(1.0),
			act:         types.Float64Value(0.9),
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
			name: "eq_ok",
			op:   validation.Equal,
			tgt:  1,
			act:  types.Int64Value(1),
		},
		{
			name:        "eq_nok",
			op:          validation.Equal,
			tgt:         1,
			act:         types.Int64Value(0),
			expectError: true,
		},
		{
			name: "neq_ok",
			op:   validation.NotEqual,
			tgt:  1,
			act:  types.Int64Value(0),
		},
		{
			name:        "neq_nok",
			op:          validation.NotEqual,
			tgt:         1,
			act:         types.Int64Value(1),
			expectError: true,
		},
		{
			name: "lt_ok",
			op:   validation.LessThan,
			tgt:  1,
			act:  types.Int64Value(0),
		},
		{
			name:        "lt_nok",
			op:          validation.LessThan,
			tgt:         1,
			act:         types.Int64Value(1),
			expectError: true,
		},
		{
			name: "lte_lt_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  1,
			act:  types.Int64Value(0),
		},
		{
			name:        "lte_lt_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         1,
			act:         types.Int64Value(2),
			expectError: true,
		},
		{
			name: "lte_eq_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  1,
			act:  types.Int64Value(1),
		},
		{
			name:        "lte_eq_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         1,
			act:         types.Int64Value(2),
			expectError: true,
		},
		{
			name: "gt_ok",
			op:   validation.GreaterThan,
			tgt:  1,
			act:  types.Int64Value(2),
		},
		{
			name:        "gt_nok",
			op:          validation.GreaterThan,
			tgt:         1,
			act:         types.Int64Value(1),
			expectError: true,
		},
		{
			name: "gte_gt_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  1,
			act:  types.Int64Value(2),
		},
		{
			name:        "gte_gt_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         1,
			act:         types.Int64Value(0),
			expectError: true,
		},
		{
			name: "gte_eq_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  1,
			act:  types.Int64Value(1),
		},
		{
			name:        "gte_eq_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         1,
			act:         types.Int64Value(0),
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
			name: "eq_ok",
			op:   validation.Equal,
			tgt:  big.NewFloat(0.1),
			act:  types.NumberValue(big.NewFloat(0.1)),
		},
		{
			name:        "eq_nok",
			op:          validation.Equal,
			tgt:         big.NewFloat(1.0),
			act:         types.NumberValue(big.NewFloat(0.1)),
			expectError: true,
		},
		{
			name: "neq_ok",
			op:   validation.NotEqual,
			tgt:  1.0,
			act:  types.NumberValue(big.NewFloat(0.1)),
		},
		{
			name:        "neq_nok",
			op:          validation.NotEqual,
			tgt:         big.NewFloat(1.0),
			act:         types.NumberValue(big.NewFloat(1.0)),
			expectError: true,
		},
		{
			name: "lt_ok",
			op:   validation.LessThan,
			tgt:  big.NewFloat(1.0),
			act:  types.NumberValue(big.NewFloat(0.9)),
		},
		{
			name:        "lt_nok",
			op:          validation.LessThan,
			tgt:         big.NewFloat(1.0),
			act:         types.NumberValue(big.NewFloat(1.0)),
			expectError: true,
		},
		{
			name: "lte_lt_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  big.NewFloat(1.0),
			act:  types.NumberValue(big.NewFloat(0.9)),
		},
		{
			name:        "lte_lt_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         big.NewFloat(1.0),
			act:         types.NumberValue(big.NewFloat(1.1)),
			expectError: true,
		},
		{
			name: "lte_eq_ok",
			op:   validation.LessThanOrEqualTo,
			tgt:  big.NewFloat(1.0),
			act:  types.NumberValue(big.NewFloat(1.0)),
		},
		{
			name:        "lte_eq_nok",
			op:          validation.LessThanOrEqualTo,
			tgt:         big.NewFloat(1.0),
			act:         types.NumberValue(big.NewFloat(1.1)),
			expectError: true,
		},
		{
			name: "gt_ok",
			op:   validation.GreaterThan,
			tgt:  big.NewFloat(1.0),
			act:  types.NumberValue(big.NewFloat(1.1)),
		},
		{
			name:        "gt_nok",
			op:          validation.GreaterThan,
			tgt:         big.NewFloat(1.0),
			act:         types.NumberValue(big.NewFloat(1.0)),
			expectError: true,
		},
		{
			name: "gte_gt_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  big.NewFloat(1.0),
			act:  types.NumberValue(big.NewFloat(1.1)),
		},
		{
			name:        "gte_gt_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         big.NewFloat(1.0),
			act:         types.NumberValue(big.NewFloat(0.9)),
			expectError: true,
		},
		{
			name: "gte_eq_ok",
			op:   validation.GreaterThanOrEqualTo,
			tgt:  big.NewFloat(1.0),
			act:  types.NumberValue(big.NewFloat(1.0)),
		},
		{
			name:        "gte_eq_nok",
			op:          validation.GreaterThanOrEqualTo,
			tgt:         big.NewFloat(1.0),
			act:         types.NumberValue(big.NewFloat(0.9)),
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
			name: "eq_ok",
			op:   validation.Equal,
			act:  types.StringValue("hi"),
			tgt:  "hi",
		},
		{
			name:        "eq_nok",
			op:          validation.Equal,
			act:         types.StringValue("hi"),
			tgt:         "ih",
			expectError: true,
		},
		{
			name: "neq_ok",
			op:   validation.NotEqual,
			act:  types.StringValue("hi"),
			tgt:  "ih",
		},
		{
			name:        "neq_nok",
			op:          validation.NotEqual,
			act:         types.StringValue("hi"),
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
			name: "list_eq_sensitive_ok",
			op:   validation.Equal,
			act:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue(one), types.StringValue(two)}),
			tgt:  targetOneTwo,
		},
		{
			name:        "list_eq_sensitive_nok_casing",
			op:          validation.Equal,
			act:         types.ListValueMust(types.StringType, []attr.Value{types.StringValue(oNe), types.StringValue(twO)}),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "list_eq_sensitive_nok_order",
			op:          validation.Equal,
			act:         types.ListValueMust(types.StringType, []attr.Value{types.StringValue(two), types.StringValue(one)}),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "list_eq_sensitive_nok_len",
			op:          validation.Equal,
			act:         types.ListValueMust(types.StringType, []attr.Value{types.StringValue(one), types.StringValue(two), types.StringValue(three)}),
			tgt:         targetOneTwo,
			expectError: true,
		},

		// list []string insensitive eq
		{
			name: "list_eq_insensitive_ok",
			op:   validation.Equal,
			act:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue(one), types.StringValue(two)}),
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name: "list_eq_insensitive_ok_casing",
			op:   validation.Equal,
			act:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue(oNe), types.StringValue(twO)}),
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name:        "list_eq_insensitive_nok_order",
			op:          validation.Equal,
			act:         types.ListValueMust(types.StringType, []attr.Value{types.StringValue(two), types.StringValue(one)}),
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},
		{
			name:        "list_eq_insensitive_nok_len",
			op:          validation.Equal,
			act:         types.ListValueMust(types.StringType, []attr.Value{types.StringValue(one), types.StringValue(two), types.StringValue(three)}),
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},

		// list []string sensitive neq
		{
			name: "list_neq_sensitive_ok_order",
			op:   validation.NotEqual,
			act:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue(two), types.StringValue(one)}),
			tgt:  targetOneTwo,
		},
		{
			name: "list_neq_sensitive_ok_casing",
			op:   validation.NotEqual,
			act:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue(oNe), types.StringValue(twO)}),
			tgt:  targetOneTwo,
		},
		{
			name: "list_neq_sensitive_ok_len",
			op:   validation.NotEqual,
			act:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue(one), types.StringValue(two), types.StringValue(three)}),
			tgt:  targetOneTwo,
		},
		{
			name:        "list_neq_sensitive_nok",
			op:          validation.NotEqual,
			act:         types.ListValueMust(types.StringType, []attr.Value{types.StringValue(one), types.StringValue(two)}),
			tgt:         targetOneTwo,
			expectError: true,
		},

		// list []string insensitive neq
		{
			name: "list_neq_insensitive_ok_order",
			op:   validation.NotEqual,
			act:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue(two), types.StringValue(one)}),
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name: "list_neq_insensitive_ok_len",
			op:   validation.NotEqual,
			act:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue(one), types.StringValue(two), types.StringValue(three)}),
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name:        "list_neq_insensitive_nok_same",
			op:          validation.NotEqual,
			act:         types.ListValueMust(types.StringType, []attr.Value{types.StringValue(one), types.StringValue(two)}),
			tgt:         targetOneTwo,
			meta:        []interface{}{true},
			expectError: true,
		},
		{
			name:        "list_neq_insensitive_nok_casing",
			op:          validation.NotEqual,
			act:         types.ListValueMust(types.StringType, []attr.Value{types.StringValue(oNe), types.StringValue(twO)}),
			tgt:         targetOneTwo,
			meta:        []interface{}{true},
			expectError: true,
		},

		// string sensitive oneof
		{
			name: "oneof_sensitive_ok_first",
			op:   validation.OneOf,
			act:  types.StringValue(one),
			tgt:  targetOneTwo,
		},
		{
			name: "oneof_sensitive_ok_last",
			op:   validation.OneOf,
			act:  types.StringValue(two),
			tgt:  targetOneTwo,
		},
		{
			name:        "oneof_sensitive_nok_extra",
			op:          validation.OneOf,
			act:         types.StringValue(three),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "oneof_sensitive_nok_casing",
			op:          validation.OneOf,
			act:         types.StringValue(oNe),
			tgt:         targetOneTwo,
			expectError: true,
		},

		// string insensitive oneof
		{
			name: "oneof_insensitive_ok_first",
			op:   validation.OneOf,
			act:  types.StringValue(one),
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name: "oneof_insensitive_ok_first_casing",
			op:   validation.OneOf,
			act:  types.StringValue(oNe),
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name: "oneof_insensitive_ok_last",
			op:   validation.OneOf,
			act:  types.StringValue(two),
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name: "oneof_insensitive_ok_last_casing",
			op:   validation.OneOf,
			act:  types.StringValue(twO),
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name:        "oneof_insensitive_nok",
			op:          validation.OneOf,
			act:         types.StringValue(three),
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},

		// string sensitive notoneof
		{
			name: "notoneof_sensitive_ok_extra",
			op:   validation.NotOneOf,
			act:  types.StringValue(three),
			tgt:  targetOneTwo,
		},
		{
			name: "notoneof_sensitive_ok_casing",
			op:   validation.NotOneOf,
			act:  types.StringValue(oNe),
			tgt:  targetOneTwo,
		},
		{
			name:        "notoneof_sensitive_nok_first",
			op:          validation.NotOneOf,
			act:         types.StringValue(one),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "notoneof_sensitive_nok_last",
			op:          validation.NotOneOf,
			act:         types.StringValue(two),
			tgt:         targetOneTwo,
			expectError: true,
		},

		// string insensitive notoneof
		{
			name: "notoneof_insensitive_ok",
			op:   validation.NotOneOf,
			act:  types.StringValue(three),
			tgt:  targetOneTwo,
			meta: []interface{}{true},
		},
		{
			name:        "notoneof_insensitive_nok_first",
			op:          validation.NotOneOf,
			act:         types.StringValue(one),
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},
		{
			name:        "notoneof_insensitive_nok_first_casing",
			op:          validation.NotOneOf,
			act:         types.StringValue(oNe),
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},
		{
			name:        "notoneof_insensitive_nok_last",
			op:          validation.NotOneOf,
			act:         types.StringValue(two),
			tgt:         targetOneTwo,
			expectError: true,
			meta:        []interface{}{true},
		},
		{
			name:        "notoneof_insensitive_nok_last_casing",
			op:          validation.NotOneOf,
			act:         types.StringValue(twO),
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

func TestComparison_Ints(t *testing.T) {
	var (
		attrInt1 = types.Int64Value(1)
		attrInt2 = types.Int64Value(2)
		attrInt3 = types.Int64Value(3)

		attrNum1 = types.NumberValue(big.NewFloat(float64(1)))
		attrNum2 = types.NumberValue(big.NewFloat(float64(2)))
		attrNum3 = types.NumberValue(big.NewFloat(float64(3)))

		targetOneTwo = []int{1, 2}
	)

	theTests := []comparisonTest{
		// list int64
		{
			name: "list_int64_eq_ok",
			op:   validation.Equal,
			act:  types.ListValueMust(types.Int64Type, []attr.Value{attrInt1, attrInt2}),
			tgt:  targetOneTwo,
		},
		{
			name:        "list_int64_eq_nok_order",
			op:          validation.Equal,
			act:         types.ListValueMust(types.Int64Type, []attr.Value{attrInt2, attrInt1}),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "list_int64_eq_nok_extra",
			op:          validation.Equal,
			act:         types.ListValueMust(types.Int64Type, []attr.Value{attrInt1, attrInt2, attrInt3}),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name: "list_int64_neq_ok_order",
			op:   validation.NotEqual,
			act:  types.ListValueMust(types.Int64Type, []attr.Value{attrInt2, attrInt1}),
			tgt:  targetOneTwo,
		},
		{
			name: "list_int64_eq_ok_extra",
			op:   validation.NotEqual,
			act:  types.ListValueMust(types.Int64Type, []attr.Value{attrInt1, attrInt2, attrInt3}),
			tgt:  targetOneTwo,
		},
		{
			name:        "list_int64_neq_nok",
			op:          validation.NotEqual,
			act:         types.ListValueMust(types.Int64Type, []attr.Value{attrInt1, attrInt2}),
			tgt:         targetOneTwo,
			expectError: true,
		},

		// list number
		{
			name: "list_number_eq_ok",
			op:   validation.Equal,
			act:  types.ListValueMust(types.NumberType, []attr.Value{attrNum1, attrNum2}),
			tgt:  targetOneTwo,
		},
		{
			name:        "list_number_eq_nok_order",
			op:          validation.Equal,
			act:         types.ListValueMust(types.NumberType, []attr.Value{attrNum2, attrNum1}),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "list_number_eq_nok_extra",
			op:          validation.Equal,
			act:         types.ListValueMust(types.NumberType, []attr.Value{attrNum1, attrNum2, attrNum3}),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name: "list_number_neq_ok_order",
			op:   validation.NotEqual,
			act:  types.ListValueMust(types.NumberType, []attr.Value{attrNum2, attrNum1}),
			tgt:  targetOneTwo,
		},
		{
			name: "list_number_neq_ok_extra",
			op:   validation.NotEqual,
			act:  types.ListValueMust(types.NumberType, []attr.Value{attrNum1, attrNum2, attrNum3}),
			tgt:  targetOneTwo,
		},
		{
			name:        "list_number_neq_nok",
			op:          validation.NotEqual,
			act:         types.ListValueMust(types.NumberType, []attr.Value{attrNum1, attrNum2}),
			tgt:         targetOneTwo,
			expectError: true,
		},

		// set int64
		{
			name: "set_int64_eq_ok",
			op:   validation.Equal,
			act:  types.SetValueMust(types.Int64Type, []attr.Value{attrInt1, attrInt2}),
			tgt:  targetOneTwo,
		},
		{
			name:        "set_int64_eq_nok_order",
			op:          validation.Equal,
			act:         types.SetValueMust(types.Int64Type, []attr.Value{attrInt2, attrInt1}),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "set_int64_eq_nok_extra",
			op:          validation.Equal,
			act:         types.SetValueMust(types.Int64Type, []attr.Value{attrInt1, attrInt2, attrInt3}),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name: "set_int64_neq_ok_order",
			op:   validation.NotEqual,
			act:  types.SetValueMust(types.Int64Type, []attr.Value{attrInt2, attrInt1}),
			tgt:  targetOneTwo,
		},
		{
			name: "set_int64_eq_ok_extra",
			op:   validation.NotEqual,
			act:  types.SetValueMust(types.Int64Type, []attr.Value{attrInt1, attrInt2, attrInt3}),
			tgt:  targetOneTwo,
		},
		{
			name:        "set_int64_neq_nok",
			op:          validation.NotEqual,
			act:         types.SetValueMust(types.Int64Type, []attr.Value{attrInt1, attrInt2}),
			tgt:         targetOneTwo,
			expectError: true,
		},

		// set number
		{
			name: "set_number_eq_ok",
			op:   validation.Equal,
			act:  types.SetValueMust(types.NumberType, []attr.Value{attrNum1, attrNum2}),
			tgt:  targetOneTwo,
		},
		{
			name:        "set_number_eq_nok_order",
			op:          validation.Equal,
			act:         types.SetValueMust(types.NumberType, []attr.Value{attrNum2, attrNum1}),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "set_number_eq_nok_extra",
			op:          validation.Equal,
			act:         types.SetValueMust(types.NumberType, []attr.Value{attrNum1, attrNum2, attrNum3}),
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name: "set_number_neq_ok_order",
			op:   validation.NotEqual,
			act:  types.SetValueMust(types.NumberType, []attr.Value{attrNum2, attrNum1}),
			tgt:  targetOneTwo,
		},
		{
			name: "set_number_neq_ok_extra",
			op:   validation.NotEqual,
			act:  types.SetValueMust(types.NumberType, []attr.Value{attrNum1, attrNum2, attrNum3}),
			tgt:  targetOneTwo,
		},
		{
			name:        "set_number_neq_nok",
			op:          validation.NotEqual,
			act:         types.SetValueMust(types.NumberType, []attr.Value{attrNum1, attrNum2}),
			tgt:         targetOneTwo,
			expectError: true,
		},

		// int64
		{
			name: "int64_oneof_ok_first",
			op:   validation.OneOf,
			act:  attrInt1,
			tgt:  targetOneTwo,
		},
		{
			name: "int64_oneof_ok_last",
			op:   validation.OneOf,
			act:  attrInt2,
			tgt:  targetOneTwo,
		},
		{
			name:        "int64_oneof_nok",
			op:          validation.OneOf,
			act:         attrInt3,
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name: "int64_notoneof_ok",
			op:   validation.NotOneOf,
			act:  attrInt3,
			tgt:  targetOneTwo,
		},
		{
			name:        "int64_notoneof_nok_first",
			op:          validation.NotOneOf,
			act:         attrInt1,
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "int64_notoneof_nok_last",
			op:          validation.NotOneOf,
			act:         attrInt2,
			tgt:         targetOneTwo,
			expectError: true,
		},

		// number
		{
			name: "number_oneof_ok_first",
			op:   validation.OneOf,
			act:  attrNum1,
			tgt:  targetOneTwo,
		},
		{
			name: "number_oneof_ok_last",
			op:   validation.OneOf,
			act:  attrNum2,
			tgt:  targetOneTwo,
		},
		{
			name:        "number_oneof_nok",
			op:          validation.OneOf,
			act:         attrNum3,
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name: "number_notoneof_ok",
			op:   validation.NotOneOf,
			act:  attrNum3,
			tgt:  targetOneTwo,
		},
		{
			name:        "number_notoneof_nok_first",
			op:          validation.NotOneOf,
			act:         attrNum1,
			tgt:         targetOneTwo,
			expectError: true,
		},
		{
			name:        "number_notoneof_nok_last",
			op:          validation.NotOneOf,
			act:         attrNum2,
			tgt:         targetOneTwo,
			expectError: true,
		},
	}

	for _, ct := range theTests {
		t.Run(ct.name, func(t *testing.T) {
			ct.do(t)
		})
	}
}
