package http_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test HTTP Suite")
}
