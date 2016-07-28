package core_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/core"
)

var _ = Describe("Pipeline", func() {
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
