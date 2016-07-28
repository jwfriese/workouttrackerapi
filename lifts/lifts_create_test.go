package lifts_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/jwfriese/workouttrackerapi"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
	"github.com/jwfriese/workouttrackerapi/test/setup"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lifts Resource: POST", func() {
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
		setup.RefreshDatabase(testConnection)
	})

	Describe("When the workout resource exists", func() {
		var (
			url      string
			response *http.Response
			err      error
		)

		BeforeEach(func() {
			url = fmt.Sprintf("%s/workouts/1/lifts", baseURL)
		})

		Describe("When given valid JSON", func() {
			BeforeEach(func() {
				body := bytes.NewBuffer([]byte(`{"name":"turtle box jumps","dataTemplate":"height/reps", "sets":[]}`))

				request, requestErr := http.NewRequest("POST", url, body)
				if requestErr != nil {
					log.Fatal(requestErr)
				}
				request.Header.Add("Content-Type", "application/json")
				response, err = testClient.Do(request)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("returns a 201", func() {
				Expect(response.StatusCode).To(Equal(http.StatusCreated))
			})

			It("returns the location of the new lift resource", func() {
				locationURL, locationErr := response.Location()
				Expect(locationErr).To(BeNil())
				Expect(locationURL.Path).To(Equal("/workouts/1/lifts/6"))
			})

			It("returns the entirety of the new lift resource in the response body", func() {
				responseBodyBytes, bodyReadErr := ioutil.ReadAll(response.Body)
				response.Body.Close()

				Expect(bodyReadErr).To(BeNil())
				Expect(responseBodyBytes).To(MatchJSON([]byte(`{"id":6,"name":"turtle box jumps","dataTemplate":"height/reps","workout":1,"sets":[]}`)))
			})

			It("creates a new lift in the database", func() {
				var (
					liftId       int
					name         string
					workoutId    int
					dataTemplate string
					setIds       sqlhelpers.IntSlice
				)

				newLiftRow := testConnection.QueryRow("SELECT * FROM lifts WHERE id=6")
				newLiftScanErr := newLiftRow.Scan(&liftId, &name, &workoutId, &dataTemplate, &setIds)

				Expect(newLiftScanErr).To(BeNil())
				Expect(liftId).To(Equal(6))
				Expect(name).To(Equal("turtle box jumps"))
				Expect(workoutId).To(Equal(1))
				Expect(dataTemplate).To(Equal("height/reps"))
				Expect(len(setIds)).To(Equal(0))
			})

			It("associates this new lift with the workout resource", func() {
				var liftIds sqlhelpers.IntSlice
				workoutRow := testConnection.QueryRow("SELECT lifts FROM workouts WHERE id=1")
				workoutScanErr := workoutRow.Scan(&liftIds)

				Expect(workoutScanErr).To(BeNil())
				Expect(len(liftIds)).To(Equal(4))
				Expect(liftIds).To(BeEquivalentTo([]int{1, 2, 3, 6}))
			})
		})

		Describe("When a required field is missing", func() {
			BeforeEach(func() {
				body := bytes.NewBuffer([]byte(`{"name":"turtle box jumps","sets":[]}`))

				request, requestErr := http.NewRequest("POST", url, body)
				if requestErr != nil {
					log.Fatal(requestErr)
				}
				request.Header.Add("Content-Type", "application/json")
				response, err = testClient.Do(request)
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
				Expect(errorResponse["error"]).To(Equal("Missing required 'dataTemplate' field in lift request JSON"))
			})
		})

		Describe("When the request cannot be read", func() {
			BeforeEach(func() {
				// Forgot to open the JSON correctly
				body := bytes.NewBuffer([]byte(`"name":"lift","dataTemplate":"weight/reps","sets":[]}`))
				request, requestErr := http.NewRequest("POST", url, body)
				if requestErr != nil {
					log.Fatal(requestErr)
				}
				request.Header.Add("Content-Type", "application/json")
				response, err = testClient.Do(request)
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
				Expect(errorResponse["error"]).To(Equal("invalid character ':' after top-level value"))
			})
		})
	})

	Describe("When the workout resource does not exist", func() {
		var (
			url      string
			response *http.Response
			err      error
		)

		BeforeEach(func() {
			url = fmt.Sprintf("%s/workouts/12345/lifts", baseURL)
			body := bytes.NewBuffer([]byte(`{"name":"turtle lift","dataTemplate":"weight/reps","sets":[]}`))

			request, requestErr := http.NewRequest("POST", url, body)
			if requestErr != nil {
				log.Fatal(requestErr)
			}

			response, err = testClient.Do(request)
		})

		It("returns a 404", func() {
			Expect(response.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("returns the error message as the 'error' value in JSON response", func() {
			responseBody, err := ioutil.ReadAll(response.Body)
			response.Body.Close()
			Expect(err).To(BeNil())

			var errorResponse map[string]interface{}
			err = json.Unmarshal(responseBody, &errorResponse)
			Expect(err).To(BeNil())

			Expect(errorResponse).ToNot(BeNil())
			Expect(errorResponse["error"]).To(Equal("Error inserting lift: Workout with id=12345 does not exist"))
		})
	})

	Describe("When the workout in the query is not the correct format", func() {
		var (
			url      string
			response *http.Response
			err      error
		)

		BeforeEach(func() {
			url = fmt.Sprintf("%s/workouts/turtles/lifts", baseURL)
			body := bytes.NewBuffer([]byte(`{"name":"turtle lift","dataTemplate":"weight/reps","sets":[]}`))

			request, requestErr := http.NewRequest("POST", url, body)
			if requestErr != nil {
				log.Fatal(requestErr)
			}

			response, err = testClient.Do(request)
		})

		It("returns a 404", func() {
			Expect(response.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("returns the error message as the 'error' value in JSON response", func() {
			responseBody, err := ioutil.ReadAll(response.Body)
			response.Body.Close()
			Expect(err).To(BeNil())

			var errorResponse map[string]interface{}
			err = json.Unmarshal(responseBody, &errorResponse)
			if err != nil {
				errString := fmt.Sprintf("Test expected to valid JSON: \n %s \n but produced error: \n %s", string(responseBody), err.Error())
				log.Fatal(errors.New(errString))
			}
			Expect(err).To(BeNil())

			Expect(errorResponse).ToNot(BeNil())
			Expect(errorResponse["error"]).To(Equal("Path component following '/workouts' must be a valid workout id"))
		})
	})
})
