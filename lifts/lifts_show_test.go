package lifts_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/jwfriese/workouttrackerapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/workouts/{id}/lifts/{id}", func() {
	var (
		testConnection *sql.DB
		testServer     *httptest.Server
		testClient     *http.Client
		baseURL        string
	)

	BeforeEach(func() {
		var err error
		testConnection, err = sql.Open("postgres", "user=postgres dbname=workout_tracker_integration_test sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		applicationHandler := workouttrackerapi.ApplicationHandler(testConnection)
		testServer = httptest.NewServer(applicationHandler)
		baseURL = testServer.URL
		testClient = http.DefaultClient
	})

	AfterEach(func() {
		testServer.Close()
	})

	Describe("GET", func() {
		var (
			response *http.Response
			err      error
			body     []byte
		)

		Context("When a lift with the given id exists in the database and on the given workout id", func() {
			BeforeEach(func() {
				url := fmt.Sprintf("%v/%v", baseURL, "/workouts/2/lifts/4")
				response, err = http.Get(url)

				if err != nil {
					log.Fatal(err)
				}

				body, err = ioutil.ReadAll(response.Body)
				response.Body.Close()

				if err != nil {
					log.Fatal(err)
				}
			})

			It("returns a 200", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns the lift JSON-ified", func() {
				Expect(body).To(MatchJSON([]byte(`{"id":4,"name":"turtle press","workout":2,"dataTemplate":"weight/reps","sets":[{"id":10,"dataTemplate":"weight/reps","lift":4,"weight":165,"reps":8},{"id":11,"dataTemplate":"weight/reps","lift":4,"weight":175,"reps":8},{"id":12,"dataTemplate":"weight/reps","lift":4,"weight":185,"reps":8}]}`)))
			})
		})

		Context("When a lift with the given id exists in the database, but it is not on the given workout", func() {
			BeforeEach(func() {
				url := fmt.Sprintf("%v%v", baseURL, "/workouts/1/lifts/4")
				response, err = http.Get(url)

				if err != nil {
					log.Fatal(err)
				}
			})

			It("returns a 404", func() {
				Expect(response.StatusCode).To(Equal(http.StatusNotFound))
			})

			It("returns an error repsonse", func() {
				defer response.Body.Close()
				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(responseBody).To(MatchJSON(`{"error":"Lift with id=4 does not exist on workout with id=1"}`))
			})
		})

		Context("When a lift with the given id exists in the database, but the workout does not", func() {
			BeforeEach(func() {
				url := fmt.Sprintf("%v%v", baseURL, "/workouts/1234/lifts/4")
				response, err = http.Get(url)

				if err != nil {
					log.Fatal(err)
				}
			})

			It("returns a 404", func() {
				Expect(response.StatusCode).To(Equal(http.StatusNotFound))
			})

			It("returns an error response", func() {
				defer response.Body.Close()
				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(responseBody).To(MatchJSON(`{"error":"Workout with id=1234 does not exist"}`))
			})
		})

		Context("When a lift with the given id does not exist in the database", func() {
			BeforeEach(func() {
				url := fmt.Sprintf("%v%v", baseURL, "/workouts/1/lifts/1000")
				response, err = http.Get(url)

				if err != nil {
					log.Fatal(err)
				}
			})

			It("returns a 404", func() {
				Expect(response.StatusCode).To(Equal(http.StatusNotFound))
			})

			It("returns an error response", func() {
				defer response.Body.Close()
				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(responseBody).To(MatchJSON(`{"error":"Lift with id=1000 does not exist"}`))
			})
		})
	})
})
