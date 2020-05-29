package carta

import (
	db "database/sql"
	"fmt"
	"log"
)

func (m *Mapper) loadRows(rows *db.Rows, dst interface{}) error {
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

func loadRow(sqlMap *SqlMap, row []interface{}, resp interface{}, rsv *resolver) error {
	var err error
	// todo: explain difference between resp and dst
	var dst interface{}

	uid := getUniqueId(row, sqlMap)
	if cachedDst, ok := rsv.Load(uid); ok {
		dst = cachedDst
	} else {
		// example if dst is *[]*User
		// after this block, grow will appens new(*User) and return that *User
		dst, err = sqlMap.mapSetter.grow(resp)
		if growErr, ok := err.(NonSlinceGrowError); ok {
			log.Println(growErr.Error()) // not breaking, but sql and/or expected response if incorrect
			return nil
		} else if err != nil {
			return err
		}
	}

	for _, field := range sqlMap.PresentColumns {
		if err := sqlMap.mapSetter.set(dst, row[field.columnIndex], field.fieldIndex, rsv, uid); err != nil {
			return err
		}
	}
	for i, subMap := range sqlMap.SubMaps {
		subMapDst := sqlMap.mapSetter.subMapByIndex(dst, i)
		if err := loadRow(subMap, subMapDst, subMapDst, rsv); err != nil {
			return err
		}
	}
	return nil
}

// Generates unique id based on the ancestors of the struct as well as currently considered colum values
func getUniqueId(row []interface{}, sqlMap *SqlMap) string {

	columnIds := []int{}
	for _, columnField := range sqlMap.PresentColumns {
		columnIds = append(columnIds, columnField.columnIndex)
	}
	for _, columnField := range sqlMap.AncestorColumns {
		columnIds = append(columnIds, columnField.columnIndex)
	}

	out := ""
	for _, i := range columnIds {
		// TODO: Implement a more advances, and better performing hashing
		out = out + fmt.Sprintf("%v|", row[i])
	}

	return out
}

// Resolver determines whether a particular set of values were already returned from sql.
// if a set of column values was previously returned by SQL,
// this is nececaty to determine whether a new instantiation of a type is necesarry
// Carta uses all present columns in a particular message to generate a unique id,
// if successive rows have the same id, it identifies the same element
// always include a uniquely identifiable column in your query
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
