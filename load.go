package carta

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/jackskj/carta/value"
)

func (m *Mapper) loadRows(rows *sql.Rows, colTyps []*sql.ColumnType) (*resolver, error) {
	defer rows.Close() // may not need
	var err error
	row := make([]interface{}, len(colTyps))
	colTypNames := make([]string, len(colTyps))
	for i := 0; i < len(colTyps); i++ {
		colTypNames[i] = colTyps[i].DatabaseTypeName()
	}
	rsv := newResolver()
	for rows.Next() {
		for i := 0; i < len(colTyps); i++ {
			row[i] = value.NewCell(colTypNames[i])
		}
		if err = rows.Scan(row...); err != nil {
			return nil, err
		}
		if err = loadRow(m, row, rsv); err != nil {
			return nil, err
		}
	}
	return rsv, nil
}

// load row maps a single sql row onto a structure that resembles the users struct
// that mapping is stored in the resolver as a pointer reference to an instance of the struct
//
// if new object is foind, create a new instance of a struct that
// maps onto that struct,
// for example, if a user maps onto:
// type Blog struct {
//          BlogId string
// }
// blogs := []Blog:
// carta.Map(rows, &blogs)
// if a new blog_id column value is found, I instantiatiate a new instance of Blog,
// set BlogId, then store the pointer referenct to this instance in the resolver
// nothins is done when the object has been already mapped in previous rows, however,
// the function contunous to recursivelly map rows for all sub mappings inside Blog
//  for example, if a blog has many Authors
// rows are actually []*Cell, theu are passed here as interface since sql scan requires []interface{}
func loadRow(m *Mapper, row []interface{}, rsv *resolver) error {
	var (
		err      error
		dstField reflect.Value // destination field to be set with
		cell     *value.Cell
		elem     *element
		found    bool
	)

	uid := getUniqueId(row, m)

	if elem, found = rsv.elements[uid]; !found {
		// unique row mapping found, new object
		loadElem := reflect.New(m.Typ).Elem()

		for _, col := range m.PresentColumns {
			var (
				kind     reflect.Kind  // kind of destination
				dst      reflect.Value // destination to set
				typ      reflect.Type  // underlying type of the destination
				isDstPtr bool          //is the destination a pointer
			)

			cell = row[col.columnIndex].(*value.Cell)

			if m.IsBasic {
				dst = loadElem
				kind = m.Kind
				typ = m.Typ
				isDstPtr = m.IsTypePtr
			} else {
				dstField = loadElem.Field(int(col.i))
				if m.Fields[col.i].IsPtr {
					dst = reflect.New(m.Fields[col.i].ElemTyp).Elem()
					kind = m.Fields[col.i].ElemKind
					typ = m.Fields[col.i].ElemTyp
					isDstPtr = true
				} else {
					dst = dstField
					kind = m.Fields[col.i].Kind
					typ = m.Fields[col.i].Typ
					isDstPtr = false
				}
			}
			if cell.IsNull() {
				_, nullable := value.NullableTypes[typ]
				if !(isDstPtr || nullable) {
					return errors.New(fmt.Sprintf("carta: cannot load null value to type %s for column %s", typ, col.name))
				}
				// no need to set destination if cell is null
			} else {
				switch kind {
				case reflect.Bool:
					if d, err := cell.Bool(); err != nil {
						return value.ConvertsionError(err, typ)
					} else {
						dst.SetBool(d)
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if d, err := cell.Uint64(); err != nil {
						return value.ConvertsionError(err, typ)
					} else {
						dst.SetUint(d)
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if d, err := cell.Int64(); err != nil {
						return value.ConvertsionError(err, typ)
					} else {
						dst.SetInt(d)
					}
				case reflect.String:
					if d, err := cell.String(); err != nil {
						return value.ConvertsionError(err, typ)
					} else {
						dst.SetString(d)
					}
				case reflect.Float32, reflect.Float64:
					if d, err := cell.Float64(); err != nil {
						return value.ConvertsionError(err, typ)
					} else {
						dst.SetFloat(d)
					}
				case reflect.Struct:
					if strTyp, ok := value.BasicTypes[typ]; ok {
						// TODO: Type asserion, prevent from calling ValueOf
						// TODO: make these stupid error checks more concise
						//  this swich statement should be optimized

						switch strTyp {
						case value.Time:
							if d, err := cell.Time(); err != nil {
								return value.ConvertsionError(err, typ)
							} else {
								dst.Set(reflect.ValueOf(d))
							}
						case value.Timestamp:
							if d, err := cell.Timestamp(); err != nil {
								return value.ConvertsionError(err, typ)
							} else {
								dst.Set(reflect.ValueOf(d))
							}
						case value.NullBool:
							if d, err := cell.NullBool(); err != nil {
								return value.ConvertsionError(err, typ)
							} else {
								dst.Set(reflect.ValueOf(d))
							}
						case value.NullFloat64:
							if d, err := cell.NullFloat64(); err != nil {
								return value.ConvertsionError(err, typ)
							} else {
								dst.Set(reflect.ValueOf(d))
							}
						case value.NullInt32:
							if d, err := cell.NullInt32(); err != nil {
								return value.ConvertsionError(err, typ)
							} else {
								dst.Set(reflect.ValueOf(d))
							}
						case value.NullInt64:
							if d, err := cell.NullInt64(); err != nil {
								return value.ConvertsionError(err, typ)
							} else {
								dst.Set(reflect.ValueOf(d))
							}
						case value.NullString:
							if d, err := cell.NullString(); err != nil {
								return value.ConvertsionError(err, typ)
							} else {
								dst.Set(reflect.ValueOf(d))
							}
						case value.NullTime:
							if d, err := cell.NullTime(); err != nil {
								return value.ConvertsionError(err, typ)
							} else {
								dst.Set(reflect.ValueOf(d))
							}
						}
					}
				}
				if !m.IsBasic && m.Fields[col.i].IsPtr {
					dstField.Set(dst.Addr())
				}
			}
		}
		elem = &element{v: loadElem}
		if len(m.SubMaps) != 0 {
			elem.subMaps = map[fieldIndex]*resolver{}
			for i, _ := range m.SubMaps {
				elem.subMaps[i] = newResolver()
			}
		}
		rsv.elements[uid] = elem
		rsv.elementOrder = append(rsv.elementOrder, uid)
	}

	for i, subMap := range m.SubMaps {
		if err = loadRow(subMap, row, elem.subMaps[i]); err != nil {
			return err
		}
	}

	return nil
}

// Generates unique id based on the ancestors of the struct as well as currently considered colum values
func getUniqueId(row []interface{}, m *Mapper) uniqueValId {
	// TODO: set capacity of the uid slice, using bytes.buffer
	uid := ""
	for _, i := range m.SortedColumnIndexes {
		uid = uid + row[i].(*value.Cell).Uid()
	}
	return uniqueValId(uid)
}
