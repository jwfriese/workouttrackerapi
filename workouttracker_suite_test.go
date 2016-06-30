package workouttrackerapi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWorkouttrackerapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workouttrackerapi Suite")
}
