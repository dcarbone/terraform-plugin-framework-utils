package acctest

import (
	"fmt"
	"strings"
)

type ConfigLiteral string

// ConfigValue attempts to convert the provided input to a Terraform config safe representation of its value
func ConfigValue(in interface{}) string {
	switch in.(type) {
	case nil:
		return "null"

	case ConfigLiteral:
		return string(in.(ConfigLiteral))

	case string:
		if strings.Contains(in.(string), "\n") {
			return fmt.Sprintf("<<EOD\n%s\nEOD", in.(string))
		} else {
			return fmt.Sprintf("%q", in.(string))
		}

	case bool:
		return fmt.Sprintf("%t", in.(bool))

	case int:
		return fmt.Sprintf("%d", in.(int))

	case float64:
		return fmt.Sprintf("%f", in.(float64))

	case []interface{}:
		formatted := make([]string, 0)
		for _, v := range in.([]interface{}) {
			formatted = append(formatted, ConfigValue(v))
		}
		return fmt.Sprintf("[\n%s]", strings.Join(formatted, ",\n"))
	case []string:
		formatted := make([]string, 0)
		for _, v := range in.([]string) {
			formatted = append(formatted, ConfigValue(v))
		}
		return fmt.Sprintf("[\n%s]", strings.Join(formatted, ",\n"))
	case []int:
		formatted := make([]string, 0)
		for _, v := range in.([]int) {
			formatted = append(formatted, ConfigValue(v))
		}
		return fmt.Sprintf("[\n%s]", strings.Join(formatted, ",\n"))
	case []float64:
		formatted := make([]string, 0)
		for _, v := range in.([]float64) {
			formatted = append(formatted, ConfigValue(v))
		}
		return fmt.Sprintf("[\n%s]", strings.Join(formatted, ",\n"))

	case map[string]interface{}:
		inner := "{"
		for k, v := range in.(map[string]interface{}) {
			inner = fmt.Sprintf("%s\n%s = %s", inner, k, ConfigValue(v))
		}
		return fmt.Sprintf("%s\n}", inner)
	case map[string]string:
		inner := "{"
		for k, v := range in.(map[string]string) {
			inner = fmt.Sprintf("%s\n%s = %s", inner, k, ConfigValue(v))
		}
		return fmt.Sprintf("%s\n}", inner)

	default:
		panic(fmt.Sprintf("Unable to handle config values of type %T", in))
	}
}

func JoinConfigs(confs ...string) string {
	return strings.Join(confs, "\n")
}

func CompileConfig(header string, fieldMaps ...map[string]interface{}) string {
	const f = `
%s {
%s
}`

	fields := ""
	for k, v := range MergeMaps(fieldMaps...) {
		fields = fmt.Sprintf("%s%s = %s\n", fields, k, ConfigValue(v))
	}

	return fmt.Sprintf(f, header, fields)
}

func ProviderHeader(name string) string {
	const f = `provider %q`
	return fmt.Sprintf(f, name)
}

func ResourceHeader(resourceType, resourceName string) string {
	const f = `resource %q %q`
	return fmt.Sprintf(f, resourceType, resourceName)
}

func DataSourceHeader(dataSourceType, dataSourceName string) string {
	const f = `data %q %q`
	return fmt.Sprintf(f, dataSourceType, dataSourceName)
}

func CompileProviderConfig(providerName string, fieldMaps ...map[string]interface{}) string {
	return CompileConfig(
		ProviderHeader(providerName),
		fieldMaps...,
	)
}

func CompileResourceConfig(resourceType, resourceName string, fieldMaps ...map[string]interface{}) string {
	return CompileConfig(
		ResourceHeader(resourceType, resourceName),
		fieldMaps...,
	)
}

func CompileDataSourceConfig(dataSourceType, dataSourceName string, fieldMaps ...map[string]interface{}) string {
	return CompileConfig(
		DataSourceHeader(dataSourceType, dataSourceName),
		fieldMaps...,
	)
}

func CompileLocalsConfig(fieldMaps ...map[string]interface{}) string {
	return CompileConfig("locals", fieldMaps...)
}
