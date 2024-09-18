package utils

import (
	"strconv"
)

// Конвертация string в uint
func ToUint(str string) (uint64, error) {
	id, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}

	return id, nil
}
