package sqlhelpers_test

import (
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SQL helpers", func() {
	Describe("SQLStringSlice interface", func() {
		var (
			err    error
			result sqlhelpers.SQLStringSlice
		)

		BeforeEach(func() {
			input := []byte(`{"turtle lift","turtle press","turtle push"}`)
			err = (&result).Scan(input)
		})

		It("pulls all of the strings out from the scanned []byte", func() {
			Expect(result).ToNot(BeNil())

			resolvedStringSlice := result.ToStringSlice()
			Expect(len(resolvedStringSlice)).To(Equal(3))
			Expect(resolvedStringSlice[0]).To(Equal("turtle lift"))
			Expect(resolvedStringSlice[1]).To(Equal("turtle press"))
			Expect(resolvedStringSlice[2]).To(Equal("turtle push"))
		})
	})
})
