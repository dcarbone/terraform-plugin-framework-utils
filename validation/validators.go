package validation

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dcarbone/terraform-plugin-framework-utils/conv"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AttributeValidator interface {
	tfsdk.AttributeValidator
	WithDescription(string) AttributeValidator
	WithMarkdownDescription(string) AttributeValidator
}

type TestFunc func(context.Context, tfsdk.ValidateAttributeRequest, *tfsdk.ValidateAttributeResponse)

type baseAttributeValidator struct {
	desc        string
	mdDesc      string
	skipNull    bool
	skipUnknown bool
	fn          TestFunc
}

func NewValidator(desc, md string, fn TestFunc, skipNull, skipUnknown bool) *baseAttributeValidator {
	v := new(baseAttributeValidator)
	v.desc = desc
	v.mdDesc = md
	v.skipNull = skipNull
	v.skipUnknown = skipUnknown
	v.fn = fn
	return v
}

func (v *baseAttributeValidator) Description(context.Context) string {
	return v.desc
}

func (v *baseAttributeValidator) MarkdownDescription(context.Context) string {
	return v.mdDesc
}

// Validate determines whether a given attribute's value was "valued" in the configuration being processed.  It will
// only allow validation to continue if there is a value to perform validation on.
func (v *baseAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	err := conv.TestAttributeValueState(req.AttributeConfig)

	if v.skipUnknown && conv.IsValueIsUnknownError(err) {
		return
	}

	if v.skipNull && conv.IsValueIsNullError(err) {
		return
	}

	// otherwise, fire away!
	v.fn(ctx, req, resp)
}

func (v *baseAttributeValidator) WithDescription(desc string) AttributeValidator {
	v.desc = desc
	return v
}

func (v *baseAttributeValidator) WithMarkdownDescription(mdDesc string) AttributeValidator {
	v.mdDesc = mdDesc
	return v
}

// RequiredTest is an AttributeValidator implementation that will register an error if the attribute has no value of
// any kind
func RequiredTest() TestFunc {
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		if conv.TestAttributeValueState(req.AttributeConfig) == nil {
			return
		}

		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Attribute must be valued",
			"Attribute must have a value configured",
		)
	}
}

// Required returns a validator that asserts a field is configured with any value at all.
func Required() AttributeValidator {
	v := NewValidator(
		"Asserts the attribute is defined and non-null",
		"Asserts the attribute is defined and non-null",
		RequiredTest(),
		false,
		true,
	)
	return v
}

// RegexpMatchTest is an AttributeValidator implementation that will first attempt to convert the value of
// a field to a string, then see if that resulting string matches the provided expression.
func RegexpMatchTest(r string) TestFunc {
	re := regexp.MustCompile(r)
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		str := conv.AttributeValueToString(req.AttributeConfig)
		if !re.MatchString(str) {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Field value does not match expression",
				fmt.Sprintf("Field value %q does not match expression %q", str, r),
			)
		}
	}
}

// RegexpMatch returns a validator that asserts an attribute's value matches the provided expression
func RegexpMatch(r string) AttributeValidator {
	v := NewValidator(
		fmt.Sprintf("Asserts attribute string value matches expression %q", r),
		fmt.Sprintf("Asserts attribute string value matches expression %q", r),
		RegexpMatchTest(r),
		true,
		true,
	)
	return v
}

// RegexpNotMatchTest is an AttributeValidator implementation that will first attempt to convert the value of
// a field to a string, then see if that resulting string matches the provided expression.
func RegexpNotMatchTest(r string) TestFunc {
	re := regexp.MustCompile(r)
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		str := conv.AttributeValueToString(req.AttributeConfig)
		if re.MatchString(str) {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Field value must NOT match expression",
				fmt.Sprintf("Field value %q matches expression %q, indicating it contains invalid characters", str, r),
			)
		}
	}
}

// RegexpNotMatch returns a validator that asserts an attribute's value does not match the provided expression
func RegexpNotMatch(r string) AttributeValidator {
	v := NewValidator(
		fmt.Sprintf("Assert attribute string value does not match expression %q", r),
		fmt.Sprintf("Assert attribute string value does not match expression %q", r),
		RegexpNotMatchTest(r),
		true,
		true,
	)
	return v
}

// LengthTest is an AttributeValidator implementation that attempts to find a length value appropriate
// for the attribute value type, asserting that it is within the specified bounds.
//
// If either min or max is set to -1, then that value is unbounded and thus not verified.
func LengthTest(minL, maxL int) TestFunc {
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		// perform some basic sanity checking
		if minL == -1 && maxL == -1 {
			resp.Diagnostics.AddAttributeWarning(
				req.AttributePath,
				"Length validation is unbounded, there is nothing to verify",
				"Both minL and maxL variables were set to -1.  This has no purpose and should be rectified.",
			)
			return
		} else if minL < -1 || maxL < -1 {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Cannot use negative value for length check",
				fmt.Sprintf("The provided minL value of %q and / or the provided maxL value %q are negative."+
					"  The only valid negative value is -1 to indicate \"unbounded\".  This should be rectified.",
					minL,
					maxL,
				))
			return
		} else if maxL != -1 && minL > maxL {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Minimum length value is greater than maximum length value",
				fmt.Sprintf(
					"The provided minimum length %q is greater than the provided maximum length of %q."+
						"  This should be rectified.",
					minL,
					maxL,
				))
			return
		}

		fl := conv.AttributeValueLength(req.AttributeConfig)
		if minL > -1 && fl < minL {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Field value length is below minimum threshold",
				fmt.Sprintf("Field value length %d is less than mininum allowed of %d", fl, minL),
			)
		}
		if maxL > -1 && fl > maxL {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Field value length is above maximum threshold",
				fmt.Sprintf("Field value length %d is greater than the maximum allowed of %d", fl, maxL),
			)
		}
	}
}

// Length asserts an attribute's length is limited to the specified bounds.  The allowed types are:
//		- string
//		- set
//		- map
//		- list
func Length(minL, maxL int) AttributeValidator {
	v := NewValidator(
		fmt.Sprintf("Asserts an attribute's value contains no less than %d and no more than %d elements, with -1 meaning unbounded", minL, maxL),
		fmt.Sprintf("Asserts an attribute's value contains no less than %d and no more than %d elements, with -1 meaning unbounded", minL, maxL),
		LengthTest(minL, maxL),
		true,
		true,
	)
	return v
}

// CompareTest executes a registered comparison function against the target attribute's value
func CompareTest(op CompareOp, target interface{}, meta ...interface{}) TestFunc {
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		err := CompareAttrValues(req.AttributeConfig, op, target, meta...)
		if err != nil {
			switch true {
			case errors.Is(err, ErrComparisonFailed):
				addComparisonFailedDiagnostic(op, target, req, resp, err)

			case errors.Is(err, ErrTypeConversionFailed):
				resp.Diagnostics.AddAttributeError(
					req.AttributePath,
					"Could not convert attribute to target type for comparison",
					fmt.Sprintf("Unable to convert attribute value type %T to %T for %q copmarison: %v", req, target, op, err))

			default:
				resp.Diagnostics.AddAttributeError(
					req.AttributePath,
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
func Compare(op CompareOp, target interface{}, meta ...interface{}) AttributeValidator {
	v := NewValidator(
		fmt.Sprintf("Asserts an attribute is %q to %T(%[2]v)", op, target),
		fmt.Sprintf("Asserts an attribute is %q to %T(%[2]v)", op, target),
		CompareTest(op, target, meta...),
		true,
		true,
	)
	return v
}

// TestIsURL asserts that the provided value can be parsed by url.Parse()
func TestIsURL() TestFunc {
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		if _, err := url.Parse(conv.AttributeValueToString(req.AttributeConfig)); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Value is not parseable as URL",
				fmt.Sprintf("Value is not parseable as url.URL: %v", err),
			)
		}
	}
}

// IsURL returns a validator that asserts a given attribute value is parseable as a URL
func IsURL() AttributeValidator {
	v := NewValidator(
		"Tests if provided value is parseable as url.URL",
		"Tests if provided value is parseable as url.URL",
		TestIsURL(),
		true,
		true,
	)

	return v
}

// TestIsDurationString asserts that a given attribute's value is a valid time.Duration string
func TestIsDurationString() TestFunc {
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		if _, err := time.ParseDuration(conv.AttributeValueToString(req.AttributeConfig)); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Value is not parseable as time.Duration",
				fmt.Sprintf("Value is not parseable as time.Duration: %v", err),
			)
		}
	}
}

// IsDurationString returns a validator that asserts a given attribute's value is a valid time.Duration string
func IsDurationString() AttributeValidator {
	v := NewValidator(
		"Tests if value is a valid time.Duration string",
		"Tests if value is a valid time.Duration string",
		TestIsDurationString(),
		true,
		true,
	)

	return v
}

// TestEnvVarValued asserts that a given attribute value is the name of an environment variable that is valued
// at runtime
func TestEnvVarValued() TestFunc {
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		varName := conv.StringValueToString(req.AttributeConfig)
		if v, ok := os.LookupEnv(varName); !ok || strings.TrimSpace(v) == "" {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Environment variable is either undefined or empty",
				fmt.Sprintf("The provided environment variable %q is not defined or empty", varName),
			)
		}
	}
}

// EnvVarValued returns a validator that asserts a given attribute's value is the name of an environment variable that
// is valued at runtime
func EnvVarValued() AttributeValidator {
	v := NewValidator(
		"Tests if value is an environment variable name that itself is valued",
		"Tests if value is an environment variable name that itself is valued",
		TestEnvVarValued(),
		true,
		true,
	)

	return v
}

// TestFileIsReadable attempts to open and subsequently read a single byte from the file at the path specified by the
// attribute value.
func TestFileIsReadable() TestFunc {
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		fname := conv.StringValueToString(req.AttributeConfig)
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
				req.AttributePath,
				"File could not be opened for reading",
				fmt.Sprintf("File %q could not be opened for reading: %v", fname, err),
			)
			return
		}

		// attempt to read 1 byte from the file
		if _, err = io.LimitReader(fh, 1).Read(make([]byte, 1, 1)); err != nil && !errors.Is(err, io.EOF) {

			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"File is not readable",
				fmt.Sprintf("File %q could not be read from: %v", fname, err),
			)
		}
	}
}

// FileIsReadable returns a validator that asserts a given attribute's value is a local file that is readable.
func FileIsReadable() AttributeValidator {
	v := NewValidator(
		"Tests if value is a file that exists and is readable",
		"Tests if value is a file that exists and is readable",
		TestFileIsReadable(),
		true,
		true,
	)

	return v
}

// MutuallyExclusiveSiblingTest ensures that a given attribute is not valued when another one is, and vice versa.
func MutuallyExclusiveSiblingTest(siblingAttr string) TestFunc {
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		// if the attribute is not valued, for whatever reason, move on.
		if conv.TestAttributeValueState(req.AttributeConfig) != nil {
			return
		}

		siblingAttrValue := types.String{}
		siblingAttrPath := req.AttributePath.WithoutLastStep().WithAttributeName(siblingAttr)

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
			req.AttributePath,
			"Mutually exclusive value error",
			fmt.Sprintf(
				"Cannot provide value to both %q and %q",
				conv.FormatAttributePathSteps(req.AttributePath.Steps()...),
				conv.FormatAttributePathSteps(siblingAttrPath.Steps()...),
			),
		)
	}
}

// MutuallyExclusiveSibling returns a validator that ensures that an attribute only carries a value when another
// sibling attribute does not.
//
// Sibling is defined as another attribute that is at the same step depth as the source attribute
func MutuallyExclusiveSibling(siblingAttr string) AttributeValidator {
	v := NewValidator(
		fmt.Sprintf("Ensures attribute is only valued if sibling attribute %q is empty", siblingAttr),
		fmt.Sprintf("Ensures attribute is only valued if sibling attribute %q is empty", siblingAttr),
		MutuallyExclusiveSiblingTest(siblingAttr),
		true,
		true,
	)

	return v
}

// MutuallyInclusiveSiblingTest ensures that a given attribute is valued when a sibling attribute is also valued
func MutuallyInclusiveSiblingTest(siblingAttr string) TestFunc {
	return func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
		// if this attribute is valued, move on
		if conv.TestAttributeValueState(req.AttributeConfig) == nil {
			return
		}

		siblingAttrValue := types.String{}
		siblingAttrPath := req.AttributePath.WithoutLastStep().WithAttributeName(siblingAttr)

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
			req.AttributePath,
			"Mutually inclusive value error",
			fmt.Sprintf(
				"Attribute %q is required when sibling attribute %q is valued",
				conv.FormatAttributePathSteps(req.AttributePath.Steps()...),
				conv.FormatAttributePathSteps(siblingAttrPath.Steps()...),
			),
		)
	}
}

func MutuallyInclusiveSibling(siblingAttr string) AttributeValidator {
	v := NewValidator(
		fmt.Sprintf("Ensure attribute is valued when sibling attribute %q is also valued", siblingAttr),
		fmt.Sprintf("Ensure attribute is valued when sibling attribute %q is also valued", siblingAttr),
		MutuallyInclusiveSiblingTest(siblingAttr),
		false,
		false,
	)

	return v
}
