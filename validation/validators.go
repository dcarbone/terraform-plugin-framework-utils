package validation

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/conv"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TestFunc func(context.Context, GenericRequest, *GenericResponse)

// GenericConfig describes the configuration of a Generic
type GenericConfig struct {
	Description         string
	MarkdownDescription string
	Describer           validator.Describer
	TestFunc            TestFunc
	SkipWhenNull        bool
	SkipWhenUnknown     bool
}

// Generic  a validator that can be applied to variable attribute types
type Generic struct {
	validator.Describer

	d           string
	md          string
	fn          TestFunc
	skipNull    bool
	skipUnknown bool
}

// NewGenericValidator returns a type implementing every validator interface
func NewGenericValidator(conf GenericConfig) Generic {
	if conf.TestFunc == nil {
		panic("test function cannot be nil")
	}
	v := Generic{
		Describer:   conf.Describer,
		d:           conf.Description,
		md:          conf.MarkdownDescription,
		fn:          conf.TestFunc,
		skipNull:    conf.SkipWhenNull,
		skipUnknown: conf.SkipWhenUnknown,
	}
	return v
}

func (g Generic) Validate(ctx context.Context, req GenericRequest, resp *GenericResponse) {
	err := conv.TestAttributeValueState(req.ConfigValue)

	if g.skipUnknown && conv.IsValueIsUnknownError(err) {
		return
	}

	if g.skipNull && conv.IsValueIsNullError(err) {
		return
	}

	// otherwise, fire away!
	g.fn(ctx, req, resp)
}

func (g Generic) Description(ctx context.Context) string {
	if g.Describer != nil {
		return g.Describer.Description(ctx)
	}
	return g.d
}

func (g Generic) MarkdownDescription(ctx context.Context) string {
	if g.Describer != nil {
		return g.Describer.MarkdownDescription(ctx)
	}
	return g.md
}

func (g Generic) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	rq, rp, err := toGenericTypes(req, resp)
	if err != nil {
		resp.Diagnostics.AddError("conversion error", err.Error())
		return
	}
	g.Validate(ctx, rq, rp)
}

func (g Generic) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	rq, rp, err := toGenericTypes(req, resp)
	if err != nil {
		resp.Diagnostics.AddError("conversion error", err.Error())
		return
	}
	g.Validate(ctx, rq, rp)
}

func (g Generic) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	rq, rp, err := toGenericTypes(req, resp)
	if err != nil {
		resp.Diagnostics.AddError("conversion error", err.Error())
		return
	}
	g.Validate(ctx, rq, rp)
}

func (g Generic) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	rq, rp, err := toGenericTypes(req, resp)
	if err != nil {
		resp.Diagnostics.AddError("conversion error", err.Error())
		return
	}
	g.Validate(ctx, rq, rp)
}

func (g Generic) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	rq, rp, err := toGenericTypes(req, resp)
	if err != nil {
		resp.Diagnostics.AddError("conversion error", err.Error())
		return
	}
	g.Validate(ctx, rq, rp)
}

func (g Generic) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	rq, rp, err := toGenericTypes(req, resp)
	if err != nil {
		resp.Diagnostics.AddError("conversion error", err.Error())
		return
	}
	g.Validate(ctx, rq, rp)
}

func (g Generic) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	rq, rp, err := toGenericTypes(req, resp)
	if err != nil {
		resp.Diagnostics.AddError("conversion error", err.Error())
		return
	}
	g.Validate(ctx, rq, rp)
}

func (g Generic) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	rq, rp, err := toGenericTypes(req, resp)
	if err != nil {
		resp.Diagnostics.AddError("conversion error", err.Error())
		return
	}
	g.Validate(ctx, rq, rp)
}

func (g Generic) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	rq, rp, err := toGenericTypes(req, resp)
	if err != nil {
		resp.Diagnostics.AddError("conversion error", err.Error())
		return
	}
	g.Validate(ctx, rq, rp)
}

// RequiredTest is an Generic implementation that will register an error if the attribute has no value of
// any kind
func RequiredTest() TestFunc {
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		if conv.TestAttributeValueState(req.ConfigValue) == nil {
			return
		}

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Attribute must be valued",
			"Attribute must have a value configured",
		)
	}
}

var requiredValidator = NewGenericValidator(GenericConfig{
	Description:         "Asserts the attribute is defined and non-null",
	MarkdownDescription: "Asserts the attribute is defined and non-null",
	TestFunc:            RequiredTest(),
	SkipWhenNull:        false,
	SkipWhenUnknown:     true,
})

// Required returns a validator that asserts a field is configured with any value at all.
func Required() Generic {
	return requiredValidator
}

// RegexpMatchTest is an Generic implementation that will first attempt to convert the value of
// a field to a string, then see if that resulting string matches the provided expression.
func RegexpMatchTest(r string) TestFunc {
	re := regexp.MustCompile(r)
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		str := conv.AttributeValueToString(req.ConfigValue)
		if !re.MatchString(str) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Field value does not match expression",
				fmt.Sprintf("Field value %q does not match expression %q", str, r),
			)
		}
	}
}

// RegexpMatch returns a validator that asserts an attribute's value matches the provided expression
func RegexpMatch(r string) Generic {
	v := NewGenericValidator(GenericConfig{
		Description:         fmt.Sprintf("Asserts attribute string value matches expression %q", r),
		MarkdownDescription: fmt.Sprintf("Asserts attribute string value matches expression %q", r),
		TestFunc:            RegexpMatchTest(r),
		SkipWhenNull:        true,
		SkipWhenUnknown:     true,
	})
	return v
}

// RegexpNotMatchTest is an Generic implementation that will first attempt to convert the value of
// a field to a string, then see if that resulting string matches the provided expression.
func RegexpNotMatchTest(r string) TestFunc {
	re := regexp.MustCompile(r)
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		str := conv.AttributeValueToString(req.ConfigValue)
		if re.MatchString(str) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Field value must NOT match expression",
				fmt.Sprintf("Field value %q matches expression %q, indicating it contains invalid characters", str, r),
			)
		}
	}
}

// RegexpNotMatch returns a validator that asserts an attribute's value does not match the provided expression
func RegexpNotMatch(r string) Generic {
	v := NewGenericValidator(GenericConfig{
		Description:         fmt.Sprintf("Assert attribute string value does not match expression %q", r),
		MarkdownDescription: fmt.Sprintf("Assert attribute string value does not match expression %q", r),
		TestFunc:            RegexpNotMatchTest(r),
		SkipWhenNull:        true,
		SkipWhenUnknown:     true,
	})
	return v
}

// LengthTest is an Generic implementation that attempts to find a length value appropriate
// for the attribute value type, asserting that it is within the specified bounds.
//
// If either min or max is set to -1, then that value is unbounded and thus not verified.
func LengthTest(minL, maxL int) TestFunc {
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		// perform some basic sanity checking
		if minL == -1 && maxL == -1 {
			resp.Diagnostics.AddAttributeWarning(
				req.Path,
				"Length validation is unbounded, there is nothing to verify",
				"Both minL and maxL variables were set to -1.  This has no purpose and should be rectified.",
			)
			return
		} else if minL < -1 || maxL < -1 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Cannot use negative value for length check",
				fmt.Sprintf("The provided minL value of %q and / or the provided maxL value %q are negative."+
					"  The only valid negative value is -1 to indicate \"unbounded\".  This should be rectified.",
					minL,
					maxL,
				))
			return
		} else if maxL != -1 && minL > maxL {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Minimum length value is greater than maximum length value",
				fmt.Sprintf(
					"The provided minimum length %q is greater than the provided maximum length of %q."+
						"  This should be rectified.",
					minL,
					maxL,
				))
			return
		}

		fl := conv.AttributeValueLength(req.ConfigValue)
		if minL > -1 && fl < minL {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Field value length is below minimum threshold",
				fmt.Sprintf("Field value length %d is less than mininum allowed of %d", fl, minL),
			)
		}
		if maxL > -1 && fl > maxL {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Field value length is above maximum threshold",
				fmt.Sprintf("Field value length %d is greater than the maximum allowed of %d", fl, maxL),
			)
		}
	}
}

// Length asserts an attribute's length is limited to the specified bounds.  The allowed types are:
//   - string
//   - set
//   - map
//   - list
func Length(minL, maxL int) Generic {
	v := NewGenericValidator(GenericConfig{
		Description:         fmt.Sprintf("Asserts an attribute's value contains no less than %d and no more than %d elements, with -1 meaning unbounded", minL, maxL),
		MarkdownDescription: fmt.Sprintf("Asserts an attribute's value contains no less than %d and no more than %d elements, with -1 meaning unbounded", minL, maxL),
		TestFunc:            LengthTest(minL, maxL),
		SkipWhenNull:        true,
		SkipWhenUnknown:     true,
	})
	return v
}

// CompareTest executes a registered comparison function against the target attribute's value
func CompareTest(op CompareOp, target interface{}, meta ...interface{}) TestFunc {
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		err := CompareAttrValues(ctx, req.ConfigValue, op, target, meta...)
		if err != nil {
			switch true {
			case errors.Is(err, ErrComparisonFailed):
				addComparisonFailedDiagnostic(op, target, req, resp, err)

			case errors.Is(err, ErrTypeConversionFailed):
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Could not convert attribute to target type for comparison",
					fmt.Sprintf("Unable to convert attribute value type %T to %T for %q copmarison: %v", req, target, op, err))

			default:
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Unexpected error during comparison",
					fmt.Sprintf("Unexpected error during comparison: %v", err),
				)
			}
		}
	}
}

// Compare executes the specified comparison to the target value for an attribute.
//
// Type comparisons
func Compare(op CompareOp, target interface{}, meta ...interface{}) Generic {
	v := NewGenericValidator(GenericConfig{
		Description:         fmt.Sprintf("Asserts an attribute is %q to %T(%[2]v)", op, target),
		MarkdownDescription: fmt.Sprintf("Asserts an attribute is %q to %T(%[2]v)", op, target),
		TestFunc:            CompareTest(op, target, meta...),
		SkipWhenNull:        true,
		SkipWhenUnknown:     true,
	})
	return v
}

// TestIsURL asserts that the provided value can be parsed by url.Parse()
func TestIsURL(requireScheme string, requirePort int) TestFunc {
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		requireScheme := requireScheme
		requirePort := strconv.Itoa(requirePort)

		validateURL := func(v string) {
			if purl, err := url.Parse(v); err != nil {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Value is not parseable as URL",
					fmt.Sprintf("Value is not parseable as url.URL: %v", err),
				)
			} else {
				if requireScheme != "" && purl.Scheme != requireScheme {
					resp.Diagnostics.AddAttributeError(
						req.Path,
						"URL scheme mismatch",
						fmt.Sprintf("Defined scheme %q does not match required %q", purl.Scheme, requireScheme),
					)
				}
				if requirePort != "" && purl.Port() != requirePort {
					resp.Diagnostics.AddAttributeError(
						req.Path,
						"URL port mismatch",
						fmt.Sprintf("Defined port %q does not match required %q", purl.Port(), requirePort),
					)
				}
			}
		}

		if lv, ok := req.ConfigValue.(types.List); ok {
			for _, v := range lv.Elements() {
				validateURL(conv.AttributeValueToString(v))
			}
		} else if sv, ok := req.ConfigValue.(types.Set); ok {
			for _, v := range sv.Elements() {
				validateURL(conv.AttributeValueToString(v))
			}
		} else if mv, ok := req.ConfigValue.(types.Map); ok {
			for _, v := range mv.Elements() {
				validateURL(conv.AttributeValueToString(v))
			}
		} else {
			validateURL(conv.AttributeValueToString(req.ConfigValue))
		}
	}
}

// IsURLWith returns a validator that asserts a given attribute's value(s) are parseable as an URL, and that it / they
// have a specific scheme and / or port
func IsURLWith(requiredScheme string, requiredPort int) Generic {
	return NewGenericValidator(GenericConfig{
		Description:         "Tests if provided value is parseable as url.URL",
		MarkdownDescription: "Tests if provided value is parseable as url.URL",
		TestFunc:            TestIsURL(requiredScheme, requiredPort),
		SkipWhenNull:        true,
		SkipWhenUnknown:     true,
	})
}

// IsURL returns a validator that asserts a given attribute's value(s) are parseable as an URL
func IsURL() Generic {
	return IsURLWith("", 0)
}

// TestIsDurationString asserts that a given attribute's value is a valid time.Duration string
func TestIsDurationString() TestFunc {
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		if _, err := time.ParseDuration(conv.AttributeValueToString(req.ConfigValue)); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Value is not parseable as time.Duration",
				fmt.Sprintf("Value is not parseable as time.Duration: %v", err),
			)
		}
	}
}

var isDurationStringValidator = NewGenericValidator(GenericConfig{
	Description:         "Tests if value is a valid time.Duration string",
	MarkdownDescription: "Tests if value is a valid time.Duration string",
	TestFunc:            TestIsDurationString(),
	SkipWhenNull:        true,
	SkipWhenUnknown:     true,
})

// IsDurationString returns a validator that asserts a given attribute's value is a valid time.Duration string
func IsDurationString() Generic {
	return isDurationStringValidator
}

// TestEnvVarValued asserts that a given attribute value is the name of an environment variable that is valued
// at runtime
func TestEnvVarValued() TestFunc {
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		varName := conv.AttributeValueToString(req.ConfigValue)
		if v, ok := os.LookupEnv(varName); !ok || strings.TrimSpace(v) == "" {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Environment variable is either undefined or empty",
				fmt.Sprintf("The provided environment variable %q is not defined or empty", varName),
			)
		}
	}
}

var envVarValuedValidator = NewGenericValidator(GenericConfig{
	Description:         "Tests if value is an environment variable name that itself is valued",
	MarkdownDescription: "Tests if value is an environment variable name that itself is valued",
	TestFunc:            TestEnvVarValued(),
	SkipWhenNull:        true,
	SkipWhenUnknown:     true,
})

// EnvVarValued returns a validator that asserts a given attribute's value is the name of an environment variable that
// is valued at runtime
func EnvVarValued() Generic {
	return envVarValuedValidator
}

// TestFileIsReadable attempts to open and subsequently read a single byte from the file at the path specified by the
// attribute value.
func TestFileIsReadable() TestFunc {
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		fname := conv.AttributeValueToString(req.ConfigValue)
		fh, err := os.Open(fname)

		if fh != nil {
			// always try to close handle
			defer func() {
				_ = fh.Close()
			}()
		}

		// if we weren't able to open the file
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"File could not be opened for reading",
				fmt.Sprintf("File %q could not be opened for reading: %v", fname, err),
			)
			return
		}

		// attempt to read 1 byte from the file
		if _, err = io.LimitReader(fh, 1).Read(make([]byte, 1, 1)); err != nil && !errors.Is(err, io.EOF) {

			resp.Diagnostics.AddAttributeError(
				req.Path,
				"File is not readable",
				fmt.Sprintf("File %q could not be read from: %v", fname, err),
			)
		}
	}
}

var fileIsReadableValidator = NewGenericValidator(GenericConfig{
	Description:         "Tests if value is a file that exists and is readable",
	MarkdownDescription: "Tests if value is a file that exists and is readable",
	TestFunc:            TestFileIsReadable(),
	SkipWhenNull:        true,
	SkipWhenUnknown:     true,
})

// FileIsReadable returns a validator that asserts a given attribute's value is a local file that is readable.
func FileIsReadable() Generic {
	return fileIsReadableValidator
}

// MutuallyExclusiveSiblingTest ensures that a given attribute is not valued when another one is, and vice versa.
func MutuallyExclusiveSiblingTest(siblingAttr string) TestFunc {
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		// if the attribute is not valued, for whatever reason, move on.
		if conv.TestAttributeValueState(req.ConfigValue) != nil {
			return
		}

		siblingAttrValue := types.String{}
		siblingAttrPath := req.Path.ParentPath().AtName(siblingAttr)

		// try to fetch value of sibling attribute
		diags := req.Config.GetAttribute(ctx, siblingAttrPath, &siblingAttrValue)
		if diags.HasError() {
			return
		}
		// if the sibling attribute is not valued, move on
		if conv.TestAttributeValueState(siblingAttrValue) != nil {
			return
		}

		// yell about things
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Mutually exclusive value error",
			fmt.Sprintf(
				"Cannot provide value to both %q and %q",
				conv.FormatPathPathSteps(req.Path.Steps()...),
				conv.FormatPathPathSteps(siblingAttrPath.Steps()...),
			),
		)
	}
}

// MutuallyExclusiveSibling returns a validator that ensures that an attribute only carries a value when another
// sibling attribute does not.
//
// Sibling is defined as another attribute that is at the same step depth as the source attribute
func MutuallyExclusiveSibling(siblingAttr string) Generic {
	v := NewGenericValidator(GenericConfig{
		Description:         fmt.Sprintf("Ensures attribute is only valued if sibling attribute %q is empty", siblingAttr),
		MarkdownDescription: fmt.Sprintf("Ensures attribute is only valued if sibling attribute %q is empty", siblingAttr),
		TestFunc:            MutuallyExclusiveSiblingTest(siblingAttr),
		SkipWhenNull:        true,
		SkipWhenUnknown:     true,
	})

	return v
}

// MutuallyInclusiveSiblingTest ensures that a given attribute is valued when a sibling attribute is also valued
func MutuallyInclusiveSiblingTest(siblingAttr string) TestFunc {
	return func(ctx context.Context, req GenericRequest, resp *GenericResponse) {
		// if this attribute is valued, move on
		if conv.TestAttributeValueState(req.ConfigValue) == nil {
			return
		}

		siblingAttrValue := types.String{}
		siblingAttrPath := req.Path.ParentPath().AtName(siblingAttr)

		// try to fetch value of sibling attribute
		diags := req.Config.GetAttribute(ctx, siblingAttrPath, &siblingAttrValue)
		if diags.HasError() {
			return
		}
		// if the sibling attribute is not valued, move on
		if conv.TestAttributeValueState(siblingAttrValue) != nil {
			return
		}

		// yell about things
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Mutually inclusive value error",
			fmt.Sprintf(
				"Attribute %q is required when sibling attribute %q is valued",
				conv.FormatPathPathSteps(req.Path.Steps()...),
				conv.FormatPathPathSteps(siblingAttrPath.Steps()...),
			),
		)
	}
}

func MutuallyInclusiveSibling(siblingAttr string) Generic {
	v := NewGenericValidator(GenericConfig{
		Description:         fmt.Sprintf("Ensure attribute is valued when sibling attribute %q is also valued", siblingAttr),
		MarkdownDescription: fmt.Sprintf("Ensure attribute is valued when sibling attribute %q is also valued", siblingAttr),
		TestFunc:            MutuallyInclusiveSiblingTest(siblingAttr),
		SkipWhenNull:        false,
		SkipWhenUnknown:     false,
	})

	return v
}
