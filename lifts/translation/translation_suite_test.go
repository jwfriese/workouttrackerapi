package translation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLiftTranslation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lift Request Translation Suite")
}
