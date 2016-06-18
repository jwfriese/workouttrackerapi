package workouts_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWorkouts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workouts Suite")
}
