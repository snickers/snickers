package db_test

import (
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/db/memory"
	"github.com/snickers/snickers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	var (
		dbInstance       db.DatabaseInterface
		runDatabaseSuite func()

		preset snickers.Preset
		job    snickers.Job
	)

	BeforeEach(func() {
		preset = snickers.Preset{
			Name:        "examplePreset",
			Description: "This is an example of preset",
			Container:   "mp4",
			RateControl: "vbr",
			Video: snickers.VideoPreset{
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
			Audio: snickers.AudioPreset{
				Codec:   "aac",
				Bitrate: "64000",
			},
		}

		job = snickers.Job{
			ID:          "123",
			Source:      "http://source.here.mp4",
			Destination: "s3://user@pass:/bucket/destination.mp4",
			Preset:      snickers.Preset{Name: "presetHere"},
			Status:      snickers.JobCreated,
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
			var anotherPreset snickers.Preset

			BeforeEach(func() {
				anotherPreset = snickers.Preset{
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
			var anotherJob snickers.Job

			BeforeEach(func() {
				anotherJob = snickers.Job{
					ID:          "321",
					Source:      "http://source2.here.mp4",
					Destination: "s3://user@pass:/bucket/destination2.mp4",
					Preset:      snickers.Preset{Name: "presetHere2"},
					Status:      snickers.JobCreated,
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
				expectedStatus := snickers.JobDownloading
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
			dbInstance, _ = memory.GetDatabase()
		})

		AfterEach(func() {
			dbInstance.ClearDatabase()
		})

		runDatabaseSuite()
	})

	// Describe("when the storage is mongodb", func() {
	// 	BeforeEach(func() {
	// 		dbInstance, _ = mongo.GetDatabase()
	// 	})

	// 	AfterEach(func() {
	// 		dbInstance.ClearDatabase()
	// 	})

	// 	runDatabaseSuite()
	// })
})
