package models

import (
	"fmt"
)

// StringFilter is a type used to filter with strings in an SQL Query
// li is for LIKE
// eq is for EQUAL
type StringFilter struct {
	li []string
	eq []string
}

// IntFilter is a type used to filter numbers in an SQL Query
// lt is for LESS THAN
// gt is for GREATER THAN
// eq if for EQUALS
type IntFilter struct {
	lt int64
	gt int64
	eq int64
}

func buildStringFilterQuery(filter StringFilter, colName string) (string, bool) {
	stmt := "("
	if len(filter.eq) != 0 {
		for ndx, val := range filter.eq {
			stmt += fmt.Sprintf("%s = \"%s\"", colName, val)
			if ndx < len(filter.eq)-1 {
				stmt += " OR "
			}
		}
		stmt += ")"
	} else if len(filter.li) != 0 {
		for ndx, val := range filter.li {
			stmt += fmt.Sprintf("%s LIKE \"%%%s%%\"", colName, val)
			if ndx < len(filter.li)-1 {
				stmt += " OR "
			}
		}
		stmt += ")"
	}

	if stmt == "(" {
		return "", false
	}
	return stmt, true
}

func buildFloat64FilterQuery(filter IntFilter, colName string) (string, bool) {
	stmt := "("
	if filter.eq != -1 {
		stmt += fmt.Sprintf("%s = %d", colName, filter.eq)
		stmt += ")"
	} else {
		and := false
		if filter.gt != -1 {
			stmt += fmt.Sprintf("%s > %d", colName, filter.gt)
			and = true
		}
		if filter.lt != -1 {
			if and {
				stmt += " AND "
			}
			stmt += fmt.Sprintf("%s < %d", colName, filter.lt)
		}
		stmt += ")"
	}
	if stmt == "(" {
		return "", false
	}
	return stmt, true
}
