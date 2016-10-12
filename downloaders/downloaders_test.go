package downloaders

import (
	"os"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/db/memory"
	"github.com/snickers/snickers/types"
)

var _ = Describe("Downloaders", func() {
	var (
		logger     *lagertest.TestLogger
		dbInstance db.Storage
		downloader DownloadFunc
		exampleJob types.Job
		configPath string
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("http-download")
		dbInstance, _ = memory.GetDatabase()
		dbInstance.ClearDatabase()
		currentDir, _ := os.Getwd()
		configPath = currentDir + "/../fixtures/config.json"
	})

	runDatabaseSuite := func() {
		It("should return an error if source couldn't be fetched", func() {
			dbInstance.StoreJob(exampleJob)
			err := downloader(logger, configPath, dbInstance, exampleJob.ID)
			Expect(err.Error()).To(SatisfyAny(
				ContainSubstring("no such host"),
				ContainSubstring("No filename could be determined"),
				ContainSubstring("The AWS Access Key Id you provided does not exist in our records")))
		})

		It("Should set the local source and local destination on Job", func() {
			dbInstance.StoreJob(exampleJob)
			downloader(logger, configPath, dbInstance, exampleJob.ID)
			changedJob, _ := dbInstance.RetrieveJob("123")

			cfg, _ := gonfig.FromJsonFile(configPath)
			swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")

			sourceExpected := swapDir + "123/src/source_here.mp4"
			Expect(changedJob.LocalSource).To(Equal(sourceExpected))

			destinationExpected := swapDir + "123/dst/source_here_240p.mp4"
			Expect(changedJob.LocalDestination).To(Equal(destinationExpected))
		})
	}

	Context("HTTP Downloader", func() {
		BeforeEach(func() {
			downloader = HTTPDownload
			exampleJob = types.Job{
				ID:          "123",
				Source:      "http://source_here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
		})

		runDatabaseSuite()
	})

	Context("FTP Downloader", func() {
		BeforeEach(func() {
			downloader = FTPDownload
			exampleJob = types.Job{
				ID:          "123",
				Source:      "ftp://login:password@host/source_here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
		})

		runDatabaseSuite()
	})

	Context("S3 Downloader", func() {
		BeforeEach(func() {
			downloader = S3Download
			exampleJob = types.Job{
				ID:          "123",
				Source:      "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/source_here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
		})

		runDatabaseSuite()
	})
})
