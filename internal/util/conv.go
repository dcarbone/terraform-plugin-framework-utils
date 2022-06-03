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

func TryCoerceToInts(in interface{}) ([]int, error) {
	switch in.(type) {
	case []int:
		out := make([]int, len(in.([]int)))
		copy(out, in.([]int))
		return out, nil
	case []int8:
		out := make([]int, len(in.([]int8)))
		for i, v := range in.([]int8) {
			out[i] = int(v)
		}
		return out, nil
	case []int16:
		out := make([]int, len(in.([]int16)))
		for i, v := range in.([]int16) {
			out[i] = int(v)
		}
		return out, nil
	case []int32:
		out := make([]int, len(in.([]int32)))
		for i, v := range in.([]int32) {
			out[i] = int(v)
		}
		return out, nil
	case []int64:
		out := make([]int, len(in.([]int64)))
		for i, v := range in.([]int64) {
			out[i] = int(v)
		}
		return out, nil

	case []uint:
		out := make([]int, len(in.([]uint)))
		for i, v := range in.([]uint) {
			out[i] = int(v)
		}
		return out, nil
	case []uint8:
		out := make([]int, len(in.([]uint8)))
		for i, v := range in.([]uint8) {
			out[i] = int(v)
		}
		return out, nil
	case []uint16:
		out := make([]int, len(in.([]uint16)))
		for i, v := range in.([]uint16) {
			out[i] = int(v)
		}
		return out, nil
	case []uint32:
		out := make([]int, len(in.([]int8)))
		for i, v := range in.([]int8) {
			out[i] = int(v)
		}
		return out, nil
	case []uint64:
		out := make([]int, len(in.([]uint64)))
		for i, v := range in.([]uint64) {
			out[i] = int(v)
		}
		return out, nil

	case []float32:
		out := make([]int, len(in.([]float32)))
		for i, v := range in.([]float32) {
			out[i] = int(v)
		}
		return out, nil
	case []float64:
		out := make([]int, len(in.([]float64)))
		for i, v := range in.([]float64) {
			out[i] = int(v)
		}
		return out, nil

	case []string:
		out := make([]int, len(in.([]string)))
		for i, v := range in.([]string) {
			if p, err := strconv.Atoi(v); err != nil {
				return nil, fmt.Errorf("offset %d(%q) cannot be parsed as int: %w", i, v, err)
			} else {
				out[i] = p
			}
		}
		return out, nil

	case []*big.Float:
		out := make([]int, 0)
		for _, v := range in.([]*big.Float) {
			if v != nil {
				vv, _ := v.Int64()
				out = append(out, int(vv))
			}
		}
		return out, nil
	case []big.Float:
		out := make([]int, len(in.([]big.Float)))
		for i, v := range in.([]big.Float) {
			vv, _ := v.Int64()
			out[i] = int(vv)
		}
		return out, nil

	default:
		return nil, fmt.Errorf("unandled type to []int conversion: %T", in)
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

func TryCoerceToInt64s(in interface{}) ([]int64, error) {
	switch in.(type) {
	case []int:
		out := make([]int64, len(in.([]int)))
		for i, v := range in.([]int) {
			out[i] = int64(v)
		}
		return out, nil
	case []int8:
		out := make([]int64, len(in.([]int8)))
		for i, v := range in.([]int8) {
			out[i] = int64(v)
		}
		return out, nil
	case []int16:
		out := make([]int64, len(in.([]int16)))
		for i, v := range in.([]int16) {
			out[i] = int64(v)
		}
		return out, nil
	case []int32:
		out := make([]int64, len(in.([]int32)))
		for i, v := range in.([]int32) {
			out[i] = int64(v)
		}
		return out, nil
	case []int64:
		out := make([]int64, len(in.([]int64)))
		copy(out, in.([]int64))
		return out, nil

	case []uint:
		out := make([]int64, len(in.([]uint)))
		for i, v := range in.([]uint) {
			out[i] = int64(v)
		}
		return out, nil
	case []uint8:
		out := make([]int64, len(in.([]uint8)))
		for i, v := range in.([]uint8) {
			out[i] = int64(v)
		}
		return out, nil
	case []uint16:
		out := make([]int64, len(in.([]uint16)))
		for i, v := range in.([]uint16) {
			out[i] = int64(v)
		}
		return out, nil
	case []uint32:
		out := make([]int64, len(in.([]int8)))
		for i, v := range in.([]int8) {
			out[i] = int64(v)
		}
		return out, nil
	case []uint64:
		out := make([]int64, len(in.([]uint64)))
		for i, v := range in.([]uint64) {
			out[i] = int64(v)
		}
		return out, nil

	case []float32:
		out := make([]int64, len(in.([]float32)))
		for i, v := range in.([]float32) {
			out[i] = int64(v)
		}
		return out, nil
	case []float64:
		out := make([]int64, len(in.([]float64)))
		for i, v := range in.([]float64) {
			out[i] = int64(v)
		}
		return out, nil

	case []string:
		out := make([]int64, len(in.([]string)))
		for i, v := range in.([]string) {
			if p, err := strconv.ParseInt(v, 10, 64); err != nil {
				return nil, fmt.Errorf("offset %d(%q) cannot be parsed as int64: %w", i, v, err)
			} else {
				out[i] = p
			}
		}
		return out, nil

	case []*big.Float:
		out := make([]int64, 0)
		for _, v := range in.([]*big.Float) {
			if v != nil {
				vv, _ := v.Int64()
				out = append(out, vv)
			}
		}
		return out, nil
	case []big.Float:
		out := make([]int64, len(in.([]big.Float)))
		for i, v := range in.([]big.Float) {
			out[i], _ = v.Int64()
		}
		return out, nil

	default:
		return nil, fmt.Errorf("unandled type to []int conversion: %T", in)
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

func TryCoerceToFloats(in interface{}) ([]float64, error) {
	switch in.(type) {
	case []int:
		out := make([]float64, len(in.([]int)))
		for i, v := range in.([]int) {
			out[i] = float64(v)
		}
		return out, nil
	case []int8:
		out := make([]float64, len(in.([]int8)))
		for i, v := range in.([]int8) {
			out[i] = float64(v)
		}
		return out, nil
	case []int16:
		out := make([]float64, len(in.([]int16)))
		for i, v := range in.([]int16) {
			out[i] = float64(v)
		}
		return out, nil
	case []int32:
		out := make([]float64, len(in.([]int32)))
		for i, v := range in.([]int32) {
			out[i] = float64(v)
		}
		return out, nil
	case []int64:
		out := make([]float64, len(in.([]int64)))
		for i, v := range in.([]int64) {
			out[i] = float64(v)
		}
		return out, nil

	case []uint:
		out := make([]float64, len(in.([]uint)))
		for i, v := range in.([]uint) {
			out[i] = float64(v)
		}
		return out, nil
	case []uint8:
		out := make([]float64, len(in.([]uint8)))
		for i, v := range in.([]uint8) {
			out[i] = float64(v)
		}
		return out, nil
	case []uint16:
		out := make([]float64, len(in.([]uint16)))
		for i, v := range in.([]uint16) {
			out[i] = float64(v)
		}
		return out, nil
	case []uint32:
		out := make([]float64, len(in.([]int8)))
		for i, v := range in.([]int8) {
			out[i] = float64(v)
		}
		return out, nil
	case []uint64:
		out := make([]float64, len(in.([]uint64)))
		for i, v := range in.([]uint64) {
			out[i] = float64(v)
		}
		return out, nil

	case []string:
		out := make([]float64, len(in.([]string)))
		for i, v := range in.([]string) {
			if p, err := strconv.ParseFloat(v, 64); err != nil {
				return nil, fmt.Errorf("offset %d(%q) cannot be parsed as float64: %w", i, v, err)
			} else {
				out[i] = p
			}
		}
		return out, nil

	case []float32:
		out := make([]float64, len(in.([]float32)))
		for i, v := range in.([]float32) {
			out[i] = float64(v)
		}
		return out, nil
	case []float64:
		out := make([]float64, len(in.([]float64)))
		copy(out, in.([]float64))
		return out, nil

	case []*big.Float:
		out := make([]float64, 0)
		for _, v := range in.([]*big.Float) {
			if v != nil {
				vv, _ := v.Float64()
				out = append(out, vv)
			}
		}
		return out, nil
	case []big.Float:
		out := make([]float64, len(in.([]big.Float)))
		for i, v := range in.([]big.Float) {
			vv, _ := v.Float64()
			out[i] = vv
		}
		return out, nil

	default:
		return nil, fmt.Errorf("unandled type to []int conversion: %T", in)
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
