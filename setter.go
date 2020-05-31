package carta

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

const (
	CartaTagKey string = "carta"
)

// SQL Map cardinality can either be:
// Association: has-one relationship, must be nested structs in the response
// Collection: had-many relationship, repeated (slice, array) nested struct or pointer to it
type Cardinality int

const (
	Unknown Cardinality = iota
	Association
	Collection
)

type loadFunc func(v reflect.Value, dst interface{}) error

type Mapper struct {
	Crd Cardinality //

	IsListPtr bool // true if destination is *[], false if destination is [], used only if cardinality is a collection

	// Basic mapper is used for collections where underlying type is basic (any field that is able to be set, look at isBasicType for more deatils )
	// for example
	// type User struct {
	//        UserId    int
	//        UserAddr  []sql.NullString // collection submap where mapper is basic
	//        UserPhone []string         // also basic mapper
	//        UserStuff *[]*string       // also basic mapper
	//        UserBlog  []*Blog          // this is NOT a basic mapper
	// }
	// basic can only be true if cardinality is collection
	IsBasic     bool
	BasicLoader loadFunc

	Typ          reflect.Type     // Underlying type to be mapped
	IsTypePtr    bool             // is the underlying type pointed to
	FieldLoaders map[int]loadFunc // setters for each fields, int is the i'th struct field

	// Columns of the SQL response which are present in this struct
	// int represents the ith struct field of this mapper where the column is to be mapped
	PresentColumns map[string]int

	// Columns of all parents structs, used to detect whether a new struct should be appended for has-many relationships
	// order is not nececary
	AncestorColumns map[string]bool

	// when reusing the same struct multiple times, you are able to specify the colimn prefix using parent structs
	// example
	// type Employee struct {
	// 	Id int
	// }
	// type Manager struct {
	// 	Employee
	// 	Employees []Employee
	// }
	// the following querry would correctly map if we were mapping to *[]Manager
	// "select id, employees_id from employees join managers"
	// employees_ is the prefix of the parent
	FieldNames    map[string]int
	AncestorNames []string

	// Nested structs which correspond to any has-one has-many relationships
	// int is the ith element of this struct where the submap exists
	SubMaps map[int]*Mapper
}

func newMapper(t reflect.Type, ancestorNames []string) (*Mapper, error) {
	var (
		crd     Cardinality
		elemTyp reflect.Type
		mapper  *Mapper
		subMaps map[int]*Mapper
		err     error
	)

	isListPtr := false
	isBasic := false
	isTypePtr := false

	if isSlicePtr(t) {
		crd = Collection
		elemTyp = t.Elem() // []interface{} to intetrface{}
		isListPtr = true
	} else if t.Kind() == reflect.Slice {
		crd = Collection
		elemTyp = t.Elem().Elem() // *[]interface{} to intetrface{}

	}

	if crd == Collection {
		isBasic = isBasicType(t)
	}

	if isStructPtr(t) {
		crd = Association
		elemTyp = t.Elem()
		isTypePtr = true
	} else if t.Kind() == reflect.Struct {
		crd = Association
		elemTyp = t
	}

	if crd == Unknown {
		return nil, errors.New("carts: unknown mapping")
	}

	mapper = &Mapper{
		Crd:           crd,
		IsListPtr:     isListPtr,
		IsBasic:       isBasic,
		Typ:           elemTyp,
		IsTypePtr:     isTypePtr,
		AncestorNames: ancestorNames,
	}

	if isBasic {
		mapper.BasicLoader = determineBasicLoaderFunc(elemTyp)
		return mapper, nil
	}

	mapper.FieldLoaders = determineLoaderFuncs(elemTyp)
	if subMaps, err = findSubMaps(elemTyp); err != nil {
		return nil, err
	}

	mapper.SubMaps = subMaps

	return mapper, nil
}

func findSubMaps(t reflect.Type, ancestorNames []string) (map[int]*Mapper, error) {
	var (
		subMap *Mapper
		err    error
		name   string
	)
	subMaps := map[int]*Mapper{}
	if t.Kind() != reflect.Struct {
		return nil, nil
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if isExported(field) && isSubMap(field.Type) {
			if tag := nameFromTag(field.Tag); tag != "" {
				name = tag
			} else {
				name = field.Name
			}
			if subMap, err = newMapper(field.Type, append(ancestorNames, name)); err != nil {
				return nil, err
			}
			subMaps[i] = subMap
		}
	}
	return subMaps, nil
}

func isExported(f reflect.StructField) bool {
	return (f.PkgPath == "")
}

// If a new mapping has been foing, grow will instantiate a new instance our type and append it
func (s *Mapper) grow(dst reflect.Value) (interface{}, interface{}, error) {
	return nil, nil, nil
}

func (s *Mapper) subMapByIndex(resp reflect.Value, i int) interface{} {
	// resp must apways be *[]intefacer{}
	return nil
}

func (ss *StructMapper) set(dst reflect.Value, v interface{}, fieldIndex int) error {
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
		if ss.IsMapped == false {
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

func nameFromTag(t reflect.StructTag) string {
	// s := t.Get(CartaTagKey)
	return ""
}

func isSubMap(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return (!isBasicType(t) && (t.Kind() == reflect.Struct))
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

// Basic types are any types that are intended to be set from sql row data
// Primative fields, sql.NullXXX, time.Time, pg timestamp qualify as basic
func isBasicType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if _, ok := basicKinds[t.Kind()]; ok {
		return true
	}
	if _, ok := basicTypes[t]; ok {
		return true
	}
	return false
}

var basicKinds = map[reflect.Kind]bool{
	reflect.Float64: true,
	reflect.Float32: true,
	reflect.Int32:   true,
	reflect.Uint32:  true,
	reflect.Int64:   true,
	reflect.Uint64:  true,
	reflect.Bool:    true,
	reflect.String:  true,
}

var basicTypes = map[reflect.Type]bool{
	reflect.TypeOf(time.Time{}):           true,
	reflect.TypeOf(timestamp.Timestamp{}): true,
	reflect.TypeOf(sql.NullBool{}):        true,
	reflect.TypeOf(sql.NullFloat64{}):     true,
	reflect.TypeOf(sql.NullInt32{}):       true,
	reflect.TypeOf(sql.NullInt64{}):       true,
	reflect.TypeOf(sql.NullString{}):      true,
	reflect.TypeOf(sql.NullTime{}):        true,
}
