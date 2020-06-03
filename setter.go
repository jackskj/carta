package carta

import (
	"errors"
	"fmt"
	"reflect"
	// "time
	// "github.com/golang/protobuf/ptypes/timestamp"
)

// If a new mapping has been foind, grow will instantiate a new instance our type and append it
func (m *Mapper) grow(dst reflect.Value) reflect.Value {
	if m.Crd == Association {
		if m.IsTypePtr {
			reflect.Indirect(dst).Set(reflect.Zero(m.Typ))
		}
		return dst
	} else if m.Crd == Collection {
		// if cardinaloty is a collection, the arrading destination MUST be a pointer to a slice
		var newDst reflect.Value
		if m.IsTypePtr {
			newDst = reflect.New(m.Typ)
		} else {
			newDst = reflect.Zero(m.Typ)
		}
		indirectSlice := reflect.Indirect(dst)
		indirectSlice = reflect.Append(indirectSlice, newDst)
		if m.IsTypePtr {
			return newDst
		} else {
			return newDst.Addr()
		}
	}
	// should never happen
	return dst
}

var NullLoad = errors.New("Null value cannot be loaded, use sql.NullX type")

func OverflowErr(i interface{}, typ reflect.Type) error {
	return fmt.Errorf("carta: value %v overflows %v", i, typ)
}

func determineBasicLoaderFunc(m *Mapper) error {
	return nil
}

var dbTypes = map[string]reflect.Kind{
	"VARCHAR":  reflect.String,
	"TEXT":     reflect.String,
	"NVARCHAR": reflect.String,
	"DECIMAL":  reflect.Float64,
	"BOOL":     reflect.Bool,
	"INT":      reflect.Int32,
	"BIGINT":   reflect.Int64,
}
