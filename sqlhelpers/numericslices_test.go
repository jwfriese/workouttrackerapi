package sqlhelpers_test

import (
	"log"

	"github.com/jwfriese/workouttrackerapi/sqlhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers for numeric types with SQL", func() {
	Describe("[]int", func() {
		Describe("Converting to a SQL-friendly string representation", func() {
			var (
				result string
			)

			Context("When there are digits in the slice", func() {
				BeforeEach(func() {
					input := sqlhelpers.IntSlice{1, 1, 3, 5, 8, 13}
					result = input.ToString()
				})

				It("creates the string list", func() {
					Expect(result).To(Equal("{1,1,3,5,8,13}"))
				})
			})

			Context("When the slice is empty", func() {
				BeforeEach(func() {
					input := sqlhelpers.IntSlice{}
					result = input.ToString()
				})

				It("creates the empty string list", func() {
					Expect(result).To(Equal("{}"))
				})
			})
		})

		Describe("Scanning from a SQL result", func() {
			var (
				err    error
				result sqlhelpers.IntSlice
			)

			Context("When given a list containing no numbers", func() {
				BeforeEach(func() {
					input := []byte("{}")
					err = (&result).Scan(input)

					if err != nil {
						log.Fatal(err)
					}
				})

				It("puts an empty int array in the result", func() {
					Expect(result).To(BeEquivalentTo([]int{}))
				})
			})

			Context("When given a list containing some numbers", func() {
				BeforeEach(func() {
					input := []byte("{1,20,5}")
					err = (&result).Scan(input)

					if err != nil {
						log.Fatal(err)
					}
				})

				It("pulls out all of the ints values from the byte slice", func() {
					Expect(result).To(BeEquivalentTo([]int{1, 20, 5}))
				})
			})
		})
	})
})
