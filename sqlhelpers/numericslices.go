package sqlhelpers

import (
	"bytes"
	"errors"
	"strconv"
)

type UIntSlice []uint

func (slice *UIntSlice) Scan(src interface{}) error {
	var newSlice []uint
	srcByteArray, castSuccessful := src.([]byte)
	if castSuccessful {
		trimmedByteSlice := bytes.Trim(srcByteArray, "{}")
		numericByteSlices := bytes.Split(trimmedByteSlice, []byte(","))
		for _, numericByteSlice := range numericByteSlices {
			numberBytes := bytes.Trim(numericByteSlice, "\" ")
			numberString := string(numberBytes[:])
			number, err := strconv.ParseUint(numberString, 10, 32)
			if err != nil {
				return err
			}

			newSlice = append(newSlice, uint(number))
		}

		*slice = newSlice
		return nil
	}

	return errors.New("Value passed into Scan could not be converted to type []byte")
}
