package downloaders

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

var _ = Describe("HTTP downloader", func() {
	var (
		logger     *lagertest.TestLogger
		dbInstance db.Storage
		cfg        gonfig.Gonfig
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("http-download")
		currentDir, _ := os.Getwd()
		cfg, _ = gonfig.FromJsonFile(currentDir + "/../fixtures/config.json")
		dbInstance, _ = db.GetDatabase(cfg)
		dbInstance.ClearDatabase()
	})

	Context("HTTP Downloader", func() {
		It("should return an error if source couldn't be fetched", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://source_here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "presetHere", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
			dbInstance.StoreJob(exampleJob)
			err := HTTPDownload(logger, cfg, dbInstance, exampleJob.ID)
			Expect(err.Error()).To(SatisfyAny(ContainSubstring("no such host"), ContainSubstring("No filename could be determined")))
		})
	})
})
