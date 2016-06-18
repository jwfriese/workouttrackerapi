package lifts_test

import (
	"github.com/jwfriese/workouttracker/httpfakes"
	"github.com/jwfriese/workouttracker/lifts"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("lifts", func() {
	Describe("its endpoint handler", func() {
		var (
			handler            http.Handler
			fakeResponseWriter *httpfakes.FakeResponseWriter
		)

		BeforeEach(func() {
			fakeResponseWriter = new(httpfakes.FakeResponseWriter)
			handler = lifts.EndpointHandler()

			handler.ServeHTTP(fakeResponseWriter, nil)
		})

		It("returns a handler", func() {
			Expect(handler).ToNot(BeNil())
		})

		It("writes static content to the response", func() {
			Expect(fakeResponseWriter.WriteArgsForCall(0)).To(Equal([]byte("All the lifts")))
		})
	})
})
