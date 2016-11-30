package encoders

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/hls/segmenter"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

var _ = Describe("HLS Encoder", func() {
	var (
		logger     *lagertest.TestLogger
		dbInstance db.Storage
		cfg        gonfig.Gonfig
		exampleJob types.Job
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("http-download")
		currentDir, _ := os.Getwd()
		cfg, _ = gonfig.FromJsonFile(currentDir + "/../fixtures/config.json")
		dbInstance, _ = db.GetDatabase(cfg)
		dbInstance.ClearDatabase()
		exampleJob = types.Job{
			ID:               "123",
			Source:           "http://source.here.mp4",
			Destination:      "s3://user@pass:/bucket/",
			Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
			Status:           types.JobCreated,
			LocalSource:      "notfound.mp4",
			LocalDestination: "ftp://login:pass@url.com/",
		}
	})

	Context("when calling buildHLSConfig()", func() {
		It("should return the HLSConfig with right SourceFile and FileBase", func() {
			hlsConfig := buildHLSConfig(exampleJob)
			expectedHlsConfig := segmenter.HLSConfig{
				FileBase:        "ftp://login:pass@url.com/",
				SourceFile:      "notfound.mp4",
				SegmentDuration: 10,
			}
			Expect(hlsConfig).To(Equal(expectedHlsConfig))
		})
	})

	Context("when calling HLSEncode()", func() {
		It("should return error if job is not existent", func() {
			err := HLSEncode(logger, dbInstance, "non-existent-id")
			Expect(err.Error()).To(Equal("job not found"))
		})

		It("should return error if segmenting non-existent source", func() {
			dbInstance.StoreJob(exampleJob)
			err := HLSEncode(logger, dbInstance, exampleJob.ID)
			Expect(err.Error()).To(Equal("Error opening input 'notfound.mp4': No such file or directory"))
		})
	})

	Context("when calling encodeInH264", func() {
		It("should succeed and change localsource if source is present", func() {
			projectPath, _ := os.Getwd()
			exampleJob.LocalSource = projectPath + "/../fixtures/videos/nyt.mp4"
			exampleJob.Preset = types.Preset{
				Container:   "m3u8",
				RateControl: "vbr",
				Video: types.VideoPreset{
					Height:        "240",
					Width:         "426",
					Codec:         "h264",
					Bitrate:       "400000",
					GopSize:       "90",
					GopMode:       "fixed",
					Profile:       "main",
					ProfileLevel:  "3.1",
					InterlaceMode: "progressive",
				},
				Audio: types.AudioPreset{
					Codec:   "aac",
					Bitrate: "64000",
				},
			}

			dbInstance.StoreJob(exampleJob)
			err := encodeInH264(logger, dbInstance, exampleJob.ID)
			Expect(err).NotTo(HaveOccurred())
			job, _ := dbInstance.RetrieveJob(exampleJob.ID)
			Expect(strings.HasSuffix(job.LocalSource, "123.mp4")).To(BeTrue())
		})
	})
})
