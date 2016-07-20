package translation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWorkoutRequestTranslation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workout Request Translation Suite")
}
