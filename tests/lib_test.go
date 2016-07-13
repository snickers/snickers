package snickers_test

import (
	"os"
	"reflect"

	"github.com/flavioribeiro/gonfig"
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/lib"
	"github.com/flavioribeiro/snickers/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Library", func() {
	Context("Pipeline", func() {
		It("Should get the HTTPDownload function if source is HTTP", func() {
			jobSource := "http://flv.io/KailuaBeach.mp4"
			downloadFunc := lib.GetDownloadFunc(jobSource)
			funcPointer := reflect.ValueOf(downloadFunc).Pointer()
			expected := reflect.ValueOf(lib.HTTPDownload).Pointer()
			Expect(funcPointer).To(BeIdenticalTo(expected))
		})

		It("Should get the S3Download function if source is S3", func() {
			jobSource := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
			downloadFunc := lib.GetDownloadFunc(jobSource)
			funcPointer := reflect.ValueOf(downloadFunc).Pointer()
			expected := reflect.ValueOf(lib.S3Download).Pointer()
			Expect(funcPointer).To(BeIdenticalTo(expected))
		})
	})

	Context("HTTP Downloader", func() {
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
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "presetHere", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
			dbInstance.StoreJob(exampleJob)

			err := lib.HTTPDownload(exampleJob.ID)
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

			lib.HTTPDownload(exampleJob.ID)
			changedJob, _ := dbInstance.RetrieveJob("123")

			swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")
			sourceExpected := swapDir + "123/src/source_here.mp4"
			Expect(changedJob.LocalSource).To(Equal(sourceExpected))

			destinationExpected := swapDir + "123/dst/source_here_240p.mp4"
			Expect(changedJob.LocalDestination).To(Equal(destinationExpected))
		})
	})

	Context("FFMPEG Encoder", func() {
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

		It("should return an error if input is not found", func() {
			exampleJob := types.Job{
				ID:               "123",
				Source:           "http://source.here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
				Status:           types.JobCreated,
				Details:          "",
				LocalSource:      "notfound.mp4",
				LocalDestination: "anywhere",
			}
			dbInstance.StoreJob(exampleJob)

			err := lib.FFMPEGEncode(exampleJob.ID)
			Expect(err.Error()).To(Equal("Error opening input 'notfound.mp4': No such file or directory"))
		})

		It("should return error if output path doesn't exists", func() {
			projectPath, _ := os.Getwd()
			exampleJob := types.Job{
				ID:               "123",
				Source:           "http://source.here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
				Status:           types.JobCreated,
				Details:          "",
				LocalSource:      projectPath + "/videos/comingsoon.mov",
				LocalDestination: "/nowhere",
			}

			dbInstance.StoreJob(exampleJob)

			err := lib.FFMPEGEncode(exampleJob.ID)
			Expect(err.Error()).To(Equal("output format is not initialized. Unable to allocate context"))
		})

		It("Should change job status and details when encoding", func() {
			projectPath, _ := os.Getwd()
			swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")
			exampleJob := types.Job{
				ID:               "123",
				Source:           "http://source.here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
				Status:           types.JobCreated,
				Details:          "",
				LocalSource:      projectPath + "/videos/nyt.mp4",
				LocalDestination: swapDir + "/output.mp4",
			}

			dbInstance.StoreJob(exampleJob)

			lib.FFMPEGEncode(exampleJob.ID)
			changedJob, _ := dbInstance.RetrieveJob("123")

			Expect(changedJob.Details).To(Equal("100%"))
			Expect(changedJob.Status).To(Equal(types.JobEncoding))
		})
	})

	Context("AWS Helpers", func() {
		var (
			dbInstance db.DatabaseInterface
		)

		BeforeEach(func() {
			dbInstance, _ = db.GetDatabase()
			dbInstance.ClearDatabase()
		})

		It("Should get bucket from URL Destination", func() {
			destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
			bucket, _ := lib.GetAWSBucket(destination)
			Expect(bucket).To(Equal("BUCKET"))
		})

		It("Should set credentials from URL Destination", func() {
			destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
			lib.SetAWSCredentials(destination)
			Expect(os.Getenv("AWS_ACCESS_KEY_ID")).To(Equal("AWSKEY"))
			Expect(os.Getenv("AWS_SECRET_ACCESS_KEY")).To(Equal("AWSSECRET"))
		})

		It("Should get path and filename from URL Destination", func() {
			destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT/HERE.mp4"
			key, _ := lib.GetAWSKey(destination)
			Expect(key).To(Equal("/OBJECT/HERE.mp4"))
		})
	})
})
