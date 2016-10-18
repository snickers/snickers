package helpers

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

var _ = Describe("Helpers", func() {
	Context("Path", func() {
		var (
			dbInstance db.Storage
			cfg        gonfig.Gonfig
		)

		BeforeEach(func() {
			currentDir, _ := os.Getwd()
			cfg, _ = gonfig.FromJsonFile(currentDir + "/../fixtures/config.json")
			dbInstance, _ = db.GetDatabase(cfg)
		})

		AfterEach(func() {
			dbInstance.ClearDatabase()
		})

		It("GetLocalSourcPath should return the correct local source path based on job", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://www.flv.io/KailuaBeach.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "presetHere", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
			dbInstance.StoreJob(exampleJob)

			Expect(GetLocalSourcePath(cfg, exampleJob.ID)).To(Equal("/tmp/123/src/"))
		})

		It("GetLocalDestination should return the correct local destination path based on job", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://www.flv.io/KailuaBeach.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "640x360", Container: "webm"},
				Status:      types.JobCreated,
				Details:     "",
			}
			dbInstance.StoreJob(exampleJob)

			Expect(GetLocalDestination(cfg, dbInstance, exampleJob.ID)).To(Equal("/tmp/123/dst/KailuaBeach_640x360.webm"))
		})

		It("GetOutputFilename should build output filename based on job and preset", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://www.flv.io/KailuaBeach.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "640x360", Container: "webm"},
				Status:      types.JobCreated,
				Details:     "",
			}
			dbInstance.StoreJob(exampleJob)

			Expect(GetOutputFilename(dbInstance, exampleJob.ID)).To(Equal("KailuaBeach_640x360.webm"))
		})
	})
})
