package value

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// TODO:  timestamp/time/ from string
// TODO:  int/float/uint/bool from string

// var Types = map[string]reflect.Type{}

type Cell struct {
	kind   reflect.Kind // data type with which Cell will be instantiated
	bits   uint64       //IEEE 754 binary representation of numeric value
	binary string       // non-numeric data as bytes for data which arrives as string
	time   time.Time    //  any data that arrives as time, that includes timestame w/ or w/o zone
	valid  bool
}

var NullSet = errors.New("Null value cannot be loaded, use sql.NullX type")

func OverflowErr(i interface{}, typ reflect.Type) error {
	return fmt.Errorf("carta: value %v overflows %v", i, typ)
}

func NewCell(v interface{}, typ *sql.ColumnType) (Cell, error) {
	var c Cell
	isSet := false
	// Types[typ.DatabaseTypeName()] = reflect.TypeOf(v)
	switch SQLTypes[typ.DatabaseTypeName()] {
	case String:
		if data, ok := v.(string); ok {
			isSet = true
			c = NewString(data)
		}
	case Float64:
		if data, ok := v.(float64); ok {
			isSet = true
			c = NewFloat64(data)
		}
	case Bool:
		if data, ok := v.(bool); ok {
			isSet = true
			c = NewBool(data)
		}
	case Int64:
		if data, ok := v.(int64); ok {
			isSet = true
			c = NewInt64(data)
		}
	case Time:
		if data, ok := v.(time.Time); ok {
			isSet = true
			c = NewTime(data)
		}
	case Uint8Slice:
		if data, ok := v.([]uint8); ok {
			isSet = true
			c = NewString(string(data))
		}
	}
	if !isSet {
		// for x, y := range Types {
		// log.Printf("%v: %v", x, y)
		// }
		return Cell{}, errors.New(fmt.Sprintf("carta: unknown data type %s for column %s ", typ.DatabaseTypeName(), typ.Name()))
	}
	return c, nil
}

func NewBool(c bool) Cell {
	if c {
		return Cell{
			kind:  reflect.Bool,
			bits:  1,
			valid: true,
		}
	}
	return Cell{
		kind:  reflect.Bool,
		bits:  0,
		valid: true,
	}
}

func NewFloat32(c float32) Cell {
	return Cell{
		kind:  reflect.Float32,
		bits:  uint64(math.Float32bits(c)),
		valid: true,
	}
}

func NewFloat64(c float64) Cell {
	return Cell{
		kind:  reflect.Float64,
		bits:  math.Float64bits(c),
		valid: true,
	}
}

func NewInt32(c int32) Cell {
	return Cell{
		kind:  reflect.Int32,
		bits:  uint64(c),
		valid: true,
	}
}

func NewUint32(c uint32) Cell {
	return Cell{
		kind:  reflect.Uint32,
		bits:  uint64(c),
		valid: true,
	}
}

func NewInt64(c int64) Cell {
	return Cell{
		kind:  reflect.Int64,
		bits:  uint64(c),
		valid: true,
	}
}

func NewUint64(c uint64) Cell {
	return Cell{
		kind:  reflect.Uint64,
		bits:  c,
		valid: true,
	}
}

func NewString(c string) Cell {
	return Cell{
		kind:   reflect.String,
		binary: c,
		valid:  true,
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
	return Cell{
		kind:  reflect.Struct,
		time:  c,
		valid: true,
	}
}

func NewNull() Cell {
	return Cell{}
}

func (c Cell) Kind() reflect.Kind {
	return c.kind
}

func (c Cell) IsNull() bool {
	return !c.valid
}

func (c Cell) IsValid() bool {
	return c.valid
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
	return c.time
}

func (c Cell) Timestamp() timestamp.Timestamp {
	if t, err := ptypes.TimestampProto(c.Time()); err == nil {
		return *t
	}
	// should not happen
	return timestamp.Timestamp{}
}

func (c Cell) NullBool() sql.NullBool {
	if !c.valid {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  c.Bool(),
		Valid: true,
	}
}

func (c Cell) NullFloat64() sql.NullFloat64 {
	if !c.valid {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: c.Float64(),
		Valid:   true,
	}
}

func (c Cell) NullInt32() sql.NullInt32 {
	if !c.valid {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: c.Int32(),
		Valid: true,
	}
}

func (c Cell) NullInt64() sql.NullInt64 {
	if !c.valid {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: c.Int64(),
		Valid: true,
	}
}

func (c Cell) NullString() sql.NullString {
	if !c.valid {
		return sql.NullString{}
	}
	return sql.NullString{
		String: c.String(),
		Valid:  true,
	}
}

func (c Cell) NullTime() sql.NullTime {
	if !c.valid {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  c.Time(),
		Valid: true,
	}
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
