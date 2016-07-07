package workouts_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/jwfriese/workouttrackerapi"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/workouts/{id}", func() {
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
	})

	Describe("GET", func() {
		var (
			response *http.Response
			getErr   error
			body     []byte
			readErr  error
		)

		Context("When a workout with the given ID exists in the database", func() {
			BeforeEach(func() {
				url := fmt.Sprintf("%v/%v", baseURL, "/workouts/3")
				response, getErr = http.Get(url)

				if getErr != nil {
					log.Fatal(getErr)
				}

				body, readErr = ioutil.ReadAll(response.Body)
				response.Body.Close()

				if readErr != nil {
					log.Fatal(readErr)
				}
			})

			It("returns a 200", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns the workout JSON-ified", func() {
				Expect(body).To(Equal([]byte(`{"id":3,"timestamp":"2016-03-11T06:12:56-08:00","lifts":[],"name":"turtle three"}`)))
			})
		})

		Context("When a workout with the given ID does not exist in the database", func() {
			BeforeEach(func() {
				url := fmt.Sprintf("%v/%v", baseURL, "/workouts/123")
				response, getErr = http.Get(url)

				if getErr != nil {
					log.Fatal(getErr)
				}
			})

			It("returns a 404", func() {
				Expect(response.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})
