package repository_test

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/sets/repository"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func populateSetsDatabase(openConnection *sql.DB) {
	_, err := openConnection.Exec("INSERT INTO sets (data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',15,100.0,NULL,NULL,10)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("INSERT INTO sets (data_template,lift,weight,height,time_in_seconds,reps) VALUES ('height/reps',1,NULL,42,NULL,8)")
	if err != nil {
		log.Fatal(err)
	}
}

func clearSetsDatabase(openConnection *sql.DB) {
	_, err := openConnection.Exec("TRUNCATE sets")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("ALTER SEQUENCE sets_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatal(err)
	}
}

var _ = Describe("SetRepository", func() {
	var (
		subject        repository.SetRepository
		testConnection *sql.DB
	)

	BeforeEach(func() {
		var err error
		testConnection, err = sql.Open("postgres", "user=postgres dbname=workout_tracker_unit_test sslmode=disable")

		if err != nil {
			log.Fatal(err)
		}

		populateSetsDatabase(testConnection)
		subject = repository.NewSetRepository(testConnection)
	})

	AfterEach(func() {
		clearSetsDatabase(testConnection)
	})

	Describe("Getting a set from the database by its id", func() {
		var (
			set *datamodel.Set
			err error
		)

		Context("When a set with the given id exists in the database", func() {
			BeforeEach(func() {
				set, err = subject.GetById(1)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("creates a set model and returns it", func() {
				Expect(set.Id).To(Equal(1))
				Expect(set.DataTemplate).To(Equal("weight/reps"))
				Expect(set.Lift).To(Equal(15))
				Expect(set.Weight).ToNot(BeNil())
				Expect(*(set.Weight)).To(BeEquivalentTo(100.0))
				Expect(set.Height).To(BeNil())
				Expect(set.TimeInSeconds).To(BeNil())
				Expect(set.Reps).ToNot(BeNil())
				Expect(*(set.Reps)).To(Equal(10))
			})
		})

		Context("When no set with the given id exists in the database", func() {
			BeforeEach(func() {
				set, err = subject.GetById(1000)
			})

			It("returns no set", func() {
				Expect(set).To(BeNil())
			})

			It("returns a descriptive error message", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Set with id=1000 does not exist"))
			})
		})
	})

	Describe("Inserting a set into the database", func() {
		var (
			createdId int
			err       error
		)

		BeforeEach(func() {
			var weight *float32 = new(float32)
			*weight = 100.0
			var height *float32 = new(float32)
			*height = 125.5
			var timeInSeconds *float32 = new(float32)
			*timeInSeconds = 60
			var reps *int = new(int)
			*reps = 8
			set := &datamodel.Set{
				Id:            -1,
				DataTemplate:  "timeInSeconds",
				Lift:          54,
				Weight:        weight,
				Height:        height,
				TimeInSeconds: timeInSeconds,
				Reps:          reps,
			}

			createdId, err = subject.Insert(set)
		})

		It("returns no error", func() {
			Expect(err).To(BeNil())
		})

		It("returns the id of the created row", func() {
			Expect(createdId).To(Equal(3))
		})

		It("creates a full set in the database", func() {
			query := fmt.Sprintf("SELECT * FROM sets WHERE id=%v", createdId)
			row := testConnection.QueryRow(query)

			var id int
			var dataTemplate string
			var lift int
			var weight sql.NullFloat64
			var height sql.NullFloat64
			var timeInSeconds sql.NullFloat64
			var reps sql.NullInt64
			scanErr := row.Scan(&id, &dataTemplate, &lift, &weight, &height, &timeInSeconds, &reps)

			Expect(scanErr).To(BeNil())
			Expect(id).To(Equal(3))
			Expect(dataTemplate).To(Equal("timeInSeconds"))
			Expect(lift).To(Equal(54))

			weightValue, weightErr := weight.Value()
			Expect(weightErr).To(BeNil())
			Expect(weightValue).To(Equal(100.0))

			heightValue, heightErr := height.Value()
			Expect(heightErr).To(BeNil())
			Expect(heightValue).To(Equal(125.5))

			timeInSecondsValue, timeInSecondsErr := timeInSeconds.Value()
			Expect(timeInSecondsErr).To(BeNil())
			Expect(timeInSecondsValue).To(BeEquivalentTo(60))

			repsValue, repsErr := reps.Value()
			Expect(repsErr).To(BeNil())
			Expect(repsValue).To(BeEquivalentTo(8))
		})
	})

	Describe("Deleting a set from the database", func() {
		var (
			err error
		)

		Context("When a set with the given id exists in the database", func() {
			BeforeEach(func() {
				err = subject.Delete(1)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("deletes the set from the database", func() {
				row := testConnection.QueryRow("SELECT id FROM sets WHERE id=1")
				var setId int
				readErr := row.Scan(&setId)

				Expect(readErr).To(Equal(sql.ErrNoRows))
			})
		})

		Context("When no set with the given id exists in the database", func() {
			BeforeEach(func() {
				err = subject.Delete(19875387)
			})

			It("returns DoesNotExist error", func() {
				Expect(err).To(Equal(repository.ErrDoesNotExist))
			})
		})
	})
})
