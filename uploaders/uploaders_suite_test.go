package uploaders

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUploaders(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Uploaders Suite")
}
