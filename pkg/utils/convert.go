package utils

import (
	"strconv"
)

type ErrConvertToInt struct {
	Reason string
}

func (e *ErrConvertToInt) Error() string {
	return e.Reason
}

// Конвертация string в uint
func ToUint(str string) (uint64, error) {
	id, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, &ErrConvertToInt{Reason: err.Error()}
	}

	return id, nil
}
