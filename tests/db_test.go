package snickers_test

import (
	"github.com/flavioribeiro/snickers/db/memory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	Context("in memory", func() {
		It("should be able to create and fetch a preset", func() {
			db, _ := memory.NewDatabase()
			Expect(db.CreatePreset()).To(Equal(0))
		})
	})

})
