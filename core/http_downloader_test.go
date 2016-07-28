package core_test

import (
	"os"

	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/core"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers"
)

var _ = Describe("HTTP Downloader", func() {
	var (
		dbInstance db.DatabaseInterface
		cfg        gonfig.Gonfig
	)

	BeforeEach(func() {
		dbInstance, _ = db.GetDatabase()
		dbInstance.ClearDatabase()
		currentDir, _ := os.Getwd()
		cfg, _ = gonfig.FromJsonFile(currentDir + "/config.json")
	})

	It("should return an error if source couldn't be fetched", func() {
		exampleJob := snickers.Job{
			ID:          "123",
			Source:      "http://source.here.mp4",
			Destination: "s3://user@pass:/bucket/",
			Preset:      snickers.Preset{Name: "presetHere", Container: "mp4"},
			Status:      snickers.JobCreated,
			Details:     "",
		}
		dbInstance.StoreJob(exampleJob)

		err := core.HTTPDownload(exampleJob.ID)
		Expect(err.Error()).To(SatisfyAny(ContainSubstring("no such host"), ContainSubstring("No filename could be determined")))
	})

	It("Should set the local source and local destination on Job", func() {
		exampleJob := snickers.Job{
			ID:          "123",
			Source:      "http://flv.io/source_here.mp4",
			Destination: "s3://user@pass:/bucket/",
			Preset:      snickers.Preset{Name: "240p", Container: "mp4"},
			Status:      snickers.JobCreated,
			Details:     "",
		}
		dbInstance.StoreJob(exampleJob)

		core.HTTPDownload(exampleJob.ID)
		changedJob, _ := dbInstance.RetrieveJob("123")

		swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")
		sourceExpected := swapDir + "123/src/source_here.mp4"
		Expect(changedJob.LocalSource).To(Equal(sourceExpected))

		destinationExpected := swapDir + "123/dst/source_here_240p.mp4"
		Expect(changedJob.LocalDestination).To(Equal(destinationExpected))
	})
})
