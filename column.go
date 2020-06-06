package carta

import (
	"database/sql"
	"log"
	"sort"
	"strings"
)

var _ = log.Fatal

// column represents the ith struct field of this mapper where the column is to be mapped
type column struct {
	typ         *sql.ColumnType
	columnIndex int
	i           fieldIndex
}

func allocateColumns(m *Mapper, columns map[string]column) error {
	var (
		candidates map[string]bool
	)
	presentColumns := map[string]column{}
	for cName, c := range columns {
		for i, field := range m.Fields {
			candidates = getColumnNameCandidates(field.Name, m.AncestorNames)
			if _, ok := candidates[cName]; ok {
				c.i = i
				presentColumns[cName] = c
				delete(columns, cName) // dealocate claimed column
			}
		}
	}
	m.PresentColumns = presentColumns

	columnIds := []int{}
	for _, columnField := range m.PresentColumns {
		columnIds = append(columnIds, columnField.columnIndex)
	}
	sort.Ints(columnIds)
	m.SortedColumnIndexes = columnIds

	ancestorColumns := map[string]column{}
	for columnName, column := range m.AncestorColumns {
		ancestorColumns[columnName] = column
	}
	for columnName, column := range m.PresentColumns {
		ancestorColumns[columnName] = column
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
	candidates := map[string]bool{fieldName: true, ToSnakeCase(fieldName): true}
	if ancestorNames == nil {
		return candidates
	}
	nameConcat := fieldName
	for i := len(ancestorNames) - 1; i >= 0; i-- {
		nameConcat = ancestorNames[i] + "_" + nameConcat
		candidates[nameConcat] = true
		candidates[strings.ToLower(nameConcat)] = true
	}
	// log.Println(candidates)
	return candidates
}

func getColumnGoType(cTyp *sql.ColumnType) {
}

func ToSnakeCase(s string) string {
	delimiter := "_"
	s = strings.Trim(s, " ")
	n := ""
	for i, v := range s {
		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		nextCaseIsChanged := false
		if i+1 < len(s) {
			next := s[i+1]
			vIsCap := v >= 'A' && v <= 'Z'
			vIsLow := v >= 'a' && v <= 'z'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			if (vIsCap && nextIsLow) || (vIsLow && nextIsCap) {
				nextCaseIsChanged = true
			}
		}

		if i > 0 && n[len(n)-1] != uint8(delimiter[0]) && nextCaseIsChanged {
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				n += string(delimiter) + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + string(delimiter)
			}
		} else if v == ' ' || v == '-' {
			n += string(delimiter)
		} else {
			n = n + string(v)
		}
	}
	return strings.ToLower(n)
}
