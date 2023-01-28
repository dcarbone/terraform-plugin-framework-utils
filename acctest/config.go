package acctest

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/internal/util"
)

type ConfigLiteral string

type ConfigValueFunc func(interface{}) string

var (
	configValueFuncsMu sync.RWMutex
	configValueFuncs   map[string]ConfigValueFunc
)

func SetConfigValueFunc(t interface{}, fn ConfigValueFunc) {
	configValueFuncsMu.Lock()
	defer configValueFuncsMu.Unlock()
	configValueFuncs[util.KeyFN(t)] = fn
}

func GetConfigValueFunc(t interface{}) (ConfigValueFunc, bool) {
	configValueFuncsMu.RLock()
	defer configValueFuncsMu.RUnlock()
	fn, ok := configValueFuncs[util.KeyFN(t)]
	return fn, ok
}

func DefaultConfigValueFuncs() map[string]ConfigValueFunc {
	return map[string]ConfigValueFunc{
		// nil
		util.KeyFN(nil): func(_ interface{}) string { return "null" },

		// values to print literally
		util.KeyFN(ConfigLiteral("")): func(v interface{}) string { return string(v.(ConfigLiteral)) },

		// simple conversions
		util.KeyFN(false):      func(v interface{}) string { return fmt.Sprintf("%t", v.(bool)) },
		util.KeyFN(0):          func(v interface{}) string { return fmt.Sprintf("%d", v.(int)) },
		util.KeyFN(float64(0)): func(v interface{}) string { return fmt.Sprintf("%f", v.(float64)) },

		// handle single and multi-line strings
		util.KeyFN(""): func(v interface{}) string {
			if strings.Contains(v.(string), "\n") {
				return fmt.Sprintf("<<EOD\n%s\nEOD", v.(string))
			} else {
				return fmt.Sprintf("%q", v.(string))
			}
		},

		// time values

		util.KeyFN(time.Nanosecond): func(v interface{}) string {
			return ConfigValue(v.(time.Duration).String())
		},

		// slices

		util.KeyFN(make([]interface{}, 0)): func(v interface{}) string {
			formatted := make([]string, 0)
			for _, v := range v.([]interface{}) {
				formatted = append(formatted, ConfigValue(v))
			}
			return fmt.Sprintf("[\n%s\n]", strings.Join(formatted, ",\n"))
		},
		util.KeyFN(make([]string, 0)): func(v interface{}) string {
			formatted := make([]string, 0)
			for _, v := range v.([]string) {
				formatted = append(formatted, ConfigValue(v))
			}
			return fmt.Sprintf("[\n%s\n]", strings.Join(formatted, ",\n"))
		},
		util.KeyFN(make([]int, 0)): func(v interface{}) string {
			formatted := make([]string, 0)
			for _, v := range v.([]int) {
				formatted = append(formatted, ConfigValue(v))
			}
			return fmt.Sprintf("[\n%s\n]", strings.Join(formatted, ",\n"))
		},
		util.KeyFN(make([]float64, 0)): func(v interface{}) string {
			formatted := make([]string, 0)
			for _, v := range v.([]float64) {
				formatted = append(formatted, ConfigValue(v))
			}
			return fmt.Sprintf("[\n%s\n]", strings.Join(formatted, ",\n"))
		},

		// maps

		util.KeyFN(make(map[string]interface{})): func(v interface{}) string {
			inner := "{"
			for k, v := range v.(map[string]interface{}) {
				inner = fmt.Sprintf("%s\n%s = %s", inner, k, ConfigValue(v))
			}
			return fmt.Sprintf("%s\n}", inner)
		},
		util.KeyFN(make(map[string]string)): func(v interface{}) string {
			inner := "{"
			for k, v := range v.(map[string]string) {
				inner = fmt.Sprintf("%s\n%s = %s", inner, k, ConfigValue(v))
			}
			return fmt.Sprintf("%s\n}", inner)
		},
	}
}

func init() {
	configValueFuncs = DefaultConfigValueFuncs()
}

// ConfigValue attempts to convert the provided input to a Terraform config safe representation of its value
func ConfigValue(in interface{}) string {
	if fn, ok := GetConfigValueFunc(in); ok {
		return fn(in)
	} else {
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
