package carta

import (
	"reflect"
	"strings"
	"sync"
)

var mapperCache = newCache()

type cache struct {
	mapCache sync.Map
}

func newCache() *cache {
	return &cache{}
}

type mapperEntry struct {
	columns []string
	dst     reflect.Type
}

func (m *mapperEntry) raw() string {
	// TODO: test how this works with unexported types
	// TODO: add a way to provide fully qualified name for the type, since m.typ is always a pointer to a struct or slice
	// return strings.Join(m.columns, ",") + "|" + m.dst.PkgPath() + "." + m.dst.String()
	return strings.Join(m.columns, ",") + "|" + m.dst.String()
}

func (c *cache) loadMap(columns []string, dst reflect.Type) (mapper *Mapper, ok bool) {
	entry := mapperEntry{columns, dst}
	vmap, ok := c.mapCache.Load(entry.raw())
	if ok {
		mapper = vmap.(*Mapper)
	}
	return
}

func (c *cache) storeMap(columns []string, dst reflect.Type, mapper *Mapper) {
	entry := mapperEntry{columns, dst}
	c.mapCache.Store(entry.raw(), mapper)
}
