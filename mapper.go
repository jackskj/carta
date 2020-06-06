package carta

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/jackskj/carta/value"
)

var _ = json.Marshal

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

type loadFunc func(src reflect.Value, dst interface{}) error
type convertFunc func(v interface{}) (reflect.Value, error)

type Field struct {
	Name string
	Typ  reflect.Type
	Kind reflect.Kind
}

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
	IsBasic bool
	// BasicConverter convertFunc

	Typ  reflect.Type // Underlying type to be mapped
	Kind reflect.Kind // Underlying Kind to be mapped

	IsTypePtr bool // is the underlying type pointed to
	// Converters map[int]convertFunc // setters for each fields, int is the i'th struct field

	// Columns of the SQL response which are present in this struct
	PresentColumns      map[string]column
	SortedColumnIndexes []int

	// Columns of all parents structs, used to detect whether a new struct should be appended for has-many relationships
	AncestorColumns map[string]column
	// TODO: SortedAncestorColumns []int

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
	// FieldNames    map[int]string
	Fields        map[fieldIndex]Field
	AncestorNames []string

	// Nested structs which correspond to any has-one has-many relationships
	// int is the ith element of this struct where the submap exists
	SubMaps map[fieldIndex]*Mapper
}

// Maps db rows onto the complex struct,
// Response must be a struct, pointer to a struct for our response, a slice of structs or slice of pointers to a struct
func Map(rows *sql.Rows, dst interface{}) error {
	var (
		mapper *Mapper
		err    error
		rsv    *resolver
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
	if !ok {
		if !(isSlicePtr(dstTyp) || isStructPtr(dstTyp)) {
			return fmt.Errorf("carta: cannot map rows onto %s, destination must be pointer to a slice(*[]) or pointer to struct", dstTyp)
		}
		// return errors.New("carta: destination pointer is nill, instantiate a new empty instance of destination before mapping")
		// }

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

		mapperCache.storeMap(columns, dstTyp, mapper)

	}

	if rsv, err = mapper.loadRows(rows, len(columns)); err != nil {
		return err
	}

	return setDst(mapper, reflect.ValueOf(dst), rsv)

}

func newMapper(t reflect.Type) (*Mapper, error) {
	var (
		crd     Cardinality
		elemTyp reflect.Type
		mapper  *Mapper
		subMaps map[fieldIndex]*Mapper
		err     error
	)

	isListPtr := false
	isBasic := false
	isTypePtr := false

	if isSlicePtr(t) {
		crd = Collection
		elemTyp = t.Elem().Elem() // *[]interface{} to intetrface{}
		isListPtr = true
	} else if t.Kind() == reflect.Slice {
		crd = Association
		crd = Collection
		elemTyp = t.Elem() // []interface{} to intetrface{}
	}

	if crd == Collection {
		isBasic = isBasicType(t)
		if elemTyp.Kind() == reflect.Ptr {
			elemTyp = elemTyp.Elem()
			isTypePtr = true
		}
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
	if subMaps, err = findSubMaps(mapper.Typ); err != nil {
		return nil, err
	}
	mapper.SubMaps = subMaps
	return mapper, nil
}

func findSubMaps(t reflect.Type) (map[fieldIndex]*Mapper, error) {
	var (
		subMap *Mapper
		err    error
	)
	subMaps := map[fieldIndex]*Mapper{}
	if t.Kind() != reflect.Struct {
		return nil, nil
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if isExported(field) && isSubMap(field.Type) {
			if subMap, err = newMapper(field.Type); err != nil {
				return nil, err
			}
			subMaps[fieldIndex(i)] = subMap
		}
	}
	return subMaps, nil
}

func determineFieldsNames(m *Mapper, ancestorNames []string) error {
	var (
		name string
	)
	fields := map[fieldIndex]Field{}

	m.AncestorNames = ancestorNames

	if m.IsBasic {
		return nil
	}
	if m.Typ.Kind() != reflect.Struct {
		log.Fatal(m.Typ)
	}

	for i := 0; i < m.Typ.NumField(); i++ {
		field := m.Typ.Field(i)
		if isExported(field) {
			if tag := nameFromTag(field.Tag); tag != "" {
				name = tag
			} else {
				name = field.Name
			}
			fields[fieldIndex(i)] = Field{
				Name: name,
				Typ:  field.Type,
				Kind: field.Type.Kind(),
			}
		}
	}
	m.Fields = fields
	if ancestorNames == nil {
		ancestorNames = []string{}
	}
	for i, subMap := range m.SubMaps {
		if err := determineFieldsNames(subMap, append(ancestorNames, fields[i].Name)); err != nil {
			return err
		}
	}
	return nil
}

// grow is pretty much an append function, input is the slice and a pointer to the new element
// dst must always be a pointer, if it is a pointer to slice
func (m *Mapper) grow(dst reflect.Value, newElemAddr reflect.Value) {
	dstIndirect := reflect.Indirect(dst)
	if m.IsTypePtr {
		dstIndirect.Set(reflect.Append(dstIndirect, newElemAddr))
	} else {
		dstIndirect.Set(reflect.Append(dstIndirect, reflect.Indirect(newElemAddr)))
	}
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
	return (!isBasicType(t) && (t.Kind() == reflect.Struct || t.Kind() == reflect.Slice))
}

// Basic types are any types that are intended to be set from sql row data
// Primative fields, sql.NullXXX, time.Time, pg timestamp qualify as basic
func isBasicType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if _, ok := value.BasicKinds[t.Kind()]; ok {
		return true
	}
	if _, ok := value.BasicTypes[t]; ok {
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
