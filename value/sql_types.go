package value

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type Value int

const (
	Invalid Value = iota
	Time
	Timestamp
	NullBool
	NullFloat64
	NullInt32
	NullInt64
	NullString
	NullTime
	Float64
	Float32
	Uint
	Int
	Int32
	Uint32
	Int64
	Uint64
	Bool
	String
)

var BasicKinds = map[reflect.Kind]bool{
	reflect.Float64: true,
	reflect.Float32: true,
	reflect.Int32:   true,
	reflect.Uint32:  true,
	reflect.Int64:   true,
	reflect.Uint64:  true,
	reflect.Bool:    true,
	reflect.String:  true,
}

var BasicTypes = map[reflect.Type]bool{
	reflect.TypeOf(time.Time{}):           true,
	reflect.TypeOf(timestamp.Timestamp{}): true,
	reflect.TypeOf(sql.NullBool{}):        true,
	reflect.TypeOf(sql.NullFloat64{}):     true,
	reflect.TypeOf(sql.NullInt32{}):       true,
	reflect.TypeOf(sql.NullInt64{}):       true,
	reflect.TypeOf(sql.NullString{}):      true,
	reflect.TypeOf(sql.NullTime{}):        true,
}

var BasicTypesByName = map[Value]reflect.Type{
	Time:        reflect.TypeOf(time.Time{}),
	Timestamp:   reflect.TypeOf(timestamp.Timestamp{}),
	NullBool:    reflect.TypeOf(sql.NullBool{}),
	NullFloat64: reflect.TypeOf(sql.NullFloat64{}),
	NullInt32:   reflect.TypeOf(sql.NullInt32{}),
	NullInt64:   reflect.TypeOf(sql.NullInt64{}),
	NullString:  reflect.TypeOf(sql.NullString{}),
	NullTime:    reflect.TypeOf(sql.NullTime{}),
	Float64:     reflect.TypeOf(float64(1)),
	Float32:     reflect.TypeOf(float32(1)),
	Uint:        reflect.TypeOf(uint(1)),
	Int:         reflect.TypeOf(int(1)),
	Int32:       reflect.TypeOf(int32(1)),
	Uint32:      reflect.TypeOf(uint32(1)),
	Int64:       reflect.TypeOf(uint64(1)),
	Uint64:      reflect.TypeOf(uint64(1)),
	Bool:        reflect.TypeOf(bool(true)),
	String:      reflect.TypeOf("string"),
}
