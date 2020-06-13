package value

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// TODO:  timestamp/time/ from string
// TODO:  int/float/uint/bool from string

var Types = map[string]reflect.Type{}

type Cell struct {
	kind       reflect.Kind // data type with which Cell will be instantiated
	bits       uint64       //IEEE 754 binary representation of numeric value
	text       string       // non-numeric data as bytes for data which arrives as string or []byte
	time       time.Time    //  any data that arrives as time, that includes timestame w/ or w/o zone
	colTypName string       // Used for parting if some data arrices in plain text format, ex, if time arrives as string
	valid      bool
}

func OverflowErr(i interface{}, typ reflect.Type) error {
	return fmt.Errorf("carta: value %v overflows %v", i, typ)
}

func ConvertsionError(convErr error, typ reflect.Type) error {
	return fmt.Errorf("carta: errors converting to %v: "+convErr.Error(), typ)
}

func NewCell(colTypName string) *Cell {
	return &Cell{colTypName: colTypName}
}

// implements database/sql scan interface
func (c *Cell) Scan(src interface{}) error {
	switch src.(type) {
	case int64:
		c.SetInt64(src.(int64))
	case float64:
		c.SetFloat64(src.(float64))
	case bool:
		c.SetBool(src.(bool))
	case []byte:
		c.SetString(string(src.([]byte)))
	case string:
		c.SetString(src.(string))
	case time.Time:
		c.SetTime(src.(time.Time))
	default:
		// src is nil
		c.SetNull()
	}
	return nil
}

func (c *Cell) SetBool(d bool) {
	c.kind = reflect.Bool
	c.valid = true
	if d {
		c.bits = 1
	} else {
		c.bits = 0
	}
}

func (c *Cell) SetFloat64(d float64) {
	c.kind = reflect.Float64
	c.valid = true
	c.bits = math.Float64bits(d)
}

func (c *Cell) SetInt64(d int64) {
	c.kind = reflect.Float64
	c.valid = true
	c.bits = uint64(d)
}

func (c *Cell) SetString(d string) {
	c.kind = reflect.String
	c.valid = true
	c.text = d
}

func (c *Cell) SetTime(d time.Time) {
	c.kind = reflect.Struct
	c.valid = true
	c.time = d
}

func (c *Cell) SetNull() {
	c.valid = false
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

func (c Cell) Bool() (bool, error) {
	return (c.bits != 0), nil
}

func (c Cell) Int32() (int32, error) {
	if c.kind == reflect.String {
		if num, err := strconv.ParseInt(c.text, 10, 32); err != nil {
			return 0, err
		} else {
			return int32(num), nil
		}
	}
	return int32(c.bits), nil
}

func (c Cell) Int64() (int64, error) {
	if c.kind == reflect.String {
		if num, err := strconv.ParseInt(c.text, 10, 64); err != nil {
			return 0, err
		} else {
			return int64(num), nil
		}
	}
	return int64(c.bits), nil
}

func (c Cell) Uint32() (uint32, error) {
	if c.kind == reflect.String {
		if num, err := strconv.ParseUint(c.text, 10, 32); err != nil {
			return 0, err
		} else {
			return uint32(num), nil
		}
	}
	return uint32(c.bits), nil
}

func (c Cell) Uint64() (uint64, error) {
	if c.kind == reflect.String {
		if num, err := strconv.ParseUint(c.text, 10, 64); err != nil {
			return 0, err
		} else {
			return uint64(num), nil
		}
	}
	return c.bits, nil
}

func (c Cell) Float32() (float32, error) {
	if c.kind == reflect.String {
		if num, err := strconv.ParseFloat(c.text, 32); err != nil {
			return 0, err
		} else {
			return float32(num), nil
		}
	}
	return math.Float32frombits(uint32(c.bits)), nil
}

func (c Cell) Float64() (float64, error) {
	if c.kind == reflect.String {
		if num, err := strconv.ParseFloat(c.text, 64); err != nil {
			return 0, err
		} else {
			return num, nil
		}
	}
	return math.Float64frombits(c.bits), nil
}

func (c Cell) String() (string, error) {
	return c.text, nil
}

func (c Cell) Time() (time.Time, error) {
	if c.kind == reflect.String {
		// TODO: Parse from string
		// switch c.colTypName {
		// }
		return time.Time{}, errors.New("cannot convert time data which arrived as string or []uint8 from sql")
	}
	return c.time, nil
}

func (c Cell) Timestamp() (timestamp.Timestamp, error) {
	var t time.Time
	var err error
	if t, err = c.Time(); err != nil {
		return timestamp.Timestamp{}, err
	}
	if ts, err := ptypes.TimestampProto(t); err != nil {
		return timestamp.Timestamp{}, err
	} else {
		return *ts, nil
	}
}

func (c Cell) NullBool() (sql.NullBool, error) {
	if !c.valid {
		return sql.NullBool{}, nil
	}
	d, err := c.Bool()
	return sql.NullBool{
		Bool:  d,
		Valid: true,
	}, err
}

func (c Cell) NullFloat64() (sql.NullFloat64, error) {
	if !c.valid {
		return sql.NullFloat64{}, nil
	}
	d, err := c.Float64()
	return sql.NullFloat64{
		Float64: d,
		Valid:   true,
	}, err
}

func (c Cell) NullInt32() (sql.NullInt32, error) {
	if !c.valid {
		return sql.NullInt32{}, nil
	}
	d, err := c.Int32()
	return sql.NullInt32{
		Int32: d,
		Valid: true,
	}, err
}

func (c Cell) NullInt64() (sql.NullInt64, error) {
	if !c.valid {
		return sql.NullInt64{}, nil
	}
	d, err := c.Int64()
	return sql.NullInt64{
		Int64: d,
		Valid: true,
	}, err
}

func (c Cell) NullString() (sql.NullString, error) {
	if !c.valid {
		return sql.NullString{}, nil
	}
	d, err := c.String()
	return sql.NullString{
		String: d,
		Valid:  true,
	}, err
}

func (c Cell) NullTime() (sql.NullTime, error) {
	if !c.valid {
		return sql.NullTime{}, nil
	}
	d, err := c.Time()
	return sql.NullTime{
		Time:  d,
		Valid: true,
	}, err
}

func (c Cell) AsInterface() (interface{}, error) {
	var i interface{}
	var err error
	switch c.kind {
	case reflect.Bool:
		i, err = c.Bool()
	case reflect.Int32:
		i, err = c.Int32()
	case reflect.Int64:
		i, err = c.Int64()
	case reflect.Uint32:
		i, err = c.Uint32()
	case reflect.Uint64:
		i, err = c.Uint64()
	case reflect.Float32:
		i, err = c.Float32()
	case reflect.Float64:
		i, err = c.Float64()
	case reflect.String:
		i, err = c.String()
	}
	return i, err
}

func (c Cell) Uid() string {
	if c.IsNull() {
		//TODO: safely represent null and bool values as string
		return "cnull"
	} else {
		switch c.kind {
		case reflect.Int64, reflect.Float64:
			return strconv.FormatUint(c.bits, 36)
		case reflect.String:
			return c.text
		case reflect.Bool:
			if c.bits != 0 {
				return "ctrue"
			} else {
				return "cfalse"
			}
		case reflect.Struct:
			return strconv.FormatInt(c.time.Unix(), 36)
		}
	}
	return ""
}

// func (c Cell) BitsAsString() string {
// switch c.kind {
// case reflect.Bool:
// return strconv.FormatBool(c.Bool())
// case reflect.Int32:
// return strconv.FormatInt(int64(c.Int32()), 10)
// case reflect.Int64:
// return strconv.FormatInt(c.Int64(), 10)
// case reflect.Uint32:
// return strconv.FormatUint(uint64(c.Uint32()), 10)
// case reflect.Uint64:
// return strconv.FormatUint(c.Uint64(), 10)
// case reflect.Float32:
// return fmt.Sprint(c.Float32())
// case reflect.Float64:
// return fmt.Sprint(c.Float64())
// default:
// return ""
// }
// }
