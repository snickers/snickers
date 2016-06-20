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
			res, _ := dbInstance.StorePreset(examplePreset)
			Expect(res).To(Equal(expected))
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

			res, _ := dbInstance.RetrievePreset("presetOne")
			Expect(res).To(Equal(preset1))
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

			presets, _ := dbInstance.GetPresets()
			Expect(len(presets)).To(Equal(0))

			dbInstance.StorePreset(preset1)
			presets, _ = dbInstance.GetPresets()
			Expect(len(presets)).To(Equal(1))

			dbInstance.StorePreset(preset2)
			presets, _ = dbInstance.GetPresets()
			Expect(len(presets)).To(Equal(2))
		})

		It("should be able to update preset", func() {
			preset := types.Preset{
				Name:        "presetOne",
				Description: "This is preset one",
			}
			dbInstance.StorePreset(preset)

			expectedDescription := "New description for this preset"
			preset.Description = expectedDescription
			dbInstance.UpdatePreset("presetOne", preset)
			res, _ := dbInstance.GetPresets()

			Expect(res[0].Description).To(Equal(expectedDescription))
		})
	})
})
