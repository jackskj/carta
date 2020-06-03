package carta

import (
	"reflect"
	"sync"
)

var mapperCache = newCache()

type cache struct {
	mapCache  sync.Map
	typeCache sync.Map
}

func newCache() *cache {
	return &cache{}
}

type MapperEntry struct {
	columns []string
	dst     reflect.Type
}

func (c *cache) loadMap(columns []string, dst reflect.Type) (mapper *Mapper, ok bool) {
	vmap, ok := c.mapCache.Load(MapperEntry{columns: columns, dst: dst})
	if ok {
		mapper = vmap.(*Mapper)
	}
	return
}

func (c *cache) storeMap(columns []string, dst reflect.Type, mapper *Mapper) {
	c.mapCache.Store(MapperEntry{columns: columns, dst: dst}, mapper)
}

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
type resolver struct {
	uniqueIds map[string]reflect.Value
}

func newResolver() *resolver {
	return &resolver{
		uniqueIds: map[string]reflect.Value{},
	}
}

func (r *resolver) Load(uid string) (cachedDst reflect.Value, ok bool) {
	cachedDst, ok = r.uniqueIds[uid]
	return
}

func (r *resolver) Store(uid string, dst reflect.Value) {
	r.uniqueIds[uid] = dst
}
