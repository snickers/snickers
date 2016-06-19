package snickers_test

import (
	"github.com/flavioribeiro/snickers/db/memory"
	"github.com/flavioribeiro/snickers/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	Context("in memory", func() {
		var (
			dbInstance *memory.Database
		)

		BeforeEach(func() {
			// TODO solve this by transforming database in interface/object (#6)
			dbInstance, _ = memory.GetDatabase()
			dbInstance.ClearDatabase()
		})

		It("should be able to store a preset", func() {
			examplePreset := types.Preset{
				Name:         "examplePreset",
				Description:  "This is an example of preset",
				Container:    "mp4",
				Profile:      "high",
				ProfileLevel: "3.1",
				RateControl:  "VBR",
				Video: types.VideoPreset{
					Width:         "720",
					Height:        "1080",
					Codec:         "h264",
					Bitrate:       "10000",
					GopSize:       "90",
					GopMode:       "fixed",
					InterlaceMode: "progressive",
				},
				Audio: types.AudioPreset{
					Codec:   "aac",
					Bitrate: "64000",
				},
			}
			expected := map[string]types.Preset{"examplePreset": examplePreset}
			Expect(dbInstance.StorePreset(examplePreset)).To(Equal(expected))
		})

		It("should be able to retrieve a preset by its name", func() {
			preset1 := types.Preset{
				Name:        "presetOne",
				Description: "This is preset one",
			}

			preset2 := types.Preset{
				Name:        "presetTwo",
				Description: "This is preset two",
			}

			dbInstance.StorePreset(preset1)
			dbInstance.StorePreset(preset2)

			Expect(dbInstance.RetrievePreset("presetOne")).To(Equal(preset1))
		})

		It("should be able to list presets", func() {
			preset1 := types.Preset{
				Name:        "presetOne",
				Description: "This is preset one",
			}

			preset2 := types.Preset{
				Name:        "presetTwo",
				Description: "This is preset two",
			}

			Expect(len(dbInstance.GetPresets())).To(Equal(0))

			dbInstance.StorePreset(preset1)
			Expect(len(dbInstance.GetPresets())).To(Equal(1))

			dbInstance.StorePreset(preset2)
			Expect(len(dbInstance.GetPresets())).To(Equal(2))
		})
	})
})
