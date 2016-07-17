package workouts_test

import (
	"bytes"
	"database/sql"
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
			buffer := bytes.NewBuffer([]byte(`{"name":"turtle abs workout","timestamp":"2016-04-25T10:12:56-08:00","lifts":[{"name":"turtle crunches","dataTemplate":"weight/reps","sets":[{"dataTemplate":"weight/reps","weight":0,"reps":25},{"dataTemplate":"weight/reps","weight":0,"reps":35}]},{"name":"turtle hollow hold", "dataTemplate":"timeInSeconds","sets":[{"dataTemplate":"timeInSeconds","timeInSeconds":65.5},{"dataTemplate":"timeInSeconds","timeInSeconds":70}]},{"name":"turtle step-ups","dataTemplate":"height/reps","sets":[{"dataTemplate":"height/reps","height":36,"reps":50}]}]}`))

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
			Expect(responseBody).To(MatchJSON([]byte(`{"id":4,"timestamp":"2016-04-25T11:12:56-07:00","lifts":[{"id":6,"name":"turtle crunches","workout":4,"dataTemplate":"weight/reps","sets":[{"weight":0,"reps":25},{"weight":0,"reps":35}]},{"id":7,"name":"turtle hollow hold","workout":4,"dataTemplate":"timeInSeconds","sets":[{"timeInSeconds":65.5},{"timeInSeconds":70}]},{"id":8,"name":"turtle step-ups","workout":4,"dataTemplate":"height/reps","sets":[{"height":36,"reps":50}]}],"name":"turtle abs workout"}`)))
		})
	})
})
