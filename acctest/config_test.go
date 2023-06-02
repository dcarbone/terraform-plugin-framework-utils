package acctest_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
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
EOD
`,
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
		{
			name: "slice-map-string-interface",
			in: []map[string]interface{}{
				{"key1": []string{"k1v1", "k1v2"}, "key2": "k2v2"},
				{"key3": 5},
			},
			out: `[
		{
		key1 = [
		"k1v1",
		"k1v2"
		]
		key2 = "k2v2"
		},
		{
		key3 = 5
		}
		]`,
		},
	}

	for _, theT := range theTests {
		t.Run(theT.name, func(t *testing.T) {
			var (
				expectedHCLFile *hcl.File
				actualHCLBlock  string
				actualHCLFile   *hcl.File
				diags           hcl.Diagnostics
				err             error

				expectedHCLName  = fmt.Sprintf("%s.hcl", theT.name)
				expectedHCLBlock = fmt.Sprintf("testvar = %s", theT.out)
				expectedData     = make(map[string]interface{})

				actualHCLName = fmt.Sprintf("%s.hcl", theT.name)
				actualData    = make(map[string]interface{})

				// create hcl parser to use diags from to test both expected and actual output
				hp = hclparse.NewParser()
			)

			// try to parse the expected output, making sure our test is valid
			if expectedHCLFile, diags = hp.ParseHCL([]byte(expectedHCLBlock), expectedHCLName); diags.HasErrors() {
				t.Logf("Expected output HCL contains errors: %v", diags.Error())
				t.Log(expectedHCLBlock)
				t.Fail()
				return
			}

			// attempt to decode expected data
			if diags = gohcl.DecodeBody(expectedHCLFile.Body, nil, &expectedData); diags.HasErrors() {
				t.Logf("Error decoding expected HCL: %v", err)
				t.Log(expectedHCLBlock)
				t.Fail()
				return
			}

			// turn output into a block definition for the parser
			actualHCLBlock = fmt.Sprintf("testvar = %s", acctest.ConfigValue(theT.in))

			// attempt parse and check diagnostics for errors
			if actualHCLFile, diags = hp.ParseHCL([]byte(actualHCLBlock), actualHCLName); diags.HasErrors() {
				t.Logf("Failed to parse generated HCL: %v", diags.Error())
				t.Log(actualHCLBlock)
				t.Fail()
				return
			}

			// decode generated hcl
			if diags = gohcl.DecodeBody(actualHCLFile.Body, nil, &actualData); diags.HasErrors() {
				t.Logf("Failed to decode generated HCL: %v", err)
				t.Fail()
				return
			}

			assert.EqualValues(t, expectedData, actualData, "Actual decoded HCL does not match expected")
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
