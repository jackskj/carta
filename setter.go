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
	for _, elem := range rsv.elements {
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
				size := len(subMapRsv.elements)
				if m.IsTypePtr {
					newChildElem = reflect.New(reflect.SliceOf(reflect.PtrTo(childTyp))).Elem()
					newChildElem.Set(reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(childTyp)), size, size))
				} else {
					newChildElem = reflect.New(reflect.SliceOf(childTyp)).Elem()
					newChildElem.Set(reflect.MakeSlice(reflect.SliceOf(childTyp), size, size))
				}
				if m.IsListPtr {
					elem.v.Field(int(fieldIndex)).Set(newChildElem)
					childDst = elem.v.Field(int(fieldIndex))
				} else {
					elem.v.Field(int(fieldIndex)).Set(reflect.Indirect(newChildElem))
					childDst = elem.v.Field(int(fieldIndex)).Addr()
				}
			} else if subMap.Crd == Association {
				newChildElem = reflect.New(childTyp).Elem()
				if m.IsTypePtr {
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

	for _, elem := range rsv.elements {
		if m.Crd == Collection {
			if m.IsTypePtr {
				dstIndirect.Set(reflect.Append(dstIndirect, elem.v.Addr()))
			} else {
				dstIndirect.Set(reflect.Append(dstIndirect, elem.v))
			}
		} else if m.Crd == Association {

			// loadElem := reflect.New(m.Typ).Elem()

			if m.IsTypePtr {
				reflect.Indirect(dst).Set(elem.v)
				log.Printf("%v %v %v", dst.Type(), elem.v.Type(), dst)
			} else {
				reflect.Indirect(dst).Set(elem.v)
			}
		}
	}

	return nil
}
