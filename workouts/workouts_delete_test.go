package workouts_test

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/jwfriese/workouttrackerapi"
	"github.com/jwfriese/workouttrackerapi/test/setup"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DELETE /workouts/:id", func() {
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

	Context("When the workout with the given id exists in the database", func() {
		var (
			response *http.Response
			err      error
		)

		BeforeEach(func() {
			url := baseURL + "/workouts/1"
			deleteRequest, requestErr := http.NewRequest("DELETE", url, nil)
			if requestErr != nil {
				log.Fatal(requestErr)
			}

			response, err = http.DefaultClient.Do(deleteRequest)
		})

		It("returns a 204", func() {
			Expect(response.StatusCode).To(Equal(http.StatusNoContent))
		})

		It("deletes the workout from the database", func() {
			deletedRow := testConnection.QueryRow("SELECT id FROM workouts WHERE id=1")
			var workoutId int
			expectedErr := deletedRow.Scan(&workoutId)

			Expect(expectedErr).To(Equal(sql.ErrNoRows))
		})
	})

	Context("When no workout with the given id exists in the database", func() {
		var (
			response *http.Response
			err      error
		)

		BeforeEach(func() {
			url := baseURL + "/workouts/1111"
			deleteRequest, requestErr := http.NewRequest("DELETE", url, nil)
			if requestErr != nil {
				log.Fatal(requestErr)
			}

			response, err = http.DefaultClient.Do(deleteRequest)
		})

		It("returns a 404", func() {
			Expect(response.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Context("When given argument for id is incorrectly formatted", func() {
		var (
			response *http.Response
			err      error
		)

		BeforeEach(func() {
			url := baseURL + "/workouts/kittens"
			deleteRequest, requestErr := http.NewRequest("DELETE", url, nil)
			if requestErr != nil {
				log.Fatal(requestErr)
			}

			response, err = http.DefaultClient.Do(deleteRequest)
		})

		It("returns a 400", func() {
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})
})
