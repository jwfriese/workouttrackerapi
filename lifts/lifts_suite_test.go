package lifts_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLiftsResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lifts Resource Suite")
}
