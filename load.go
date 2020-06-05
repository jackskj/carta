package carta

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"github.com/jackskj/carta/value"
)

var _ = log.Fatal

func (m *Mapper) loadRows(rows *sql.Rows, dst interface{}) error {
	defer rows.Close() // may not need
	var err error
	row := make([]interface{}, len(m.PresentColumns))
	rsv := newResolver()
	dstV := reflect.ValueOf(dst)
	for rows.Next() {
		for i, _ := range row {
			row[i] = new(interface{})
		}
		if err = rows.Scan(row...); err != nil {
			return err
		}
		if err = loadRow(m, row, dstV, rsv); err != nil {
			return err
		}
	}
	return nil
}

func loadRow(m *Mapper, row []interface{}, dst reflect.Value, rsv *resolver) error {
	var (
		err      error
		newElem  reflect.Value // destination to be set with dstV
		dstField reflect.Value // destination to be set with dstV
		cell     value.Cell
	)
	uid := getUniqueId(row, m)
	if cachedDst, ok := rsv.Load(uid); ok {
		// mapping of this object has already happened
		// this is is either due to has-many relationship or duplicate row
		dst = cachedDst
	} else {
		// unique row mapping found, new object
		newElem = reflect.New(m.Typ).Elem()

		for _, field := range m.PresentColumns {
			if m.IsBasic {
				// if v, err = m.BasicConverter(srcV); err != nil {
				// return err
				// }
				dstField = newElem
			} else {
				// if v, err = m.Converters[field.fieldIndex](srcV); err != nil {
				// return err
				// }
				dstField = newElem.Field(field.fieldIndex)
			}
			// sql.Row.Scan() requires pointers for each cell
			srcI := row[field.columnIndex].(*interface{})
			if srcI == nil { // sql cell is nil
				cell = value.NewNull()
			} else {
				if cell, err = value.NewCell(*srcI, field.typ); err != nil {
					return err
				}
			}
			dstField.Set(reflect.ValueOf(cell.AsInterface()))
		}
		// grow() will return append a new zero value to the slice, and return that value to be set
		dst = m.grow(dst, newElem)

	}
	log.Println(dst)
	rsv.Store(uid, dst)
	for i, subMap := range m.SubMaps {
		if err = loadRow(subMap, row, dst.Field(i), rsv); err != nil {
			return err
		}
	}
	return nil
}

// Generates unique id based on the ancestors of the struct as well as currently considered colum values
func getUniqueId(row []interface{}, m *Mapper) string {

	columnIds := []int{}
	for _, columnField := range m.PresentColumns {
		columnIds = append(columnIds, columnField.columnIndex)
	}
	for _, column := range m.AncestorColumns {
		columnIds = append(columnIds, column.columnIndex)
	}

	out := ""
	for _, i := range columnIds {
		// TODO: Implement a more advanced and better performing hashing
		// this can be done by dumping row values into []byte by using
		// sql.driver.valuer interface
		r := row[i].(*interface{})
		out = out + fmt.Sprintf("%v|", *r)
	}

	return out
}
