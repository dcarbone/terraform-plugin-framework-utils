package validation_test

import (
	"testing"

	"github.com/dcarbone/terraform-plugin-framework-utils/validation"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type comparisonTest struct {
	name        string
	op          validation.CompareOp
	in          attr.Value
	tgt         interface{}
	expectPanic bool
	expectError bool
}

func TestComparison_String(t *testing.T) {
	theTests := []comparisonTest{
		{
			name: "string_eq_ok",
			op:   validation.Equal,
			in:   types.String{Value: "hi"},
			tgt:  "hi",
		},
		{
			name:        "string_eq_nok",
			op:          validation.Equal,
			in:          types.String{Value: "hi"},
			tgt:         "ih",
			expectError: true,
		},
		{
			name: "string_neq_ok",
			op:   validation.NotEqual,
			in:   types.String{Value: "hi"},
			tgt:  "ih",
		},
		{
			name:        "string_neq_nok",
			op:          validation.NotEqual,
			in:          types.String{Value: "hi"},
			tgt:         "hi",
			expectError: true,
		},
	}

	for _, at := range theTests {
		t.Run(at.name, func(t *testing.T) {
			if at.expectPanic {
				defer func() {
					if v := recover(); v != nil {
						t.Logf("Panic seen: %v", v)
					} else {
						t.Log("panic expected")
						t.Fail()
					}
				}()
			}
			err := validation.CompareAttrValues(at.in, at.op, at.tgt)
			if err != nil {
				t.Logf("comparison error seen: %v", err)
				if !at.expectError {
					t.Fail()
				}
			}
		})
	}
}
