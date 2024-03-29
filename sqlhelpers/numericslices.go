package sqlhelpers

import (
	"bytes"
	"errors"
	"log"
	"strconv"
)

type IntSlice []int

func (slice *IntSlice) ToString() string {
	var buffer bytes.Buffer
	err := buffer.WriteByte('{')
	if err != nil {
		log.Fatal(err)
	}

	for index, integer := range *slice {
		if index > 0 {
			err = buffer.WriteByte(',')
			if err != nil {
				log.Fatal(err)
			}
		}

		_, err = buffer.WriteString(strconv.Itoa(integer))
		if err != nil {
			log.Fatal(err)
		}
	}

	err = buffer.WriteByte('}')
	if err != nil {
		log.Fatal(err)
	}

	return buffer.String()
}

func (slice *IntSlice) Scan(src interface{}) error {
	var newSlice []int
	srcByteArray, castSuccessful := src.([]byte)
	if castSuccessful {
		trimmedByteSlice := bytes.Trim(srcByteArray, "{}")
		if len(trimmedByteSlice) == 0 {
			*slice = []int{}
			return nil
		}

		numericByteSlices := bytes.Split(trimmedByteSlice, []byte(","))
		for _, numericByteSlice := range numericByteSlices {
			numberBytes := bytes.Trim(numericByteSlice, "\" ")
			numberString := string(numberBytes[:])
			number, err := strconv.ParseInt(numberString, 10, 32)
			if err != nil {
				return err
			}

			newSlice = append(newSlice, int(number))
		}

		*slice = newSlice
		return nil
	}

	return errors.New("Value passed into Scan could not be converted to type []byte")
}
