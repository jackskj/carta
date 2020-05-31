package carta

import (
	"database/sql"
	"strings"
)

// column represents the ith struct field of this mapper where the column is to be mapped
type column struct {
	typ         *sql.ColumnType
	columnIndex int
	fieldIndex  int
}

func allocateColumns(m *Mapper, columns map[string]column) error {
	var (
		candidates map[string]bool
	)
	presentColumns := map[string]column{}
	for cName, c := range columns {
		for i, fieldName := range m.FieldNames {
			candidates = getColumnNameCandidates(fieldName, m.AncestorNames)
			if _, ok := candidates[cName]; ok {
				c.fieldIndex = i
				presentColumns[cName] = c
				delete(columns, cName) // dealocate claimed column
			}
		}
	}
	m.PresentColumns = presentColumns
	ancestorColumns := map[string]bool{}
	for columnName, _ := range m.AncestorColumns {
		ancestorColumns[columnName] = true
	}
	for columnName, _ := range m.PresentColumns {
		ancestorColumns[columnName] = true
	}
	for _, subMap := range m.SubMaps {
		subMap.AncestorColumns = ancestorColumns
		if err := allocateColumns(subMap, columns); err != nil {
			return err
		}
	}
	return nil
}

func getColumnNameCandidates(fieldName string, ancestorNames []string) map[string]bool {
	candidates := map[string]bool{fieldName: true, strings.ToLower(fieldName): true}
	if ancestorNames == nil {
		return candidates
	}
	nameConcat := fieldName
	for i := len(ancestorNames) - 1; i >= 0; i-- {
		nameConcat = ancestorNames[i] + "_" + nameConcat
		candidates[nameConcat] = true
		candidates[strings.ToLower(nameConcat)] = true
	}
	return candidates
}
