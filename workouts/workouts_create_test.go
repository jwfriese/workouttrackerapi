package workouts_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/jwfriese/workouttrackerapi"
	"github.com/jwfriese/workouttrackerapi/test/setup"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /workouts", func() {
	var (
		applicationHandler http.Handler
		testServer         *httptest.Server
		testConnection     *sql.DB
		baseURL            string
	)

	BeforeEach(func() {
		var err error
		testConnection, err = sql.Open("postgres", "user=postgres dbname=workout_tracker_integration_test sslmode=disable")

		if err != nil {
			log.Fatal(err)
		}

		applicationHandler = workouttrackerapi.ApplicationHandler(testConnection)
		testServer = httptest.NewServer(applicationHandler)
		baseURL = testServer.URL
	})

	AfterEach(func() {
		testServer.Close()
		setup.RefreshDatabase(testConnection)
	})

	Context("When the posted body contains valid JSON", func() {
		var (
			response *http.Response
		)

		BeforeEach(func() {
			buffer := bytes.NewBuffer([]byte(`{"name":"turtle abs workout","timestamp":"2016-04-25T10:12:56-08:00"}`))

			url := fmt.Sprintf("%v/%v", baseURL, "workouts")

			var err error
			response, err = http.Post(url, "application/json", buffer)
			if err != nil {
				log.Fatal(err)
			}
		})

		It("returns a 201", func() {
			Expect(response.StatusCode).To(Equal(http.StatusCreated))
		})

		It("returns the Location of the new workout resource", func() {
			locationURL, err := response.Location()
			Expect(err).To(BeNil())
			Expect(locationURL.Path).To(Equal("/workouts/4"))
		})

		It("returns the entirety of the new workout resource in the body", func() {
			responseBody, err := ioutil.ReadAll(response.Body)
			response.Body.Close()
			Expect(err).To(BeNil())
			Expect(responseBody).To(MatchJSON([]byte(`{"id":4,"timestamp":"2016-04-25T11:12:56-07:00","lifts":[],"name":"turtle abs workout"}`)))
		})
	})

	Context("When the request body has some error", func() {
		var (
			response *http.Response
		)

		BeforeEach(func() {
			buffer := bytes.NewBuffer([]byte(`{"name":"turtle abs workout","timestamp":"ill-formatted timestamp"}`))

			url := fmt.Sprintf("%v/%v", baseURL, "workouts")

			var err error
			response, err = http.Post(url, "application/json", buffer)
			if err != nil {
				log.Fatal(err)
			}
		})

		It("returns a 400", func() {
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("returns the error message as the 'error' value in JSON response", func() {
			responseBody, err := ioutil.ReadAll(response.Body)
			response.Body.Close()
			Expect(err).To(BeNil())

			var errorResponse map[string]interface{}
			err = json.Unmarshal(responseBody, &errorResponse)
			Expect(err).To(BeNil())

			Expect(errorResponse).ToNot(BeNil())
			Expect(errorResponse["error"]).To(Equal("pq: invalid input syntax for type timestamp with time zone: 'ill-formatted timestamp'"))
		})
	})
})
