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

var BasicKinds = map[reflect.Kind]Value{
	reflect.Float64: Float64,
	reflect.Float32: Float32,
	reflect.Int32:   Int32,
	reflect.Uint32:  Uint32,
	reflect.Int64:   Int64,
	reflect.Uint64:  Uint64,
	reflect.Bool:    Bool,
	reflect.String:  String,
}

var BasicTypes = map[reflect.Type]Value{
	reflect.TypeOf(time.Time{}):           Time,
	reflect.TypeOf(timestamp.Timestamp{}): Timestamp,
	reflect.TypeOf(sql.NullBool{}):        NullBool,
	reflect.TypeOf(sql.NullFloat64{}):     NullFloat64,
	reflect.TypeOf(sql.NullInt32{}):       NullInt32,
	reflect.TypeOf(sql.NullInt64{}):       NullInt64,
	reflect.TypeOf(sql.NullString{}):      NullString,
	reflect.TypeOf(sql.NullTime{}):        NullTime,
}
