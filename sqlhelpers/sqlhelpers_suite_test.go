package sqlhelpers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSqlhelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sqlhelpers Suite")
}
