package carta

import (
	"errors"
	"log"
	"reflect"
)

var _ = log.Panic

func setDst(m *Mapper, dst reflect.Value, rsv *resolver) error {
	// dst is  always a pointer
	dstIndirect := reflect.Indirect(dst)

	// post order traversal, first set all submap structs, then the struct itself
	for _, uid := range rsv.elementOrder {
		elem := rsv.elements[uid]

		//set childeren first
		for fieldIndex, subMapRsv := range elem.subMaps {
			var (
				subMap       *Mapper
				childTyp     reflect.Type
				childDst     reflect.Value
				newChildElem reflect.Value
				ok           bool
			)

			if subMap, ok = m.SubMaps[fieldIndex]; !ok {
				// this should never happen
				return errors.New("carta: sub map not found")
			}
			if f, ok := m.SubMaps[fieldIndex]; ok {
				childTyp = f.Typ
			} else {
				// this should never happen
				return errors.New("carta: field not found")
			}

			if subMap.Crd == Collection {
				capacity := len(subMapRsv.elements)
				if subMap.IsTypePtr {
					newChildElem = reflect.New(reflect.SliceOf(reflect.PtrTo(childTyp))).Elem()
					newChildElem.Set(reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(childTyp)), 0, capacity))
				} else {
					newChildElem = reflect.New(reflect.SliceOf(childTyp)).Elem()
					newChildElem.Set(reflect.MakeSlice(reflect.SliceOf(childTyp), 0, capacity))
				}
				if subMap.IsListPtr {
					elem.v.Field(int(fieldIndex)).Set(newChildElem.Addr())
					childDst = elem.v.Field(int(fieldIndex))
				} else {
					elem.v.Field(int(fieldIndex)).Set(newChildElem)
					childDst = elem.v.Field(int(fieldIndex)).Addr()
				}
			} else if subMap.Crd == Association {
				newChildElem = reflect.New(childTyp).Elem()
				if subMap.IsTypePtr {
					elem.v.Field(int(fieldIndex)).Set(newChildElem.Addr())
					childDst = elem.v.Field(int(fieldIndex))
				} else {
					elem.v.Field(int(fieldIndex)).Set(newChildElem)
					childDst = elem.v.Field(int(fieldIndex)).Addr()
				}
			}

			// setting the child
			setDst(subMap, childDst, subMapRsv)
		}
	}

	for _, uid := range rsv.elementOrder {
		elem := rsv.elements[uid]
		if m.Crd == Collection {
			if m.IsTypePtr {
				dstIndirect.Set(reflect.Append(dstIndirect, elem.v.Addr()))
			} else {
				dstIndirect.Set(reflect.Append(dstIndirect, elem.v))
			}
		} else if m.Crd == Association {
			dstIndirect.Set(elem.v)
		}
	}
	return nil
}
