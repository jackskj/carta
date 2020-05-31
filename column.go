package carta

import (
	"database/sql"
)

type column struct {
	typ *sql.ColumnType
	i   int
}

func allocateColumns(m *Mapper, columns map[string]column) error {

	return nil
}
