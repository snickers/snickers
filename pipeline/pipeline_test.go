package pipeline

import (
	"io"
	"os"
	"reflect"

	"code.cloudfoundry.org/lager/lagertest"

	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/db/memory"
	"github.com/snickers/snickers/downloaders"
	"github.com/snickers/snickers/types"
)

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

var _ = Describe("Pipeline", func() {
	var (
		logger *lagertest.TestLogger
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
	})

	Context("Pipeline", func() {
		It("Should get the HTTPDownload function if source is HTTP", func() {
			jobSource := "http://flv.io/KailuaBeach.mp4"
			downloadFunc := GetDownloadFunc(jobSource)
			funcPointer := reflect.ValueOf(downloadFunc).Pointer()
			expected := reflect.ValueOf(downloaders.HTTPDownload).Pointer()
			Expect(funcPointer).To(BeIdenticalTo(expected))
		})

		It("Should get the S3Download function if source is S3", func() {
			jobSource := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
			downloadFunc := GetDownloadFunc(jobSource)
			funcPointer := reflect.ValueOf(downloadFunc).Pointer()
			expected := reflect.ValueOf(downloaders.S3Download).Pointer()
			Expect(funcPointer).To(BeIdenticalTo(expected))
		})
	})

	Context("when calling Swap Cleaner", func() {
		It("should remove local source and local destination", func() {
			dbInstance, _ := memory.GetDatabase()
			dbInstance.ClearDatabase()

			exampleJob := types.Job{
				ID:               "123",
				Source:           "http://source.here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
				Status:           types.JobCreated,
				Details:          "",
				LocalSource:      "/tmp/123/src/KailuaBeach.mp4",
				LocalDestination: "/tmp/123/dst/KailuaBeach.webm",
			}

			dbInstance.StoreJob(exampleJob)

			os.MkdirAll("/tmp/123/src/", 0777)
			os.MkdirAll("/tmp/123/dst/", 0777)

			cp(exampleJob.LocalSource, "./videos/nyt.mp4")
			cp(exampleJob.LocalDestination, "./videos/nyt.mp4")

			Expect(exampleJob.LocalSource).To(BeAnExistingFile())
			Expect(exampleJob.LocalDestination).To(BeAnExistingFile())

			CleanSwap(dbInstance, exampleJob.ID)

			Expect(exampleJob.LocalSource).To(Not(BeAnExistingFile()))
			Expect(exampleJob.LocalDestination).To(Not(BeAnExistingFile()))
		})
	})

	Context("HTTP Downloader", func() {
		var (
			dbInstance db.Storage
			cfg        gonfig.Gonfig
		)

		BeforeEach(func() {
			dbInstance, _ = memory.GetDatabase()
			dbInstance.ClearDatabase()
			currentDir, _ := os.Getwd()
			cfg, _ = gonfig.FromJsonFile(currentDir + "/config.json")
		})

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

			err := downloaders.HTTPDownload(logger, dbInstance, exampleJob.ID)
			Expect(err.Error()).To(SatisfyAny(ContainSubstring("no such host"), ContainSubstring("No filename could be determined")))
		})

		It("Should set the local source and local destination on Job", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://flv.io/source_here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
			dbInstance.StoreJob(exampleJob)

			downloaders.HTTPDownload(logger, dbInstance, exampleJob.ID)
			changedJob, _ := dbInstance.RetrieveJob("123")

			swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")
			sourceExpected := swapDir + "123/src/source_here.mp4"
			Expect(changedJob.LocalSource).To(Equal(sourceExpected))

			destinationExpected := swapDir + "123/dst/source_here_240p.mp4"
			Expect(changedJob.LocalDestination).To(Equal(destinationExpected))
		})
	})

})
