package conv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FormatPathPathSteps takes one or more path steps and joins them together with "."
func FormatPathPathSteps(pathSteps ...path.PathStep) string {
	bits := make([]string, 0)
	for _, pathStep := range pathSteps {
		bits = append(bits, pathStep.String())
	}
	return strings.Join(bits, ".")
}

// FormatPathPaths takes one or more path.Path types and returns a pretty-printable string.
func FormatPathPaths(paths ...path.Path) string {
	out := "["
	for i, o := range paths {
		if i > 0 {
			out = fmt.Sprintf("%s, ", out)
		}
		out = fmt.Sprintf("%s%q", out, FormatPathPathSteps(o.Steps()...))
	}
	return fmt.Sprintf("%s]", out)
}

// FormatAttributePathSteps takes one or more path steps and joins them together with "."
func FormatAttributePathSteps(pathSteps ...tftypes.AttributePathStep) string {
	bits := make([]string, 0)
	for _, pathStep := range pathSteps {
		switch pathStep.(type) {
		case tftypes.AttributeName:
			bits = append(bits, string(pathStep.(tftypes.AttributeName)))
		case tftypes.ElementKeyString:
			bits = append(bits, string(pathStep.(tftypes.ElementKeyString)))
		case tftypes.ElementKeyInt:
			bits = append(bits, strconv.FormatInt(int64(pathStep.(tftypes.ElementKeyInt)), 10))
		case tftypes.ElementKeyValue:
			bits = append(bits, (tftypes.Value)(pathStep.(tftypes.ElementKeyValue)).String())

		default:
			// if this is reached, a new path step implementation has been created
			panic(fmt.Sprintf("no case to convert type %T (%[1]v) to string, please create issue with this error message", pathStep))
		}
	}
	return strings.Join(bits, ".")
}

// FormatAttributePaths takes one or more *tftypes.AttributePaths and returns a pretty-printable string.
func FormatAttributePaths(paths ...*tftypes.AttributePath) string {
	out := "["
	for i, o := range paths {
		if i > 0 {
			out = fmt.Sprintf("%s, ", out)
		}
		out = fmt.Sprintf("%s%q", out, FormatAttributePathSteps(o.Steps()...))
	}
	return fmt.Sprintf("%s]", out)
}
