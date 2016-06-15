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

		It("should be able to store a preset", func() {
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
			expected := map[string]db.Preset{"examplePreset": examplePreset}
			Expect(dbInstance.StorePreset(examplePreset)).To(Equal(expected))
		})

		It("should be able to retrieve a preset by its name", func() {
			preset1 := db.Preset{
				Name:        "presetOne",
				Description: "This is preset one",
			}

			preset2 := db.Preset{
				Name:        "presetTwo",
				Description: "This is preset two",
			}

			dbInstance.StorePreset(preset1)
			dbInstance.StorePreset(preset2)

			Expect(dbInstance.RetrievePreset("presetOne")).To(Equal(preset1))
		})
	})

})
