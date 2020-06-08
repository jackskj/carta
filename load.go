package carta

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/jackskj/carta/value"
)

func (m *Mapper) loadRows(rows *sql.Rows, columnNum int) (*resolver, error) {
	defer rows.Close() // may not need
	var err error
	row := make([]interface{}, columnNum)
	rsv := newResolver()
	for rows.Next() {
		for i, _ := range row {
			row[i] = new(interface{})
		}
		if err = rows.Scan(row...); err != nil {
			return nil, err
		}
		if err = loadRow(m, row, rsv); err != nil {
			return nil, err
		}
	}

	// e, err := json.Marshal(rsv.uniqueVals)
	// if err != nil {
	// log.Fatalf("as" + err.Error())
	// }
	// log.Println("a" + string(e))
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
func loadRow(m *Mapper, row []interface{}, rsv *resolver) error {
	var (
		err      error
		dstField reflect.Value // destination field to be set with
		cell     value.Cell
		elem     *element
		found    bool
	)

	uid := getUniqueId(row, m)

	if elem, found = rsv.elements[uid]; !found {
		// unique row mapping found, new object
		loadElem := reflect.New(m.Typ).Elem()

		for _, field := range m.PresentColumns {
			if m.IsBasic {
				dstField = loadElem
			} else {
				dstField = loadElem.Field(int(field.i))
			}
			// sql.Row.Scan() retuens pointers in each cell, I have to use pointer indirection here
			srcI := *row[field.columnIndex].(*interface{})

			if srcI == nil { // returned sql cell is nil
				cell = value.NewNull()
			} else {
				if cell, err = value.NewCell(srcI, field.typ); err != nil {
					return err
				}
			}

			//setting sql value onto the field
			if !m.IsBasic {
				var (
					kind reflect.Kind
					dst  reflect.Value
					typ  reflect.Type
				)
				if m.Fields[field.i].IsPtr && cell.IsValid() {
					dst = reflect.New(m.Fields[field.i].ElemTyp).Elem()
					kind = m.Fields[field.i].ElemKind
					typ = m.Fields[field.i].ElemTyp
				} else {
					kind = m.Fields[field.i].Kind
					typ = m.Fields[field.i].Typ
					dst = dstField
				}
				switch kind {
				case reflect.Bool:
					dst.SetBool(cell.Bool())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					dst.SetUint(cell.Uint64())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					dst.SetInt(cell.Int64())
				case reflect.String:
					dst.SetString(cell.String())
				case reflect.Float32, reflect.Float64:
					dst.SetFloat(cell.Float64())
				case reflect.Struct:
					if strTyp, ok := value.BasicTypes[typ]; ok {
						// TODO: Type asserion, prevent from calling ValueOf
						switch strTyp {
						case value.Time:
							dst.Set(reflect.ValueOf(cell.Time()))
						case value.Timestamp:
							dst.Set(reflect.ValueOf(cell.Timestamp()))
						case value.NullBool:
							dst.Set(reflect.ValueOf(cell.NullBool()))
						case value.NullFloat64:
							dst.Set(reflect.ValueOf(cell.NullFloat64()))
						case value.NullInt32:
							dst.Set(reflect.ValueOf(cell.NullInt32()))
						case value.NullInt64:
							dst.Set(reflect.ValueOf(cell.NullInt64()))
						case value.NullString:
							dst.Set(reflect.ValueOf(cell.NullString()))
						case value.NullTime:
							dst.Set(reflect.ValueOf(cell.NullTime()))
						}
					}
				}
				if m.Fields[field.i].IsPtr && cell.IsValid() {
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
func getUniqueId(row []interface{}, m *Mapper) (uid uniqueValId) {
	uid = ""
	for _, i := range m.SortedColumnIndexes {
		// TODO: Implement a more advanced and better performing hashing
		// this can be done by dumping row values into []byte by using
		// sql.driver.valuer interface
		r := row[i].(*interface{})
		uid = uid + uniqueValId(fmt.Sprintf("%v|", *r))
	}
	return
}
