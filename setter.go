package carta

import (
	"errors"
	"fmt"
	"reflect"
)

// If a new mapping has been foing, grow will instantiate a new instance our type and append it
func (s *Mapper) grow(dst reflect.Value) (interface{}, interface{}, error) {
	return nil, nil, nil
}

func (s *Mapper) subMapByIndex(resp reflect.Value, i int) interface{} {
	// resp must apways be *[]intefacer{}
	return nil
}

func (s *Mapper) set(dst reflect.Value, v interface{}, fieldIndex int) error {
	if err := ss.Mapper.loaders[fieldIndex](dst, v); err != nil {
		return err
	}
	return nil
}

func NonSlinceGrowError(typ reflect.Type) error {
	return fmt.Errorf("carta: Multiple mappings were found for non-slice type, %T. A portion of returned sql data is therefore omitted. Consider debugging your query, verifying your sql relational integrity, or chaging the type to a slice.", typ)
}

func (s *Mapper) grow(dst reflect.Value) (interface{}, error) {
	if s.crd == Association {
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

func (ss *StructMapper) subMapByIndex(dst interface{}, i int) interface{} {
	return nil
}

var NullLoad = errors.New("Null value cannot be loaded, use sql.NullX type")

func OverflowErr(i interface{}, typ reflect.Type) error {
	return fmt.Errorf("carta: value %v overflows %v", i, typ)
}

func determineBasicLoaderFunc(m *Mapper) error {
	return nil
}

func determineLoaderFuncs(m *Mapper) error {
	if m.IsBasic {
		return determineBasicLoaderFunc(m)
	}
	for _, c := range m.PresentColumns {
		m.FieldLoaders[c.fieldIndex] = getLoaderFunc()
	}
	return nil
}

// Mapper functions inspired by BQ api
func setInt(v reflect.Value, x interface{}) error {
	if x == nil {
		return NullLoad
	}
	xx := x.(int64)
	if v.OverflowInt(xx) {
		return OverflowErr(xx, v.Type())
	}
	v.SetInt(xx)
	return nil
}

func setUint(v reflect.Value, x interface{}) error {
	if x == nil {
		return NullLoad
	}
	xx := x.(int64)
	if xx < 0 || v.OverflowUint(uint64(xx)) {
		return OverflowErr(xx, v.Type())
	}
	v.SetUint(uint64(xx))
	return nil
}

func setFloat(v reflect.Value, x interface{}) error {
	if x == nil {
		return NullLoad
	}
	xx := x.(float64)
	if v.OverflowFloat(xx) {
		return OverflowErr(xx, v.Type())
	}
	v.SetFloat(xx)
	return nil
}

func setBool(v reflect.Value, x interface{}) error {
	if x == nil {
		return NullLoad
	}
	v.SetBool(x.(bool))
	return nil
}

func setString(v reflect.Value, x interface{}) error {
	if x == nil {
		return NullLoad
	}
	v.SetString(x.(string))
	return nil
}

func setBytes(v reflect.Value, x interface{}) error {
	if x == nil {
		v.SetBytes(nil)
	} else {
		v.SetBytes(x.([]byte))
	}
	return nil
}

func setNull(v reflect.Value, x interface{}, build func() interface{}) error {
	if x == nil {
		v.Set(reflect.Zero(v.Type()))
	} else {
		n := build()
		v.Set(reflect.ValueOf(n))
	}
	return nil
}

func setNested(ops []structLoaderOp, v reflect.Value, val interface{}) error {
	// v is either a struct or a pointer to a struct.
	if v.Kind() == reflect.Ptr {
		// If the value is nil, set the pointer to nil.
		if val == nil {
			v.Set(reflect.Zero(v.Type()))
			return nil
		}
		// If the pointer is nil, set it to a zero struct value.
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return runOps(ops, v, val.([]Value))
}

func setRepeated(field reflect.Value, vslice []Value, setElem setFunc) error {
	vlen := len(vslice)
	var flen int
	switch field.Type().Kind() {
	case reflect.Slice:
		// Make a slice of the right size, avoiding allocation if possible.
		switch {
		case field.Len() < vlen:
			field.Set(reflect.MakeSlice(field.Type(), vlen, vlen))
		case field.Len() > vlen:
			field.SetLen(vlen)
		}
		flen = vlen

	case reflect.Array:
		flen = field.Len()
		if flen > vlen {
			// Set extra elements to their zero value.
			z := reflect.Zero(field.Type().Elem())
			for i := vlen; i < flen; i++ {
				field.Index(i).Set(z)
			}
		}
	default:
		return fmt.Errorf("bigquery: impossible field type %s", field.Type())
	}
	for i, val := range vslice {
		if i < flen { // avoid writing past the end of a short array
			if err := setElem(field.Index(i), val); err != nil {
				return err
			}
		}
	}
	return nil
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
