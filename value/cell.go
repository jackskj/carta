package value

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/golang/protobuf/ptypes/timestamp"
)

// TODO:  timestamp/time/ from string
// NillXXX  from bits or binary

// TODOLater:  int/float/uint/bool from string

var _ = log.Fatal

type Cell struct {
	kind   reflect.Kind // data type with which Cell will be instantiated
	bits   uint64       //IEEE 754 binary representation of numeric value
	binary string       // non-numeric data as bytes
	isNull bool
}

var NullSet = errors.New("NULL value cannot be set")

func NewCell(v interface{}, typ *sql.ColumnType) (Cell, error) {
	var c Cell
	isSet := false

	switch typ.DatabaseTypeName() {
	case "VARCHAR", "TEXT", "NVARCHAR":
		if data, ok := v.(string); ok {
			isSet = true
			c = NewString(data)
		}
	case "DECIMAL":
		c = NewFloat64(v.(float64))
	case "BOOL":
		c = NewBool(v.(bool))
	case "INT":
		c = NewInt(v.(int))
	case "INT4":
		if data, ok := v.(int64); ok {
			isSet = true
			c = NewInt64(data)
		}
	case "BIGINT":
		c = NewInt(v.(int))
	}
	if !isSet {
		return Cell{}, errors.New(fmt.Sprintf("carta: unknown data type %s for column %s ", typ.DatabaseTypeName(), typ.Name()))
	}
	return c, nil

}

func NewBool(c bool) Cell {
	if c {
		return Cell{
			kind: reflect.Bool,
			bits: 1,
		}
	}
	return Cell{
		kind: reflect.Bool,
		bits: 0,
	}
}

func NewFloat32(c float32) Cell {
	return Cell{
		kind: reflect.Float32,
		bits: uint64(math.Float32bits(c)),
	}
}

func NewFloat64(c float64) Cell {
	return Cell{
		kind: reflect.Float64,
		bits: math.Float64bits(c),
	}
}

func NewInt32(c int32) Cell {
	return Cell{
		kind: reflect.Int32,
		bits: uint64(c),
	}
}

func NewUint32(c uint32) Cell {
	return Cell{
		kind: reflect.Uint32,
		bits: uint64(c),
	}
}

func NewInt64(c int64) Cell {
	return Cell{
		kind: reflect.Int64,
		bits: uint64(c),
	}
}

func NewUint64(c uint64) Cell {
	return Cell{
		kind: reflect.Uint64,
		bits: c,
	}
}

func NewString(c string) Cell {
	return Cell{
		kind:   reflect.String,
		binary: c,
	}
}

func NewInt(c int) Cell {
	if unsafe.Sizeof(c) == 4 {
		return NewInt32(int32(c))
	}
	return NewInt64(int64(c))
}

func NewUint(c uint) Cell {
	if unsafe.Sizeof(c) == 4 {
		return NewUint32(uint32(c))
	}
	return NewUint64(uint64(c))
}

func NewTime(c time.Time) Cell {
	return Cell{}
}

func NewTimestamp(c timestamp.Timestamp) Cell {
	return Cell{}
}

func NewNull() Cell {
	return Cell{isNull: true}
}

func (c Cell) Kind() reflect.Kind {
	return c.kind
}

func (c Cell) Bool() bool {
	return (c.bits != 0)
}

func (c Cell) Int32() int32 {
	return int32(c.bits)
}

func (c Cell) Int64() int64 {
	return int64(c.bits)
}

func (c Cell) Uint32() uint32 {
	return uint32(c.bits)
}

func (c Cell) Uint64() uint64 {
	return c.bits
}

func (c Cell) Float32() float32 {
	return math.Float32frombits(uint32(c.bits))
}

func (c Cell) Float64() float64 {
	return math.Float64frombits(c.bits)
}

func (c Cell) String() string {
	return c.binary
}

func (c Cell) Time() time.Time {
	return time.Time{}
}

func (c Cell) Timestamp() timestamp.Timestamp {
	return timestamp.Timestamp{}
}

func (c Cell) AsInterface() interface{} {
	switch c.kind {
	case reflect.Bool:
		return c.Bool()
	case reflect.Int32:
		return c.Int32()
	case reflect.Int64:
		return c.Int64()
	case reflect.Uint32:
		return c.Uint32()
	case reflect.Uint64:
		return c.Uint64()
	case reflect.Float32:
		return c.Float32()
	case reflect.Float64:
		return c.Float64()
	case reflect.String:
		return c.binary
	}
	return nil
}

func (c Cell) BitsAsString() string {
	switch c.kind {
	case reflect.Bool:
		return strconv.FormatBool(c.Bool())
	case reflect.Int32:
		return strconv.FormatInt(int64(c.Int32()), 10)
	case reflect.Int64:
		return strconv.FormatInt(c.Int64(), 10)
	case reflect.Uint32:
		return strconv.FormatUint(uint64(c.Uint32()), 10)
	case reflect.Uint64:
		return strconv.FormatUint(c.Uint64(), 10)
	case reflect.Float32:
		return fmt.Sprint(c.Float32())
	case reflect.Float64:
		return fmt.Sprint(c.Float64())
	default:
		return ""
	}
}

func (c Cell) SetInt(v reflect.Value) error {
	if c.isNull {
		return NullSet
	}
	if v.OverflowInt(c.Int64()) {
		return OverflowErr(c.Int64(), v.Type())
	}
	v.SetInt(c.Int64())
	return nil
}

func (c Cell) SetUint(v reflect.Value) error {
	if c.isNull {
		return NullSet
	}
	if v.OverflowUint(c.Uint64()) {
		return OverflowErr(c.Uint64(), v.Type())
	}
	v.SetUint(c.Uint64())
	return nil
}

func (c Cell) SetFloat(v reflect.Value) error {
	if c.isNull {
		return NullSet
	}
	if v.OverflowFloat(c.Float64()) {
		return OverflowErr(c.Float64(), v.Type())
	}
	v.SetFloat(c.Float64())
	return nil
}

func (c Cell) SetBool(v reflect.Value) error {
	if c.isNull {
		return NullSet
	}
	v.SetBool(c.Bool())
	return nil
}

func (c Cell) SetString(v reflect.Value) error {
	if c.isNull {
		return NullSet
	}
	v.SetString(c.String())
	return nil
}

func (c Cell) SetTime(v reflect.Value) error {
	return nil
}

func (c Cell) SetTimestamp(v reflect.Value) error {
	return nil
}
func (c Cell) SetNullBool(v reflect.Value) error {
	return nil
}
func (c Cell) SetNullFloat64(v reflect.Value) error {
	return nil
}
func (c Cell) SetNullInt32(v reflect.Value) error {
	return nil
}
func (c Cell) SetNullInt64(v reflect.Value) error {
	return nil
}
func (c Cell) SetNullString(v reflect.Value) error {
	return nil
}
func (c Cell) SetNullTime(v reflect.Value) error {
	return nil
}

// func determineConvertFunc(m *Mapper) error {
// var conv convertFunc
// for _, c := range m.PresentColumns {
// sourceKind := getColumnGoType(c.typ)
// conv = func(v interface{}) (interface{}, error) {
// var c Cell
// switch sourceKind {
// case reflect.Bool:
// c = NewBool(v.(bool))
// case reflect.Float32:
// c = NewFloat32(v.(float32))
// case reflect.Float64:
// c = NewFloat64(v.(float64))
// case reflect.Int32:
// c = NewInt32(v.(int32))
// case reflect.Uint32:
// c = NewUint32(v.(uint32))
// case reflect.Int64:
// c = NewInt64(v.(int64))
// case reflect.Uint64:
// c = NewUint64(v.(uint64))
// case reflect.String:
// c = NewString(v.(string))
// case reflect.Int:
// c = NewInt(v.(int))
// case reflect.Uint:
// c = NewUint(v.(uint))
// }
// if dstKind != reflect.Struct {
// switch dstKind {
// case reflect.Bool:
// return c.Bool(), nil
// case reflect.Float32:
// return c.Float32(), nil
// case reflect.Float64:
// return c.Float64(), nil
// case reflect.Int32:
// return c.Int32(), nil
// case reflect.Uint32:
// return c.Uint32(), nil
// case reflect.Int64:
// return c.Int64(), nil
// case reflect.Uint64:
// return c.Uint64(), nil
// case reflect.String:
// return c.String(), nil
// case reflect.Int:
// return c.Int(), nil
// case reflect.Uint:
// return c.Uint(), nil
// }
// } else {
// switch dstTyp {
// case basicTypesByName["Time"]:
// return c.Time()
// case basicTypesByName["Timestamp"]:
// return c.Timestamp()
// case basicTypesByName["NullBool"]:
// return c.NullBool()
// case basicTypesByName["NullFloat64"]:
// return c.NullFloat64()
// case basicTypesByName["NullInt32"]:
// return c.NullInt32()
// case basicTypesByName["NullInt64"]:
// return c.NullInt64()
// case basicTypesByName["NullString"]:
// return c.NullString()
// case basicTypesByName["NullTime"]:
// return c.NullTime()
// }
// }
// return nil, fmt.Errorf("carta: error canverting from %T to %T", sourceTyp, dstTyp)
// }
// if m.IsBasic {
// m.BasicConverter = conv
// } else {
// m.Converters = append(m.Converters, conv)
// }
// }
// return nil
// }
