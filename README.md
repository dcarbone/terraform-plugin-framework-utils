# terraform-plugin-framework-utils
Utilities for use with the
[HashiCorp Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)

This project, much like the framework itself, is a work in progress.  I will try to keep it as up to date with upstream
changes as possible but, as always, community help is appreciated!

# Index

* [Type Conversion](#type-conversion)
* [Attribute Validation](#attribute-validation)
* [Test Utilities](#test-utilities)

# Type Conversion

Converting between types used internally by Terraform and typical Go types can be somewhat tricky and / or tedious.

To help with this, I have created a small suite of type conversion utilities designed to make converting to
and from Terraform and Go easy and obvious.

You can see the complete list of available conversions here: 
[terraform-plugin-framework-utils/conv](https://github.com/dcarbone/terraform-plugin-framework-utils/blob/main/conv)

# Attribute Validation

The Terraform Plugin Framework has a great [validation](https://www.terraform.io/plugin/framework/validation) interface
that makes it easy to create custom validation rules for the various schemas you define in your provider.

Being in beta, the HashiCorp team has yet to provide a built-in set of providers.  I have created a few that I have
found useful when creating my own providers, and provided a small wrapper to make creating new providers as
simple as [defining a function](validation/validators.go#25).

## Provided Validators

### Required

Fails validation if the attribute is null or unknown

```go
{
	Validators: []tfsdk.AttributeValidator{
        validation.Required()
    },
}
```

### RegexpMatch

Fails validation if the attribute's value that does not match the user-defined regular expression.  This validator
will attempt to convert the attribute to a string first. 

```go
{
    Validators: []tfsdk.AttributeValidator{
        validation.RegexpMatch("{{ your regex here }}")
    },
}
```

### RegexpNotMatch

Fails validation if the attribute's value matches the user-defined regular expression.  This validator
will attempt to convert the attribute to a string first.

```go
v := validation.RegexpNotMatch("{{ your regex here }}")
```

### Length

Fails validation if the attribute's value's length is not within the specified bounds.

```go
{
    Validators: []tfsdk.AttributeValidator{
        // lower limit
        validation.Lenght(5, -1),

        // upper limit
        validation.Length(-1, 10),

        // lower and upper limit
        validation.length(5, 10),
    },
}
```

### Compare

Fails validation if the attribute's value does not match the configured comparison operation.

See [comparison.go](validation/comparison.go) for details on what comparison operations are available.  You can add
your own [ComparisonFunc](validation/comparison.go#44) using [SetComparisonFunc](validation/comparison.go#229)

```go
{
    Validators: []tfsdk.AttributeValidator{
        // equal
        validation.Compare(validation.Equal, 5),

        // less than
        validation.Compare(validation.LessThan, 10),

        // less than or equal to
        validation.Compare(validation.LessThanOrEqualTo, 10),

        // greater than
        validation.Compare(validation.GreaterThan, 5),

        // greater than or equal to
        validation.Compare(validation.GreaterThanOrEqualTo, 5),

        // not equal
        validation.Compare(validation.NotEqual, 10).
    }
}
```

### IsURL

Fails validation if the attribute's value is not parseable by `url.Parse`

```go
{
    Validators: []tfsdk.AttributeValidator{
        validation.IsUrl()
    }
}
```

### IsDurationString

Fails validation if the attribute's value is not parseable by `time.ParseDuration`

```go
{
    Validators: []tfsdk.AttributeValidator{
        validation.IsDurationString()
    }
}
```

### EnvVarValued

Fails validation if the environment variable name defined by the attribute's value is, itself, not valued at runtime.

```go
{
    Validators: []tfsdk.AttributeValidator{
        validation.EnvVarValued()
    }
}
```

### FileIsReadable

Fails validation if the file at the path defined in the attribute's value is not readable at runtime.

```go
{
    Validators: []tfsdk.AttributeValidator{
        validation.FileIsReadable()
    }
}
```

### MutuallyExclusiveSibling

Fails validation if the attribute is valued and the configured sibling attribute is also valued.

```go
{
    Validators: []tfsdk.AttributeValidator{
        validation.MutuallyExclusiveSibling("{{ sibling field name }}")
    }
}
```

#### Example

```hcl
# Example provider Terraform HCL
provider "whatever" {
  address = "http://example.com"
  address_env = "EXAMPLE_ADDR"
}
```

```go
// Example validators list defined on the `address` attribute's schema
{
    Validators: []tfsdk.AttributeValidator{
        validation.MutuallyExclusiveSibling("address_env")
    }
}
```

Adding the above validator to the `address` attribute's `Validators` list above will require that the `address_env`
field must be empty when `address` is defined.  You may also add same validator to the `address_env` attribute, this
time pointing at the `address` field.

### MutuallyInclusiveSibling

Requires that two sibling attributes either both be valued or not valued.

```go
{
    Validators: []tfsdk.AttributeValidator{
        validation.MutuallyInclusiveSibling("{{ sibling field name }}")
    }
}
```

#### Example

```hcl
# Example provider Terraform HCL
provider "whatever" {
  ssh_key_file = file("local/filepath/ssh.key")
  ssh_key_password = null
}
```

```go
// Example validators list defined on the `ssh_key_password` attribute's schema
{
    Validators: []tfsdk.AttributeValidator{
        validation.MutuallyInclusiveSibling("ssh_key")
    }
}
```

Adding the above validator to the `ssh_key_password` attribute's `Validators` list will require that, if the
`ssh_key_file` attribute is defined so, too, must the `ssh_key_password` attribute be valued.

# Test Utilities

The Terraform Provider Framework provides an excellent suite of 
[test tools](https://www.terraform.io/plugin/framework/acctests) to use when creating unit and acceptance tests for
provider.

For my uses, I wanted a way to construct hcl config blocks without having to define a heredoc string for each one.

So I created a few [config utilities](acctest/config.go) to assist with this.

## Example

```go
fieldMap := map[string]interface{}{
	"address": "http://example.com",
	"token": acctest.Literal("file(\"/location/on/disk/token\")"),
	"number_of_fish_in_the_sea": 3500000000000,
}
confHCL := acctest.CompileProviderConfig("my_provider", map[string]interface{}{})
```

```hcl
provider "my_provider" {
  address = "http://example.com"
  token = file("/location/on/disk/token")
  number_of_fish_in_the_sea = 3500000000000
}
```

This can be used with the `acctest.JoinConfigs` func to bring together multiple reusable configuration blocks for 
different tests.