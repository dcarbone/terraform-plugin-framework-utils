package conv

import (
	"fmt"
	"math/big"
	"strconv"
)

// GoNumberToString is a laziness helper to convert several types to a usable string value.
func GoNumberToString(num interface{}) string {
	switch num.(type) {
	case int:
		return strconv.Itoa(num.(int))
	case int32:
		return strconv.Itoa(int(num.(int32)))
	case int64:
		return strconv.Itoa(int(num.(int64)))
	case float64:
		return strconv.FormatFloat(num.(float64), 'g', int(FloatPrecision), 64)
	case *big.Float:
		return num.(*big.Float).String()

	default:
		return fmt.Sprintf("%T", num)
	}
}
