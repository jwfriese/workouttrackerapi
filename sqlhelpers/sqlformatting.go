package sqlhelpers

import (
	"strconv"
)

func Float32PointerToSQLString(ptr *float32) string {
	sqlString := "NULL"
	if ptr != nil {
		sqlString = strconv.FormatFloat(float64(*ptr), 'f', -1, 32)
	}
	return sqlString
}

func IntPointerToSQLString(ptr *int) string {
	sqlString := "NULL"
	if ptr != nil {
		sqlString = strconv.Itoa(*ptr)
	}
	return sqlString
}
