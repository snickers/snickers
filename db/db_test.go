package db

import (
	"errors"
	"strings"

	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	var (
		dbInstance       Storage
		runDatabaseSuite func()

		preset types.Preset
		job    types.Job
	)

	BeforeEach(func() {
		preset = types.Preset{
			Name:        "examplePreset",
			Description: "This is an example of preset",
			Container:   "mp4",
			RateControl: "vbr",
			Video: types.VideoPreset{
				Width:         "720",
				Height:        "1080",
				Codec:         "h264",
				Bitrate:       "10000",
				GopSize:       "90",
				GopMode:       "fixed",
				Profile:       "high",
				ProfileLevel:  "3.1",
				InterlaceMode: "progressive",
			},
			Audio: types.AudioPreset{
				Codec:   "aac",
				Bitrate: "64000",
			},
		}

		job = types.Job{
			ID:          "123",
			Source:      "http://source.here.mp4",
			Destination: "s3://user@pass:/bucket/destination.mp4",
			Preset:      types.Preset{Name: "presetHere"},
			Status:      types.JobCreated,
			Details:     "0%",
		}
	})

	runDatabaseSuite = func() {
		Describe("StorePreset", func() {
			It("should be able to store a preset", func() {
				res, err := dbInstance.StorePreset(preset)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(preset))
			})

			It("should fail to store a preset if it already exists", func() {
				dbInstance.StorePreset(preset)
				res, err := dbInstance.StorePreset(preset)
				Expect(res).To(Equal(types.Preset{}))
				Expect(err).To(Equal(errors.New("Error 409: Preset already exists, please update instead.")))

			})
		})

		Describe("RetrievePreset", func() {
			JustBeforeEach(func() {
				dbInstance.StorePreset(preset)
			})

			It("should be able to retrieve a preset by its name", func() {
				res, err := dbInstance.RetrievePreset("examplePreset")
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(preset))
			})

			Context("when the present does not exist", func() {
				It("should return an error", func() {
					_, err := dbInstance.RetrievePreset("invalid-preset")
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("GetPresets", func() {
			var anotherPreset types.Preset

			BeforeEach(func() {
				anotherPreset = types.Preset{
					Name:        "anotherPreset",
					Description: "This is another preset",
				}
			})

			JustBeforeEach(func() {
				dbInstance.StorePreset(preset)
				dbInstance.StorePreset(anotherPreset)
			})

			It("should be able to list presets", func() {
				presets, err := dbInstance.GetPresets()
				Expect(err).NotTo(HaveOccurred())
				Expect(presets).To(ConsistOf(preset, anotherPreset))
			})
		})

		Describe("UpdatePreset", func() {
			JustBeforeEach(func() {
				dbInstance.StorePreset(preset)
			})

			It("should be able to update preset", func() {
				expectedDescription := "New description for this preset"
				preset.Description = expectedDescription
				dbInstance.UpdatePreset(preset.Name, preset)
				res, err := dbInstance.GetPresets()
				Expect(err).NotTo(HaveOccurred())

				Expect(res[0].Description).To(Equal(expectedDescription))
			})

			Context("when the present does not exist", func() {
				It("should return an error", func() {
					_, err := dbInstance.RetrievePreset("invalid-preset")
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("DeletePreset", func() {
			JustBeforeEach(func() {
				dbInstance.StorePreset(preset)
			})

			It("should be able to delete a preset", func() {
				dbInstance.DeletePreset("examplePreset")
				presets, err := dbInstance.GetPresets()
				Expect(err).NotTo(HaveOccurred())
				Expect(presets).Should(BeEmpty())
			})
		})

		Describe("StoreJob", func() {
			It("should be able to store a job", func() {
				res, err := dbInstance.StoreJob(job)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(job))
			})
		})

		Describe("RetrieveJob", func() {
			JustBeforeEach(func() {
				dbInstance.StoreJob(job)
			})

			It("should be able to retrieve a job by its name", func() {
				res, err := dbInstance.RetrieveJob("123")
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(job))
			})

			Context("when the job does not exist", func() {
				It("should return an error", func() {
					_, err := dbInstance.RetrieveJob("invalid-job")
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("GetJobs", func() {
			var anotherJob types.Job

			BeforeEach(func() {
				anotherJob = types.Job{
					ID:          "321",
					Source:      "http://source2.here.mp4",
					Destination: "s3://user@pass:/bucket/destination2.mp4",
					Preset:      types.Preset{Name: "presetHere2"},
					Status:      types.JobCreated,
					Details:     "0%",
				}
			})

			JustBeforeEach(func() {
				dbInstance.StoreJob(job)
				dbInstance.StoreJob(anotherJob)
			})

			It("should be able to list jobs", func() {
				jobs, err := dbInstance.GetJobs()
				Expect(err).NotTo(HaveOccurred())
				Expect(jobs).To(ConsistOf(job, anotherJob))
			})
		})

		Describe("UpdateJob", func() {
			JustBeforeEach(func() {
				dbInstance.StoreJob(job)
			})

			It("should be able to update job", func() {
				expectedStatus := types.JobDownloading
				job.Status = expectedStatus
				dbInstance.UpdateJob(job.ID, job)

				res, err := dbInstance.GetJobs()
				Expect(err).NotTo(HaveOccurred())
				Expect(res[0].Status).To(Equal(expectedStatus))
			})
		})
	}

	Describe("when the storage is in memory", func() {
		BeforeEach(func() {
			cfg, _ := gonfig.FromJson(strings.NewReader(`{"DATABASE_DRIVER":"memory"}`))
			dbInstance, _ = GetDatabase(cfg)
		})

		AfterEach(func() {
			dbInstance.ClearDatabase()
		})

		runDatabaseSuite()
	})

	Describe("when the storage is mongodb", func() {
		Describe("When it connects", func() {
			BeforeEach(func() {
				cfg, _ := gonfig.FromJson(strings.NewReader(`{"DATABASE_DRIVER":"mongo", "MONGODB_HOST":"localhost"}`))
				dbInstance, _ = GetDatabase(cfg)
			})

			AfterEach(func() {
				dbInstance.ClearDatabase()
			})

			runDatabaseSuite()
		})

		Describe("When it fail to connect", func() {
			It("should not connect on mongo", func() {
				failedCfg, _ := gonfig.FromJson(strings.NewReader(`{"DATABASE_DRIVER":"mongo", "MONGODB_HOST":"invalid.ip.address"}`))
				_, err := GetDatabase(failedCfg)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
