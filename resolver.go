package carta

import (
	"reflect"
)

// Resolver determines whether an object has already appeared in past rows.
// ie, if a set of column values was previously returned by SQL,
// this is nececaty to determine whether a new instantiation of a type is necesarry
// Carta uses all present columns in a particular message to generate a unique id,
// if successive rows have the same id, it identifies the same element
// always include a uniquely identifiable column in your query
// resolver cannot be stored in pointer reciver, this would result in concurrency bugs,
//
// for example, if user requests mapping to *[]*User  where
// type User  struct {
//     userId    int
//     Addresses []*Address
// }
// if sql query returns multiple rows with the same userId, resolver will return
// a pointer to the *User with that id so that furhter mapping can continue, in this case, mapping of address
//
// TODO: consider passing resover in context value

type (
	uniqueValId string
	fieldIndex  int
)

type element struct {
	v       reflect.Value // value of a struct that is mapped, this is never a pointer, its either a primative or struct
	subMaps map[fieldIndex]*resolver
}

type resolver struct {
	elements map[uniqueValId]*element
}

func newResolver() *resolver {
	return &resolver{
		elements: map[uniqueValId]*element{},
	}
}
