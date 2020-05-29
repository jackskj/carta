package carta

import (
	"fmt"
	"reflect"
)

type mapSetter interface {
	set(dst, v interface{}, i int, res *resolver, uid string) error
	grow(dst interface{}) (interface{}, error)
	subMapByIndex(dst interface{}, i int) interface{}
}

type loadFunc func(v, val interface{}) error

type Setter struct {
	typ     reflect.Type
	loaders map[int]loadFunc
}

// List setter is used for setting has-many relationships, aka Collections
// this is used when the destination is *[]Struct
type ListSetter struct {
	Setter
}

func newListSetter(typ reflect.Type) *ListSetter {
	return &ListSetter{}
}

func (ls *ListSetter) set(dst, v interface{}, i int, res *resolver, uid string) error {
	if err := ls.loaders[i](dst, v); err != nil {
		return err
	}
	return nil
}

func (ls *ListSetter) grow(dst interface{}) (interface{}, error) {
	return nil, nil
}
func (ls *ListSetter) subMapByIndex(dst interface{}, i int) interface{} {
	return nil
}

// Struct setter is used for setting has-one relationships, aka Associations
// this is used when the destination is *Struct
type StructSetter struct {
	Setter
}

func newStructSetter(typ reflect.Type) *StructSetter {
	return &StructSetter{}
}

func (ss *StructSetter) set(dst, v interface{}, fieldIndex int, res *resolver, uid string) error {
	if err := ss.loaders[fieldIndex](dst, v); err != nil {
		return err
	}
	return nil
}

type NonSlinceGrowError struct {
	typ reflect.Type
}

func (f NonSlinceGrowError) Error() string {
	return fmt.Sprintf("carta: Multiple mappings were found for non-slice type, %T. Sql data is therefore ommitted. Consider debugging your query, verifying your sql relational integrity, or chading the type to a slice.", f.typ)
}

// Grow is really tricky for structs, invaction of this funvtion more than indicated either user's mistake, or broker referencial integrity
// To exaplain, if a user asks to map sql response to *User, he/she expects only one user resonse
// However, invocattion of this function more thank once indicates that carta logic determined that the sql response actually would map to many Users,
func (ss *StructSetter) grow(dst interface{}) (interface{}, error) {
	// s
	return nil, NonSlinceGrowError{ss.typ}
}

func (ss *StructSetter) subMapByIndex(dst interface{}, i int) interface{} {
	return nil
}
