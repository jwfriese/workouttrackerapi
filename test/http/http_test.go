package http_test

import (
	"github.com/jwfriese/workouttrackerapi/test/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("http", func() {
	Describe("Building a full URL for test", func() {
		var (
			fullURL string
		)

		Context("When forward slash is included in the given endpoint", func() {
			BeforeEach(func() {
				fullURL = http.CreateFullTestURLForEndpoint("/testendpoint")
			})

			It("builds the valid URL", func() {
				Expect(fullURL).To(Equal("localhost:8080/testendpoint"))
			})
		})

		Context("When forward slash is omitted from the given endpoint", func() {
			BeforeEach(func() {
				fullURL = http.CreateFullTestURLForEndpoint("testendpoint")
			})

			It("builds the valid URL, inserting the needed forward slash", func() {
				Expect(fullURL).To(Equal("localhost:8080/testendpoint"))
			})
		})
	})
})
