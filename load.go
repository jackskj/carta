package carta

import (
	"database/sql"
	"fmt"
	"log"
)

func (m *Mapper) loadRows(rows *sql.Rows, dst interface{}) error {
	defer rows.Close() // may not need
	var err error
	row := make([]interface{}, len(m.columns))
	rsv := newResolver()
	for rows.Next() {
		// for i, _ := range row {
		// row[i] =interface{}{}
		// }
		if err = rows.Scan(row...); err != nil {
			return err
		}
		if err = loadRow(m.SqlMap, row, dst, rsv); err != nil {
			return err
		}
	}
	return nil
}

func loadRow(m *Mapper, row []interface{}, resp interface{}, rsv *resolver) error {
	var err error
	// todo: explain difference between resp and dst
	var dst interface{}

	uid := getUniqueId(row, m)
	if cachedDst, ok := rsv.Load(uid); ok {
		dst = cachedDst
	} else {
		// uniqur row mapping found
		// example if resp is *[]*User, dst is *User
		// after this block, grow will append new *User and return that *User as dst
		resp, dst, err = m.grow(resp)
		if growErr, ok := err.(NonSlinceGrowError); ok {
			log.Println(growErr.Error()) // not breaking, but sql and/or expected response is incorrect
			// todo: consider this breaking, and simply return err
			return nil
		} else if err != nil {
			return err
		}

		for _, field := range m.PresentColumns {
			if err := m.set(dst, row[field.columnIndex], field.fieldIndex); err != nil {
				return err
			}
		}
		rsv.Store(uid, dst)
	}

	for i, subMap := range m.SubMaps {
		subMapDst := m.subMapByIndex(resp, i)
		if err := loadRow(subMap, row, subMapDst, rsv); err != nil {
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
// userId int
// Addresses []*Address
// }
// if sql query returns multiple rows with the same userId, resolver will return
// a pointer to the *User with that id so that furhter mapping can continue, in this case, mapping of address
//
// TODO: consider passing resover in context value
type resolver struct {
	uniqueIds map[string]interface{}
}

func newResolver() *resolver {
	return &resolver{
		uniqueIds: map[string]interface{}{},
	}
}

func (r *resolver) Load(uid string) (cachedDst interface{}, ok bool) {
	cachedDst, ok = r.uniqueIds[uid]
	return
}

func (r *resolver) Store(uid string, dst interface{}) {
	r.uniqueIds[uid] = dst
}
