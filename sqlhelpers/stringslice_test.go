package sqlhelpers_test

import (
	"database/sql"
	"github.com/jwfriese/workouttracker/sqlhelpers"
	_ "github.com/lib/pq"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SQL helpers", func() {
	Describe("SQLStringSlice interface", func() {
		var (
			testConnection *sql.DB
		)

		BeforeEach(func() {
			var err error
			testConnection, err = sql.Open("postgres", "user=postgres dbname=workout_tracker_test sslmode=disable")
			if err != nil {
				log.Fatal(err)
			}
		})

		It("can scan a collection of strings from a SQL row", func() {
			rows, err := testConnection.Query("SELECT * FROM workouts WHERE id = 1")
			if err != nil {
				log.Fatal(err)
			}

			defer rows.Close()
			var id int
			var name string
			var timestamp string
			var stringSlice sqlhelpers.SQLStringSlice
			for rows.Next() {
				err = rows.Scan(&id, &name, &timestamp, &stringSlice)
				if err != nil {
					log.Fatal(err)
				}
			}

			Expect(stringSlice).ToNot(BeNil())

			resolvedStringSlice := stringSlice.ToStringSlice()
			Expect(len(resolvedStringSlice)).To(Equal(3))
			Expect(resolvedStringSlice[0]).To(Equal("turtle lift"))
			Expect(resolvedStringSlice[1]).To(Equal("turtle press"))
			Expect(resolvedStringSlice[2]).To(Equal("turtle push"))
		})
	})
})
