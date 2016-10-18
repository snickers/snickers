package uploaders

import (
	"reflect"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uploaders", func() {
	Context("GetUploadFunc", func() {
		It("should return S3Upload if source has amazonaws", func() {
			jobDestination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/source_here.mp4"
			uploadFunc := GetUploadFunc(jobDestination)
			funcName := runtime.FuncForPC(reflect.ValueOf(uploadFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/snickers/snickers/uploaders.S3Upload"))
		})

		It("should return FTPUpload if source starts with ftp://", func() {
			jobDestination := "ftp://login:password@host/source_here.mp4"
			uploadFunc := GetUploadFunc(jobDestination)
			funcName := runtime.FuncForPC(reflect.ValueOf(uploadFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/snickers/snickers/uploaders.FTPUpload"))
		})
	})
})
