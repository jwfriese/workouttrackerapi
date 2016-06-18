package workouttracker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWorkouttracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workouttracker Suite")
}
