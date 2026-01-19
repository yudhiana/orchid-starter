package common

import (
	"errors"
	"strings"
)

func GetChainError(err error) (msg string) {
	chainString := " -> "
	for err != nil {
		msg += err.Error() + chainString
		err = errors.Unwrap(err)
	}
	return strings.TrimSuffix(msg, chainString)
}
