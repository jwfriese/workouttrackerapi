package sqlhelpers_test

import (
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Formatting numeric pointers", func() {
	var (
		result string
	)

	Describe("*float32", func() {
		Context("nonnil value", func() {
			BeforeEach(func() {
				ptr := new(float32)
				*ptr = 123.4
				result = sqlhelpers.Float32PointerToSQLString(ptr)
			})

			It("stringifies the float value", func() {
				Expect(result).To(Equal("123.4"))
			})
		})

		Context("nil value", func() {
			BeforeEach(func() {
				var ptr *float32
				result = sqlhelpers.Float32PointerToSQLString(ptr)
			})

			It("stringifies the NULL", func() {
				Expect(result).To(Equal("NULL"))
			})
		})
	})

	Describe("*int", func() {
		Context("nonnil value", func() {
			BeforeEach(func() {
				ptr := new(int)
				*ptr = 12
				result = sqlhelpers.IntPointerToSQLString(ptr)
			})

			It("stringifies the int value", func() {
				Expect(result).To(Equal("12"))
			})
		})

		Context("nil value", func() {
			BeforeEach(func() {
				var ptr *int
				result = sqlhelpers.IntPointerToSQLString(ptr)
			})

			It("stringifies the NULL", func() {
				Expect(result).To(Equal("NULL"))
			})
		})
	})
})
