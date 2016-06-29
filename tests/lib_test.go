package snickers_test

import (
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/lib"
	"github.com/flavioribeiro/snickers/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func NextStep(job types.Job) {}

var _ = Describe("Library", func() {
	Context("Helpers", func() {
		var (
			dbInstance db.DatabaseInterface
		)

		BeforeEach(func() {
			dbInstance, _ = db.GetDatabase()
			dbInstance.ClearDatabase()
		})

		It("ChangeJobStatus should change job status", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/destination.mp4",
				Preset:      types.Preset{Name: "presetHere"},
				Status:      types.JobCreated,
				Details:     "",
			}
			dbInstance.StoreJob(exampleJob)

			lib.ChangeJobStatus(exampleJob.ID, types.JobEncoding)
			changedJob, _ := dbInstance.RetrieveJob(exampleJob.ID)

			Expect(changedJob.Status).To(Equal(types.JobEncoding))
		})

		It("ChangeJobDetails should change job details", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/destination.mp4",
				Preset:      types.Preset{Name: "presetHere"},
				Status:      types.JobCreated,
				Details:     "0%",
			}
			dbInstance.StoreJob(exampleJob)

			lib.ChangeJobDetails(exampleJob.ID, "100%")
			changedJob, _ := dbInstance.RetrieveJob(exampleJob.ID)

			Expect(changedJob.Details).To(Equal("100%"))
		})
	})

	Context("Download", func() {
		var (
			dbInstance db.DatabaseInterface
		)

		BeforeEach(func() {
			dbInstance, _ = db.GetDatabase()
			dbInstance.ClearDatabase()
		})

		It("Should change job status and details on error", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/destination.mp4",
				Preset:      types.Preset{Name: "presetHere"},
				Status:      types.JobCreated,
				Details:     "",
			}
			dbInstance.StoreJob(exampleJob)

			lib.Download(exampleJob.ID, NextStep)
			changedJob, _ := dbInstance.RetrieveJob("123")

			Expect(changedJob.Status).To(Equal(types.JobError))
			Expect(changedJob.Details).To(Equal("Head http://source.here.mp4: dial tcp: lookup source.here.mp4: no such host"))
		})
	})
})
