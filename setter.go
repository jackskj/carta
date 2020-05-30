package carta

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type mapSetter interface {
	set(dst reflect.Value, v interface{}, i int) error
	grow(dst reflect.Value) (interface{}, error)
	subMapByIndex(dst interface{}, i int) interface{}
}

type loadFunc func(v reflect.Value, dst interface{}) error

type Setter struct {
	crd Cardinality //

	listPtr bool // true *[], false [] used only if cardinality is a collection

	isMapped bool // true if  we previous SQL row has already mapped this struct, used only if cardinality is association

	// Basic loader is usef if the type is primative, sql.NullX, time.timestamp, to pointer to those
	// for example
	// type User struct {
	//        UserId   int
	//        UserAddr []sql.NullString //collection submap where setter is basic
	// }
	// it is not possibe for both isList and Is basic to be true, this would mean that we are type is myltidimentional array
	// for example [][]string, which does not make any sence for sql mapping
	isBasic     bool
	basicLoader loadFunc

	typ          reflect.Type     // Underlying type
	typePtr     bool // id the underlying type pointed to 
	fieldLoaders map[int]loadFunc // int is the ith struct field
}

func newStter(t reflect.Type) *Setter {
	var crd int
	listPtr := false
	isBasic := false
	// elemTyp :=
	if t.Kind() == reflect.Slice {
		crd = Collection
	} else if isSlicePtr(t) {
		crd = Collection
		listPtr = true
	} else if isBasicType(t) {
		isBasic := true
	}

	elemTyp := t.Elem()
	ptr := false
	if t.Kind() == reflect.Ptr {
		elemTyp = t.Elem()
		ptr = true

	}
	loaders := map[int]loadFunc{}
	for typ {

	}
	// if determineAllowedType(elemTyp) {
	// }
	// is the underlying type a struct or *struct
	if isStructPtr(t) || t.Kind() == reflect.Struct {
	}
	return &Setter{
		isMapped: false,
	}
}

// If a new mapping has been foing, grow will instantiate a new instance our type and append it
func (s *Setter) grow(dst reflect.Value) (interface{}, interface{}, error) {
	return nil, nil, nil
}

func (s *Setter) subMapByIndex(dst reflect.Value, i int) interface{} {
	return nil
}

func (ss *StructSetter) set(dst reflect.Value, v interface{}, fieldIndex int) error {
	if err := ss.Setter.loaders[fieldIndex](dst, v); err != nil {
		return err
	}
	return nil
}

func NonSlinceGrowError(typ reflect.Type) error {
	return fmt.Errorf("carta: Multiple mappings were found for non-slice type, %T. A portion of returned sql data is therefore omitted. Consider debugging your query, verifying your sql relational integrity, or chaging the type to a slice.", typ)
}

func (s *Setter) grow(dst reflect.Value) (interface{}, error) {
	if s.crd == Association {
		// Grow is tricky for structs, invocation of this function where cardinality is association more then once indicates either user's mistake or broker referencial integrity
		// To explain, if a user asks to map sql response to *User, he/she expects only one user
		// However, if carta calls of this function more thank once, it indicates that carta logic determined that the sql response actually would map to many Users,
		if ss.IsMapped == false {
			ss.IsMapped = true
			if dst.IsNil
			return dst, dst, nil
		}
		return nil, nil, NonSlinceGrowError(ss.typ)
	}
	
}

func (ss *StructSetter) subMapByIndex(dst interface{}, i int) interface{} {
	return nil
}

var NullLoad = errors.New("Null value cannot be loaded, use sql.NullX type")

func OverflowErr(i interface{}, typ reflect.Type) error {
	return fmt.Errorf("carta: value %v overflows %v", i, typ)
}

// Setter functions inspired by BQ api
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
