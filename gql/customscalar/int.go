package customscalar

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalInt8 --
func MarshalInt8(i int8) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(int64(i), 10))
	})
}

// UnmarshalInt8 --
func UnmarshalInt8(v interface{}) (int8, error) {
	switch v := v.(type) {
	case string:
		iv, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return 0, err
		}
		return int8(iv), nil
	case int:
		return int8(v), nil
	case int64:
		return int8(v), nil
	case json.Number:
		iv, err := strconv.ParseInt(string(v), 10, 8)
		if err != nil {
			return 0, err
		}
		return int8(iv), nil
	default:
		return 0, fmt.Errorf("%T is not an int", v)
	}
}

// MarshalInt64 --
func MarshalInt64(i int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(i, 10))
	})
}

// UnmarshalInt64 --
func UnmarshalInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case string:
		iv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return iv, nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case json.Number:
		iv, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return 0, err
		}
		return iv, nil
	default:
		return 0, fmt.Errorf("%T is not an int", v)
	}
}
