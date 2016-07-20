package translation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSetTranslation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Set Request Translation Suite")
}
