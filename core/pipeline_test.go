package core_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/core"
)

var _ = Describe("Pipeline", func() {
	var (
		jobSource    string
		downloadFunc core.DownloadFunc
		funcPointer  uintptr
	)

	JustBeforeEach(func() {
		downloadFunc = core.GetDownloadFunc(jobSource)
		funcPointer = reflect.ValueOf(downloadFunc).Pointer()
	})

	Context("when the source is http", func() {
		BeforeEach(func() {
			jobSource = "http://flv.io/KailuaBeach.mp4"
		})

		It("returns an HTTPDownload function", func() {
			expected := reflect.ValueOf(core.HTTPDownload).Pointer()
			Expect(funcPointer).To(BeIdenticalTo(expected))
		})
	})

	Context("when the source is S3", func() {
		BeforeEach(func() {
			jobSource = "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
		})

		It("returns S3Download function", func() {
			expected := reflect.ValueOf(core.S3Download).Pointer()
			Expect(funcPointer).To(BeIdenticalTo(expected))
		})
	})
})
