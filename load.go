package carta

import (
	"database/sql"
	"fmt"
	// "log"
	"reflect"
)

func (m *Mapper) loadRows(rows *sql.Rows, dst interface{}) error {
	defer rows.Close() // may not need
	var err error
	row := make([]interface{}, len(m.PresentColumns))
	rsv := newResolver()
	dstV := reflect.ValueOf(dst)
	for rows.Next() {
		// for i, _ := range row {
		// row[i] =interface{}{}
		// }
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
		srcV     reflect.Value // whatever value sql gives in one cell
		dstField reflect.Value // destination to be set with dstV
	)

	uid := getUniqueId(row, m)
	if cachedDst, ok := rsv.Load(uid); ok {
		// mapping of this object has already happened
		// this is is either due to has-many relationship or duplicate row
		dst = cachedDst
	} else {
		// uniqur row mapping found, new object
		// example if destination is []*User
		// grow() will return append a new zero value to the slice, and return that value to be set
		dst = m.grow(dst)
		for _, field := range m.PresentColumns {
			if m.IsBasic {
				// if v, err = m.BasicConverter(srcV); err != nil {
				// return err
				// }
				dstField = dst
			} else {
				// if v, err = m.Converters[field.fieldIndex](srcV); err != nil {
				// return err
				// }
				dstField = dst.Field(field.fieldIndex)
			}
			srcV = reflect.ValueOf(row[field.columnIndex])
			dstField.Set(srcV)
		}
	}
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
		out = out + fmt.Sprintf("%v|", row[i])
	}

	return out
}
