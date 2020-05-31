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
	PresentColumns map[string]column

	// Columns of all parents structs, used to detect whether a new struct should be appended for has-many relationships
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
	// employees_ is the prefix of the parent (lower case of the parent with "_")
	FieldNames    map[int]string
	AncestorNames []string

	// Nested structs which correspond to any has-one has-many relationships
	// int is the ith element of this struct where the submap exists
	SubMaps map[int]*Mapper
}

// Maps db rows onto the complex struct,
// Response must be a struct, pointer to a struct for our response, a slice of structs or slice of pointers to a struct
func Map(rows *sql.Rows, dst interface{}) error {
	var (
		mapper *Mapper
		err    error
	)
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	dstTyp := reflect.TypeOf(dst)
	mapper, ok := mapperCache.loadMap(columns, dstTyp)
	if ok {
		return mapper.loadRows(rows, dst)
	} else {
		if isSlicePtr(dstTyp) || isStructPtr(dstTyp) {
			return fmt.Errorf("carta: cannot map rows onto %T, destination must be pointer to a slice(*[]) or pointer to struct", dstTyp)
		}

		// generate new mapper
		if mapper, err = newMapper(dstTyp); err != nil {
			return err
		}

		// determine field names
		if err = determineFieldsNames(mapper, nil); err != nil {
			return err
		}

		// Allocate columns
		columnsByName := map[string]column{}
		for i, columnName := range columns {
			columnsByName[columnName] = column{
				typ:         columnTypes[i],
				columnIndex: i,
			}
		}
		if err = allocateColumns(mapper, columnsByName); err != nil {
			return err
		}

		if err = determineLoaderFuncs(mapper); err != nil {
			return err
		}

		mapperCache.storeMap(columns, dstTyp, mapper)
	}
	return mapper.loadRows(rows, dst)
}

func newMapper(t reflect.Type) (*Mapper, error) {
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
		Crd:       crd,
		IsListPtr: isListPtr,
		IsBasic:   isBasic,
		Typ:       elemTyp,
		IsTypePtr: isTypePtr,
	}

	mapper.SubMaps = subMaps

	return mapper, nil
}

func findSubMaps(t reflect.Type) (map[int]*Mapper, error) {
	var (
		subMap *Mapper
		err    error
	)
	subMaps := map[int]*Mapper{}
	if t.Kind() != reflect.Struct {
		return nil, nil
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if isExported(field) && isSubMap(field.Type) {
			if subMap, err = newMapper(field.Type); err != nil {
				return nil, err
			}
			subMaps[i] = subMap
		}
	}
	return subMaps, nil
}

func determineFieldsNames(m *Mapper, ancestorNames []string) error {
	var (
		t    reflect.Type
		name string
	)
	names := map[int]string{}

	m.AncestorNames = ancestorNames

	if m.IsBasic {
		return nil
	}
	if m.IsTypePtr {
		t = m.Typ.Elem()
	} else {
		t = m.Typ
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if isExported(field) {
			if tag := nameFromTag(field.Tag); tag != "" {
				name = tag
			} else {
				name = field.Name
			}
			names[i] = name
		}
	}
	m.FieldNames = names
	if ancestorNames == nil {
		ancestorNames = []string{}
	}
	for i, subMap := range m.SubMaps {
		if err := determineFieldsNames(subMap, append(ancestorNames, names[i])); err != nil {
			return err
		}
	}
	return nil
}

func isExported(f reflect.StructField) bool {
	return (f.PkgPath == "")
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

// test wether incoming type is a pointer to a struct, courtesy of BQ api
func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}
func isSlicePtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice
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
