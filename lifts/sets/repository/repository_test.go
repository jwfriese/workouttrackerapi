package repository_test

import (
	"database/sql"
	"log"

	"github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/sets/repository"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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

		subject = repository.NewSetRepository(testConnection)
	})

	Describe("Getting a set from the DB by its Id", func() {
		var (
			set *datamodel.Set
		)

		BeforeEach(func() {
			set = subject.GetById(1)
		})

		It("creates a Set model and returns it", func() {
			Expect(set.Id).To(Equal(1))
			Expect(set.DataTemplate).To(Equal("weight/reps"))
			Expect(set.Lift).To(Equal(1))
			Expect(set.Weight).ToNot(BeNil())
			Expect(*(set.Weight)).To(BeEquivalentTo(100.0))
			Expect(set.Height).To(BeNil())
			Expect(set.TimeInSeconds).To(BeNil())
			Expect(set.Reps).ToNot(BeNil())
			Expect(*(set.Reps)).To(Equal(10))
		})
	})
})
