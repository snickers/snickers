package snickers_test

import (
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/db/memory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	Context("in memory", func() {
		var (
			dbInstance *memory.Database
		)

		BeforeEach(func() {
			dbInstance, _ = memory.NewDatabase()
		})

		It("should be able to create a preset", func() {
			examplePreset := db.Preset{
				Name:         "examplePreset",
				Description:  "This is an example of preset",
				Container:    "mp4",
				Profile:      "high",
				ProfileLevel: "3.1",
				RateControl:  "VBR",
				Video: db.VideoPreset{
					Width:         "720",
					Height:        "1080",
					Codec:         "h264",
					Bitrate:       "10000",
					GopSize:       "90",
					GopMode:       "fixed",
					InterlaceMode: "progressive",
				},
				Audio: db.AudioPreset{
					Codec:   "aac",
					Bitrate: "64000",
				},
			}
			expected := []db.Preset{examplePreset}
			Expect(dbInstance.CreatePreset(examplePreset)).To(Equal(expected))
		})
	})

})
