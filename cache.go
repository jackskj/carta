package carta

import (
	"reflect"
	"sync"
)

var mapperCache = newCache()

type cache struct {
	typeCache sync.Map
	mapCache  sync.Map
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
