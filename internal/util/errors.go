package util

import "errors"

func MatchError(src, tgt error) bool {
	for src != nil {
		if errors.Is(src, tgt) {
			return true
		}
		src = errors.Unwrap(src)
	}
	return false
}
