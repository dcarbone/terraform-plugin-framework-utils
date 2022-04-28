package acctest_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/dcarbone/terraform-plugin-framework-utils/acctest"
)

func TestConfigValue_Defaults(t *testing.T) {
	type convTest struct {
		name string
		in   interface{}
		out  interface{}
	}

	theTests := []convTest{
		{
			name: "nil",
			in:   nil,
			out:  "null",
		},
		{
			name: "config-literal",
			in:   acctest.ConfigLiteral("local.whateverthing"),
			out:  "local.whateverthing",
		},
		{
			name: "bool-true",
			in:   true,
			out:  "true",
		},
		{
			name: "bool-false",
			in:   false,
			out:  "false",
		},
		{
			name: "int-positive",
			in:   math.MaxInt16,
			out:  fmt.Sprintf("%d", math.MaxInt16),
		},
		{
			name: "int-negative",
			in:   -5,
			out:  "-5",
		},
		{
			name: "float-positive",
			in:   0.5,
			out:  fmt.Sprintf("%f", 0.5),
		},
		{
			name: "float-negative",
			in:   -0.5,
			out:  fmt.Sprintf("%f", -0.5),
		},
		{
			name: "string-literal",
			in:   "single-line",
			out:  `"single-line"`,
		},
		{
			name: "string-heredoc",
			in: `multi
line`,
			out: `<<EOD
multi
line
EOD`,
		},
		{
			name: "duration-to-string",
			in:   time.Nanosecond,
			out:  `"1ns"`,
		},
		{
			name: "slice-interface",
			in:   []interface{}{"hello", 5},
			out: `[
"hello",
5
]`,
		},
		{
			name: "slice-string",
			in:   []string{"hello", "there"},
			out: `[
"hello",
"there"
]`,
		},
		{
			name: "slice-int",
			in:   []int{1, 5},
			out: `[
1,
5
]`,
		},
		{
			name: "slice-float64",
			in:   []float64{1, 5},
			out: `[
` + fmt.Sprintf("%f", float64(1)) + `,
` + fmt.Sprintf("%f", float64(5)) + `
]`,
		},
	}

	// todo: use biggerer brain to figure out how to test map -> string verification

	for _, theT := range theTests {
		t.Run(theT.name, func(t *testing.T) {
			out := acctest.ConfigValue(theT.in)
			if !reflect.DeepEqual(out, theT.out) {
				t.Log("Output does not match expected")
				t.Logf("Input: %v", theT.in)
				t.Logf("Expected: %v", theT.out)
				t.Logf("Actual: %v", out)
				t.Fail()
			}
		})
	}
}

func TestConfigValue_Set(t *testing.T) {
	int32fn := func(v interface{}) string {
		return fmt.Sprintf("%d", v.(int32))
	}
	acctest.SetConfigValueFunc(int32(0), int32fn)
	v := acctest.ConfigValue(int32(1))
	if v != "1" {
		t.Logf("Expected int32(1) to produce 1, saw %v", v)
		t.Fail()
	}
}
