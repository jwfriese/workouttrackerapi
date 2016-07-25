package lifts_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/jwfriese/workouttrackerapi"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
	"github.com/jwfriese/workouttrackerapi/test/setup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /lifts", func() {
	var (
		testConnection *sql.DB
		testServer     *httptest.Server
		testClient     *http.Client
		baseURL        string
	)

	BeforeEach(func() {
		var setupErr error
		testConnection, setupErr = sql.Open("postgres", "user=postgres dbname=workout_tracker_integration_test sslmode=disable")
		if setupErr != nil {
			log.Fatal(setupErr)
		}

		applicationHandler := workouttrackerapi.ApplicationHandler(testConnection)
		testServer = httptest.NewServer(applicationHandler)
		baseURL = testServer.URL
		testClient = http.DefaultClient
	})

	AfterEach(func() {
		testServer.Close()
		setup.RefreshDatabase(testConnection)
	})

	Describe("When the request is valid", func() {
		var (
			response *http.Response
			err      error
		)

		BeforeEach(func() {
			url := fmt.Sprintf("%s/lifts", baseURL)
			body := bytes.NewBuffer([]byte(`{"name":"turtle weighted pull-ups","dataTemplate":"weight/reps","workout":2,"sets":[{"dataTemplate":"weight/reps","weight":125.0,"reps":10}]}`))
			request, requestErr := http.NewRequest("POST", url, body)
			if requestErr != nil {
				log.Fatal(requestErr)
			}
			request.Header.Add("ContentType", "application/json")
			response, err = testClient.Do(request)
		})

		It("returns no error", func() {
			Expect(err).To(BeNil())
		})

		It("returns a 201", func() {
			Expect(response.StatusCode).To(Equal(http.StatusCreated))
		})

		It("returns the Location of the new lift resource", func() {
			locationURL, err := response.Location()
			Expect(err).To(BeNil())
			Expect(locationURL.Path).To(Equal("/lifts/6"))
		})

		It("returns a full JSON representation of the created lift in the response", func() {
			bodyBytes, readErr := ioutil.ReadAll(response.Body)
			response.Body.Close()

			Expect(readErr).To(BeNil())
			Expect(bodyBytes).To(MatchJSON([]byte(`{"id":6,"name":"turtle weighted pull-ups","workout":2,"dataTemplate":"weight/reps","sets":[{"weight":125.0,"reps":10}]}`)))
		})

		It("creates a new lift with the given information in the database", func() {
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
			Expect(name).To(Equal("turtle weighted pull-ups"))
			Expect(workoutId).To(Equal(2))
			Expect(dataTemplate).To(Equal("weight/reps"))
			Expect(setIds).To(BeEquivalentTo([]int{16}))
		})

		It("associates the new lift with the workout specified by the request", func() {
			workoutRow := testConnection.QueryRow("SELECT lifts FROM workouts WHERE id=2")
			var liftIds sqlhelpers.IntSlice

			workoutScanErr := workoutRow.Scan(&liftIds)
			Expect(workoutScanErr).To(BeNil())
			Expect(liftIds).To(BeEquivalentTo([]int{4, 5, 6}))
		})

		It("creates sets in the database for all those included in the request", func() {
			setRow := testConnection.QueryRow("SELECT id, weight, reps FROM sets WHERE id=16")
			var (
				setId          int
				nullableWeight sql.NullFloat64
				nullableReps   sql.NullInt64
			)
			setScanErr := setRow.Scan(&setId, &nullableWeight, &nullableReps)
			Expect(setScanErr).To(BeNil())
			Expect(setId).To(Equal(16))

			weightValue, weightErr := nullableWeight.Value()
			Expect(weightErr).To(BeNil())
			Expect(weightValue).To(Equal(125.0))

			repsValue, repsErr := nullableReps.Value()
			Expect(repsErr).To(BeNil())
			Expect(repsValue).To(BeEquivalentTo(10))
		})
	})

	Describe("When the request has some error", func() {

	})
})
