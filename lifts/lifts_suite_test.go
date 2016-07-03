package lifts_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLifts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lifts Endpoint Suite")
}
