package acctest_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/dcarbone/terraform-plugin-framework-utils/acctest"
)

func TestConfigValue_Defaults(t *testing.T) {
	type convTest struct {
		in  interface{}
		out interface{}
	}

	theTests := []convTest{
		{
			in:  nil,
			out: "null",
		},
		{
			in:  acctest.ConfigLiteral("local.whateverthing"),
			out: "local.whateverthing",
		},
		{
			in:  true,
			out: "true",
		},
		{
			in:  false,
			out: "false",
		},
		{
			in:  math.MaxInt16,
			out: fmt.Sprintf("%d", math.MaxInt16),
		},
		{
			in:  -5,
			out: "-5",
		},
		{
			in:  0.5,
			out: fmt.Sprintf("%f", 0.5),
		},
		{
			in:  "single-line",
			out: `"single-line"`,
		},
		{
			in: `multi
line`,
			out: `<<EOD
multi
line
EOD`,
		},
		{
			in: []interface{}{"hello", 5},
			out: `[
"hello",
5
]`,
		},
		{
			in: []string{"hello", "there"},
			out: `[
"hello",
"there"
]`,
		},
		{
			in: []int{1, 5},
			out: `[
1,
5
]`,
		},
		{
			in: []float64{1, 5},
			out: `[
` + fmt.Sprintf("%f", float64(1)) + `,
` + fmt.Sprintf("%f", float64(5)) + `
]`,
		},
	}

	// todo: use biggerer brain to figure out how to test map -> string verification

	for _, theT := range theTests {
		out := acctest.ConfigValue(theT.in)
		if !reflect.DeepEqual(out, theT.out) {
			t.Log("Output does not match expected")
			t.Logf("Input: %v", theT.in)
			t.Logf("Expected: %v", theT.out)
			t.Logf("Actual: %v", out)
			t.Fail()
		}
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
