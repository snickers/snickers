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

		It("should be able to store a job", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/destination.mp4",
				Preset:      types.Preset{Name: "presetHere"},
				Status:      "created",
			}
			expected := map[string]types.Job{"123": exampleJob}
			res, _ := dbInstance.StoreJob(exampleJob)
			Expect(res).To(Equal(expected))
		})

		It("should be able to retrieve a job by its name", func() {
			job1 := types.Job{
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/destination.mp4",
				Preset:      types.Preset{Name: "presetHere"},
				Status:      "created",
			}

			job2 := types.Job{
				ID:          "321",
				Source:      "http://source2.here.mp4",
				Destination: "s3://user@pass:/bucket/destination2.mp4",
				Preset:      types.Preset{Name: "presetHere2"},
				Status:      "created",
			}

			dbInstance.StoreJob(job1)
			dbInstance.StoreJob(job2)

			res, _ := dbInstance.RetrieveJob("123")
			Expect(res).To(Equal(job1))
		})

		It("should be able to list jobs", func() {
			job1 := types.Job{
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/destination.mp4",
				Preset:      types.Preset{Name: "presetHere"},
				Status:      "created",
			}

			job2 := types.Job{
				ID:          "321",
				Source:      "http://source2.here.mp4",
				Destination: "s3://user@pass:/bucket/destination2.mp4",
				Preset:      types.Preset{Name: "presetHere2"},
				Status:      "created",
			}

			jobs, _ := dbInstance.GetJobs()
			Expect(len(jobs)).To(Equal(0))

			dbInstance.StoreJob(job1)
			jobs, _ = dbInstance.GetJobs()
			Expect(len(jobs)).To(Equal(1))

			dbInstance.StoreJob(job2)
			jobs, _ = dbInstance.GetJobs()
			Expect(len(jobs)).To(Equal(2))
		})

		It("should be able to update job", func() {
			job1 := types.Job{
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/destination.mp4",
				Preset:      types.Preset{Name: "presetHere"},
				Status:      "created",
			}

			dbInstance.StoreJob(job1)

			expectedStatus := "downloading"
			job1.Status = expectedStatus
			dbInstance.UpdateJob("123", job1)
			res, _ := dbInstance.GetJobs()

			Expect(res[0].Status).To(Equal(expectedStatus))
		})
	})
})
