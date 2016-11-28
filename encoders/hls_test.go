package encoders

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/snickers/hls/segmenter"
	"github.com/snickers/snickers/types"
)

var _ = Describe("HLS Encoder", func() {
	Context("when calling buildHLSConfig()", func() {
		It("should return the HLSConfig with right SourceFile and FileBase", func() {
			exampleJob := types.Job{
				ID:               "123",
				Source:           "http://source.here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
				Status:           types.JobCreated,
				LocalSource:      "notfound.mp4",
				LocalDestination: "ftp://login:pass@url.com/",
			}
			hlsConfig := buildHLSConfig(exampleJob)
			expectedHlsConfig := segmenter.HLSConfig{
				FileBase:        "ftp://login:pass@url.com/",
				SourceFile:      "notfound.mp4",
				SegmentDuration: 10,
			}
			Expect(hlsConfig).To(Equal(expectedHlsConfig))
		})
	})
})
