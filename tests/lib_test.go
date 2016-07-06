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
			Expect(changedJob.Details).To(SatisfyAny(ContainSubstring("no such host"), ContainSubstring("No filename could be determined")))
		})
	})
})
