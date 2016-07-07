package sqlhelpers_test

import (
	"log"

	"github.com/jwfriese/workouttrackerapi/sqlhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers for numeric types with SQL", func() {
	Describe("[]int", func() {
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
