package translation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTranslation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Translation Suite")
}
