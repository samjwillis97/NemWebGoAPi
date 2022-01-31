package models

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
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

func (f *StringFilter) fromFilterMap(filterMap map[string][]string, param string) {
	if val, ok := filterMap[param+".li"]; ok {
		f.li = val
	}
	if val, ok := filterMap[param+".eq"]; ok {
		f.eq = val
	}
}

func (f *IntFilter) fromFilterMap(filterMap map[string][]string, param string) {
	if val, ok := filterMap[param+".lt"]; ok {
		f.lt, _ = strconv.ParseInt(val[0], 10, 64)
	} else {
		f.lt = -1
	}
	if val, ok := filterMap[param+".gt"]; ok {
		f.gt, _ = strconv.ParseInt(val[0], 10, 64)
	} else {
		f.gt = -1
	}
	if val, ok := filterMap[param+".eq"]; ok {
		f.eq, _ = strconv.ParseInt(val[0], 10, 64)
	} else {
		f.eq = -1
	}
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

func buildInt64FilterQuery(filter IntFilter, colName string) (string, bool) {
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

// buildFilterQuery takes a struct that consists of filters like StringFilter and IntFilter
// It then generates the SQL query to apply these filters
// requires the struct fields have a tag of "col" to be able to generate query effectively
func buildFilterQuery(filter interface{}) string {
	var stmt string
	filterArr := make([]string, 0)

	v := reflect.Indirect(reflect.ValueOf(filter))
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		col := t.Field(i).Tag.Get("col")

		if col == "" {
			continue
		}

		switch reflect.TypeOf(filter).Field(i).Type {
		case reflect.TypeOf(StringFilter{}):
			concreteVal, _ := fieldVal.Interface().(StringFilter)
			if val, ok := buildStringFilterQuery(concreteVal, col); ok {
				filterArr = append(filterArr, val)
			}
		case reflect.TypeOf(IntFilter{}):
			concreteVal, _ := fieldVal.Interface().(IntFilter)
			if val, ok := buildInt64FilterQuery(concreteVal, col); ok {
				filterArr = append(filterArr, val)
			}
		}
	}

	if len(filterArr) == 0 {
		return ""
	}

	for ndx, val := range filterArr {
		if ndx == 0 {
			stmt += "\nWHERE "
		} else {
			stmt += "\nAND "
		}
		stmt += val
	}
	return stmt
}

// ParseFilterMap converts the query parameters returned by net/http into a filter as defined in the destination
// TODO: Finish - may not be possible easily
func ParseFilterMap(filterMap map[string][]string, dest interface{}) {
	v := reflect.Indirect(reflect.ValueOf(dest))
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		param := t.Field(i).Tag.Get("param")

		if param == "" {
			continue
		}

		switch t.Field(i).Type {
		case reflect.TypeOf(StringFilter{}):
			elementName := param + ".li"
			if val, ok := filterMap[elementName]; ok {
				fmt.Println("Found: ", elementName)
				fmt.Println(val)
			}

		case reflect.TypeOf(IntFilter{}):

		}
	}

}

func getFloatReflectOnly(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(reflect.TypeOf(float64(0))) {
		return math.NaN(), fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(reflect.TypeOf(float64(0)))
	return fv.Float(), nil
}
