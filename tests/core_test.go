package snickers_test

import (
	"os"
	"reflect"

	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/core"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers"
)

var _ = Describe("Core", func() {
	Context("Pipeline", func() {
		It("Should get the HTTPDownload function if source is HTTP", func() {
			jobSource := "http://flv.io/KailuaBeach.mp4"
			downloadFunc := core.GetDownloadFunc(jobSource)
			funcPointer := reflect.ValueOf(downloadFunc).Pointer()
			expected := reflect.ValueOf(core.HTTPDownload).Pointer()
			Expect(funcPointer).To(BeIdenticalTo(expected))
		})

		It("Should get the S3Download function if source is S3", func() {
			jobSource := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
			downloadFunc := core.GetDownloadFunc(jobSource)
			funcPointer := reflect.ValueOf(downloadFunc).Pointer()
			expected := reflect.ValueOf(core.S3Download).Pointer()
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
			bucket, _ := core.GetAWSBucket(destination)
			Expect(bucket).To(Equal("BUCKET"))
		})

		It("Should set credentials from URL Destination", func() {
			destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
			core.SetAWSCredentials(destination)
			Expect(os.Getenv("AWS_ACCESS_KEY_ID")).To(Equal("AWSKEY"))
			Expect(os.Getenv("AWS_SECRET_ACCESS_KEY")).To(Equal("AWSSECRET"))
		})

		It("Should get path and filename from URL Destination", func() {
			destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT/HERE.mp4"
			key, _ := core.GetAWSKey(destination)
			Expect(key).To(Equal("/OBJECT/HERE.mp4"))
		})
	})
})
