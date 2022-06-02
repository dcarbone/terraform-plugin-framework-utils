package util

import (
	"fmt"
	"math/big"
	"strconv"
)

func TryCoerceToBool(in interface{}) (bool, error) {
	switch in.(type) {
	case bool:
		return in.(bool), nil
	case string:
		return strconv.ParseBool(in.(string))

	default:
		return false, fmt.Errorf("unandled type to bool conversion: %T", in)
	}
}

func TryCoerceToInt(in interface{}) (int, error) {
	switch in.(type) {
	case int:
		return in.(int), nil
	case int8:
		return int(in.(int8)), nil
	case int16:
		return int(in.(int16)), nil
	case int32:
		return int(in.(int32)), nil
	case int64:
		return int(in.(int64)), nil

	case float32:
		return int(in.(float32)), nil
	case float64:
		return int(in.(float64)), nil

	case string:
		return strconv.Atoi(in.(string))

	case *big.Float:
		bf := in.(*big.Float)
		if bf == nil {
			return 0, nil
		}
		out, _ := bf.Int64()
		return int(out), nil

	default:
		return 0, fmt.Errorf("unandled type to int conversion: %T", in)
	}
}

func TryCoerceToInt64(in interface{}) (int64, error) {
	switch in.(type) {
	case int:
		return int64(in.(int)), nil
	case int8:
		return int64(in.(int8)), nil
	case int16:
		return int64(in.(int16)), nil
	case int32:
		return int64(in.(int32)), nil
	case int64:
		return in.(int64), nil

	case float32:
		return int64(in.(float32)), nil
	case float64:
		return int64(in.(float64)), nil

	case string:
		return strconv.ParseInt(in.(string), 10, 64)

	case *big.Float:
		bf := in.(*big.Float)
		if bf == nil {
			return 0, nil
		}
		out, _ := bf.Int64()
		return out, nil

	default:
		return 0, fmt.Errorf("unandled type to int64 conversion: %T", in)
	}
}

func TryCoerceToFloat64(in interface{}) (float64, error) {
	switch in.(type) {
	case int:
		return float64(in.(int)), nil
	case int8:
		return float64(in.(int8)), nil
	case int16:
		return float64(in.(int16)), nil
	case int32:
		return float64(in.(int32)), nil
	case int64:
		return float64(in.(int64)), nil

	case float32:
		return float64(in.(float32)), nil
	case float64:
		return in.(float64), nil

	case string:
		return strconv.ParseFloat(in.(string), 64)

	case *big.Float:
		bf := in.(*big.Float)
		if bf == nil {
			return 0, nil
		}
		out, _ := bf.Float64()
		return out, nil

	default:
		return 0, fmt.Errorf("unandled type to float64 conversion: %T", in)
	}
}

func TryCoerceToBigFloat(in interface{}) (*big.Float, error) {
	f64, err := TryCoerceToFloat64(in)
	if err != nil {
		return nil, err
	}
	return big.NewFloat(f64), nil
}

func GetPrintableTypeWithValue(in interface{}) string {
	switch in.(type) {
	case string:
		return fmt.Sprintf("%[1]T(%[1]q)", in)

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%[1]T(%[1]d)", in)

	case float32, float64:
		return fmt.Sprintf("%[1]T(%[1]f)", in)

	case bool:
		return fmt.Sprintf("%[1]T(%[1]f)", in)

	default:
		return fmt.Sprintf(fmt.Sprintf("%[1]T(%[1]v)", in))
	}
}
