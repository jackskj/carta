package carta

import (
	"errors"
	"fmt"
	"reflect"
)

func (s *Mapper) subMapByIndex(resp reflect.Value, i int) interface{} {
	// resp must apways be *[]intefacer{}
	return nil
}

func (s *Mapper) set(dst reflect.Value, v interface{}, fieldIndex int) error {
	if err := s.FieldLoaders[fieldIndex](dst, v); err != nil {
		return err
	}
	return nil
}

func NonSlinceGrowError(typ reflect.Type) error {
	return fmt.Errorf("carta: Multiple mappings were found for non-slice type, %T. A portion of returned sql data is therefore omitted. Consider debugging your query, verifying your sql relational integrity, or chaging the type to a slice.", typ)
}

// If a new mapping has been foing, grow will instantiate a new instance our type and append it
func (s *Mapper) grow(dst reflect.Value) (interface{}, error) {
	if s.Crd == Association {
		// Grow is tricky for structs, invocation of this function where cardinality is association more then once indicates either user's mistake or broker referencial integrity
		// To explain, if a user asks to map sql response to *User, he/she expects only one user
		// However, if carta calls of this function more thank once, it indicates that carta logic determined that the sql response actually would map to many Users,
		if s.IsMapped == false {
			// ss.IsMapped = true
			// if dst.IsNil
			return dst, dst, nil
		}
		return nil, nil, NonSlinceGrowError(ss.typ)
	}

}

var NullLoad = errors.New("Null value cannot be loaded, use sql.NullX type")

func OverflowErr(i interface{}, typ reflect.Type) error {
	return fmt.Errorf("carta: value %v overflows %v", i, typ)
}

func determineBasicLoaderFunc(m *Mapper) error {
	return nil
}

func determineLoaderFuncs(m *Mapper) error {
	var (
		err error
		lf  loadFunc
	)
	if m.IsBasic {
		return determineBasicLoaderFunc(m)
	}
	for _, c := range m.PresentColumns {
		if lf, err = getLoaderFunc(m.Typ, c); err != nil {
			return err
		}
		m.FieldLoaders[c.fieldIndex] = lf
	}
	return nil
}

func getLoaderFunc(t reflect.Type, c column) (loadFunc, error) {
	isTypePtr := false
	//basic mapper logic here

	if t.Kind() == reflect.Ptr {
		isTypePtr = true
		t = t.Elem()
	}
	fieldType := t.Field(c.fieldIndex).Type()
	return nil, nil

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

func determineSetFunc(ftyp reflect.Type, sqlTyp reflect.Type) loadFunc {
	switch sqlType {
	case reflect.String:
		if ftyp.Kind() == reflect.String {
			return setString
		}
		if ftyp == sql.NullString {
			return setString
		}
	}
	return nil
}
