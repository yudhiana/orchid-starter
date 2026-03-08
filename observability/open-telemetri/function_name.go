package openTelemetri

import (
	"runtime"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/attribute"
)

// cache untuk menyimpan nama fungsi yang sudah diambil
var funcNameCache sync.Map

// GetFuncName mengambil nama fungsi yang memanggilnya
// akan disimpan di cache agar tidak dihitung berulang
func GetFuncName() string {
	// ambil program counter dari caller 1 level di atas
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	// cek cache dulu
	if name, ok := funcNameCache.Load(pc); ok {
		return name.(string)
	}

	// ambil nama fungsi dari PC
	fn := runtime.FuncForPC(pc)
	var name string
	if fn != nil {
		name = fn.Name()
		name = name[strings.LastIndex(name, "/")+1:]
	} else {
		name = "unknown"
	}

	// simpan ke cache
	funcNameCache.Store(pc, name)
	return name
}

func MakeTags(attrs map[string]any) []attribute.KeyValue {
	var kvs []attribute.KeyValue
	for k, v := range attrs {
		switch val := v.(type) {
		case string:
			kvs = append(kvs, attribute.String(k, val))
		case int:
			kvs = append(kvs, attribute.Int(k, val))
		case bool:
			kvs = append(kvs, attribute.Bool(k, val))
		}
	}
	return kvs
}
