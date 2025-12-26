package customscalar

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

//MarshalFloat32 --
func MarshalFloat32(f float32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%g", f))
	})
}

//UnmarshalFloat32 --
func UnmarshalFloat32(v interface{}) (float32, error) {
	switch v := v.(type) {
	case string:
		conv, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return 0, err
		}
		return float32(conv), nil
	case int:
		return float32(v), nil
	case int64:
		return float32(v), nil
	case float64:
		return float32(v), nil
	case float32:
		return v, nil
	case json.Number:
		conv, err := strconv.ParseFloat(string(v), 32)
		if err != nil {
			return 0, err
		}
		return float32(conv), nil
	default:
		return 0, fmt.Errorf("%T is not an float", v)
	}
}

//MarshalFloat64 --
func MarshalFloat64(f float64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%g", f))
	})
}

//UnmarshalFloat64 --
func UnmarshalFloat64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case json.Number:
		return strconv.ParseFloat(string(v), 64)
	default:
		return 0, fmt.Errorf("%T is not an float", v)
	}
}
