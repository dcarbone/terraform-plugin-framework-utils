package conv

import "errors"

var (
	ErrValueIsNull    = errors.New("value is nil")
	ErrValueIsUnknown = errors.New("value is unknown")
	ErrValueIsEmpty   = errors.New("value is empty")
)

func IsValueIsNullError(err error) bool {
	for err != nil {
		if errors.Is(err, ErrValueIsNull) {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

func IsValueIsUnknownError(err error) bool {
	for err != nil {
		if errors.Is(err, ErrValueIsUnknown) {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

func IsValueIsEmptyError(err error) bool {
	for err != nil {
		if errors.Is(err, ErrValueIsEmpty) {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}
