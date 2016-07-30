package core_test

import (
	"os"

	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/core"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

var _ = Describe("HTTP Downloader", func() {
	var (
		dbInstance db.DatabaseInterface
		err        error
		cfg        gonfig.Gonfig
		job        types.Job
		jobSource  string
	)

	BeforeEach(func() {
		dbInstance, err = db.GetDatabase()
		Expect(err).NotTo(HaveOccurred())

		currentDir, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		cfg, err = gonfig.FromJsonFile(currentDir + "/config.json")
		Expect(err).NotTo(HaveOccurred())
		jobSource = "http://flv.io/source_here.mp4"
	})

	JustBeforeEach(func() {
		job = types.Job{
			ID:          "123",
			Source:      jobSource,
			Destination: "s3://user@pass:/bucket/",
			Preset:      types.Preset{Name: "240p", Container: "mp4"},
			Status:      types.JobCreated,
			Details:     "",
		}
		dbInstance.StoreJob(job)
	})

	AfterEach(func() {
		dbInstance.ClearDatabase()
	})

	It("Should set the local source and local destination on Job", func() {
		core.HTTPDownload(job.ID)
		changedJob, _ := dbInstance.RetrieveJob("123")

		swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")
		sourceExpected := swapDir + "123/src/source_here.mp4"
		Expect(changedJob.LocalSource).To(Equal(sourceExpected))

		destinationExpected := swapDir + "123/dst/source_here_240p.mp4"
		Expect(changedJob.LocalDestination).To(Equal(destinationExpected))
	})

	Context("when could not be fetched", func() {
		BeforeEach(func() {
			jobSource = "http://source.here.mp4"
		})

		It("returns an error", func() {
			err := core.HTTPDownload(job.ID)
			Expect(err.Error()).To(SatisfyAny(ContainSubstring("no such host"), ContainSubstring("No filename could be determined")))
		})
	})
})
