package value

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	// "time
	// "github.com/golang/protobuf/ptypes/timestamp"
)

var _ = log.Fatal

// If a new mapping has been foind, grow will instantiate a new instance our type and append it

var NullLoad = errors.New("Null value cannot be loaded, use sql.NullX type")

func OverflowErr(i interface{}, typ reflect.Type) error {
	return fmt.Errorf("carta: value %v overflows %v", i, typ)
}

// func determineBasicLoaderFunc(m *Mapper) error {
// return nil
// }

func SetField(field reflect.Value, c Cell) {

}
