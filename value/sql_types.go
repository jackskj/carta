package value

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

// Value represents go data types which carta supports for loading as well as what data types arrive from the sql driver
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

	// special case, carta does NOT support loading []uint8, any data that arrives from sql database as []uint8
	//is converted to bytes and expected field type is a string or *string
	Uint8Slice
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

// Map of database data types to go types
var SQLTypes = map[string]Value{
	"VARCHAR":  String,
	"TEXT":     String,
	"NVARCHAR": String,

	"DECIMAL": Float64,
	"FLOAT8":  Float64,
	"FLOAT4":  Float64,

	"BOOL": Bool,

	"INT":  Int64,
	"INT2": Int64,
	"INT4": Int64,
	"INT8": Int64,

	"TIME":        Time,
	"DATE":        Time,
	"TIMESTAMP":   Time,
	"TIMESTAMPZ":  Time,
	"TIMETZ":      Time,
	"TIMESTAMPTZ": Time,

	"NUMERIC": Uint8Slice,
	"UUID":    Uint8Slice,
	"BPCHAR":  Uint8Slice,
	"BIT":     Uint8Slice,
	"CIDR":    Uint8Slice,
	"XML":     Uint8Slice,
	"OID":     Uint8Slice,
}
