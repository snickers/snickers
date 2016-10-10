package downloaders

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDownloaders(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Downloaders Suite")
}
