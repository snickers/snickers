package encoders

import (
	"reflect"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/types"
)

var _ = Describe("Encoders", func() {
	Context("GetEncodeFunc", func() {
		It("should return HLSEncode if container is m3u8", func() {
			job := types.Job{
				ID:          "123",
				Source:      "ftp://login:password@host/source_here.mov",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "m3u8"},
				Status:      types.JobCreated,
			}
			encodeFunc := GetEncodeFunc(job)
			funcName := runtime.FuncForPC(reflect.ValueOf(encodeFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/snickers/snickers/encoders.HLSEncode"))
		})

		It("should return FFMPEGEncode if source is not m3u8", func() {
			job := types.Job{
				ID:          "123",
				Source:      "ftp://login:password@host/source_here.mov",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "mp4"},
				Status:      types.JobCreated,
			}
			encodeFunc := GetEncodeFunc(job)
			funcName := runtime.FuncForPC(reflect.ValueOf(encodeFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/snickers/snickers/encoders.FFMPEGEncode"))
		})
	})
})
