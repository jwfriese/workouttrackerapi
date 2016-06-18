package sqlhelpers

import (
	"errors"
	"strings"
)

type SQLStringSlice struct {
	internalStringSlice []string
}

func (s *SQLStringSlice) ToStringSlice() []string {
	return s.internalStringSlice
}

func (s *SQLStringSlice) Scan(src interface{}) error {
	srcByteArray, castSuccessful := src.([]byte)
	srcString := string(srcByteArray)
	if castSuccessful {
		s.internalStringSlice = []string{}
		stringBracketsTrimmed := strings.Trim(srcString, "{}")
		collectedStrings := strings.Split(stringBracketsTrimmed, ",")
		for _, collectedString := range collectedStrings {
			s.internalStringSlice = append(s.internalStringSlice, strings.Trim(collectedString, "\""))
		}

		return nil
	}

	return errors.New("Value passed into Scan could not be converted to type string")
}
