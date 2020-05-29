package carta

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log"
	"reflect"
	"strconv"
	"time"
)

// This set of functions casts a single value from the DB
// onto a single proto field type and send the proto field.
// Inspired by the fmt package
// Neither proto filed nor db value are known types, hence I use type assertions

// TODOs:
// - reflect each column once column once, then reuse the reflection results
// - return errors for faulty conversions, for example, casting string to int
// - convert "null" responses (mssql) to nil responses

type Value struct {
	rfield   reflect.Value
	rvalue   reflect.Value
	ivalue   interface{}
	intVal   int64
	uintVal  uint64
	floatVal float64
	strVal   string
	boolVal  bool
	timeVal  reflect.Value
	err      error
}

func setProto(field reflect.Value, value interface{}) error {
	v := Value{
		rfield: field,
		rvalue: reflect.ValueOf(value),
		ivalue: value,
	}
	switch field.Interface().(type) {
	case int, int8, int16, int32, int64:
		v.castVal("int")
		field.SetInt(v.intVal)
	case uint, uint8, uint16, uint32, uint64:
		v.castVal("uint")
		field.SetUint(v.uintVal)
	case float32, float64:
		v.castVal("float")
		field.SetFloat(v.floatVal)
	case string:
		v.castVal("string")
		field.SetString(v.strVal)
	case bool:
		v.castVal("bool")
		field.SetBool(v.boolVal)
	case *timestamp.Timestamp:
		v.castVal("timestamp")
		field.Set(v.timeVal)
	default:
		//mapping enums
		if field.Kind() == reflect.Int32 {
			v.castVal("enum")
			field.SetInt(v.intVal)
		}
	}
	// TODO, return errors for more faulty conversions, for example, invalid string to int
	// As on now, errors are reurn for invalid datetime
	if v.err != nil {
		return v.err
	} else {
		return nil
	}
}

func (v *Value) castVal(respType string) {
	switch respType {
	case "int":
		v.castInt(false)
	case "enum":
		v.castInt(true)
	case "uint":
		v.castUint()
	case "float":
		v.castFLoat()
	case "string":
		v.castString()
	case "bool":
		v.castBool()
	case "timestamp":
		v.castTimestamp()
	}
}

func (v *Value) castInt(isEnum bool) {
	switch v.ivalue.(type) {
	case int, int8, int16, int32, int64:
		v.intVal = v.rvalue.Int()
	case uint, uint8, uint16, uint32, uint64:
		v.intVal = int64(v.rvalue.Uint())
	case float32, float64:
		v.intVal = int64(v.rvalue.Float())
	case string:
		if isEnum {
			valuesMapName := v.rfield.Type().Name()
			enumVal := reflect.ValueOf(v.ivalue).String()
			//Get the enums from registered proto enumse
			if intVal, found := EnumVals[valuesMapName][enumVal]; found {
				v.intVal = int64(intVal)
			} else {
				log.Println(EnumVals[valuesMapName])
				log.Println(EnumVals)
				v.err = errors.New("Value \"" + enumVal + "\" not found in " + valuesMapName + " enum.")
			}
		} else if s, err := strconv.Atoi(v.rvalue.String()); err == nil {
			v.intVal = int64(s)
		} else {
			v.intVal = int64(0)
		}
	case bool:
		if v.rvalue.Bool() {
			v.intVal = 1
		} else {
			v.intVal = 0
		}
	case time.Time:
		v.intVal = v.ivalue.(time.Time).Unix()
	}
}

func (v *Value) castUint() {
	switch v.ivalue.(type) {
	case int, int8, int16, int32, int64:
		v.uintVal = uint64(v.rvalue.Int())
	case uint, uint8, uint16, uint32, uint64:
		v.uintVal = v.rvalue.Uint()
	case float32, float64:
		v.uintVal = uint64(v.rvalue.Float())
	case string:
		if s, err := strconv.Atoi(v.rvalue.String()); err == nil {
			v.uintVal = uint64(s)
		} else {
			v.uintVal = uint64(0)
		}
	case bool:
		if v.rvalue.Bool() {
			v.uintVal = uint64(1)
		} else {
			v.uintVal = uint64(0)
		}
	case time.Time:
		v.uintVal = uint64(v.ivalue.(time.Time).Unix())
	}
}

func (v *Value) castFLoat() {
	switch v.ivalue.(type) {
	case int, int8, int16, int32, int64:
		v.floatVal = float64(v.rvalue.Int())
	case uint, uint8, uint16, uint32, uint64:
		v.floatVal = float64(v.rvalue.Uint())
	case float32, float64:
		v.floatVal = v.rvalue.Float()
	case string:
		if s, err := strconv.Atoi(v.rvalue.String()); err == nil {
			v.floatVal = float64(s)
		} else {
			v.floatVal = float64(0)
		}
	case bool:
		if v.rvalue.Bool() {
			v.floatVal = float64(1)
		} else {
			v.floatVal = float64(0)
		}
	case time.Time:
		v.floatVal = float64(v.ivalue.(time.Time).Unix())
	}
}

func (v *Value) castString() {
	switch v.ivalue.(type) {
	case int, int8, int16, int32, int64:
		v.strVal = strconv.FormatInt(v.rvalue.Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		v.strVal = strconv.FormatUint(v.rvalue.Uint(), 10)
	case float32, float64:
		v.strVal = strconv.FormatFloat(v.rvalue.Float(), 'E', -1, 64)
	case string:
		v.strVal = v.rvalue.String()
	case bool:
		if v.rvalue.Bool() {
			v.strVal = "true"
		} else {
			v.strVal = "false"
		}
	case time.Time:
		v.strVal = v.ivalue.(time.Time).String()
	}
}

func (v *Value) castBool() {
	switch v.ivalue.(type) {
	case int, int8, int16, int32, int64:
		if v.rvalue.Int() == 0 {
			v.boolVal = false
		} else {
			v.boolVal = true
		}
	case uint, uint8, uint16, uint32, uint64:
		if v.rvalue.Uint() == uint64(0) {
			v.boolVal = false
		} else {
			v.boolVal = true
		}
	case float32, float64:
		if v.rvalue.Float() == float64(0) {
			v.boolVal = false
		} else {
			v.boolVal = true
		}
	case string:
		if v.rvalue.String() == "" {
			v.boolVal = false
		} else {
			v.boolVal = true
		}
	case bool:
		v.boolVal = v.rvalue.Bool()
	case time.Time:
		v.boolVal = true
	}
}

func (v *Value) castTimestamp() {
	switch v.ivalue.(type) {
	case time.Time:
		if sqlTime, err := ptypes.TimestampProto(v.ivalue.(time.Time)); err == nil {
			v.timeVal = reflect.ValueOf(sqlTime)
		} else {
			v.timeVal = reflect.ValueOf(ptypes.TimestampNow())
			v.err = err
		}
	default:
		v.timeVal = reflect.ValueOf(ptypes.TimestampNow())
		v.err = errors.New(fmt.Sprintf(
			"cannot convert %s of type %s to time.Time",
			v.ivalue, reflect.TypeOf(v.ivalue),
		))
	}
}
