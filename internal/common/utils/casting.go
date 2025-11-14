package utils

import (
	"strconv"
)

func StringToUint(s string) (uint, error) {
	// Implementation of string to uint conversion
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
