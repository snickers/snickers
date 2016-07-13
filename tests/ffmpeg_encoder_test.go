package snickers_test

import (
	"os"
	"os/exec"
	"strings"

	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/lib"
	"github.com/snickers/snickers/types"
)

var _ = Describe("FFmpeg Encoder", func() {
	Context("when calling", func() {
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
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset: types.Preset{
					Name:         "presetHere",
					Container:    "mp4",
					Profile:      "main",
					ProfileLevel: "3.1",
					RateControl:  "VBR",
					Video: types.VideoPreset{
						Height:        "240",
						Width:         "426",
						Codec:         "h264",
						Bitrate:       "1000000",
						GopSize:       "90",
						GopMode:       "fixed",
						InterlaceMode: "progressive",
					},
					Audio: types.AudioPreset{
						Codec:   "aac",
						Bitrate: "64000",
					},
				},
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

	Context("Regarding the application of presets", func() {
		It("should create h264/mp4 output", func() {
			currentDir, _ := os.Getwd()

			job := types.Job{
				ID: "123",
				Preset: types.Preset{
					Container:    "mp4", // OK
					Profile:      "main",
					ProfileLevel: "3.1",
					RateControl:  "VBR",
					Video: types.VideoPreset{
						Height:        "240",
						Width:         "426",
						Codec:         "h264", // OK
						Bitrate:       "1000000",
						GopSize:       "90",
						GopMode:       "fixed",
						InterlaceMode: "progressive",
					},
					Audio: types.AudioPreset{
						Codec:   "aac", // OK
						Bitrate: "64000",
					},
				},
				Status:           types.JobCreated,
				Details:          "0%",
				LocalSource:      currentDir + "/videos/nyt.mp4",
				LocalDestination: "/tmp/o.mp4",
			}

			dbInstance, _ := db.GetDatabase()
			dbInstance.StoreJob(job)
			lib.FFMPEGEncode(job.ID)

			out, _ := exec.Command("mediainfo", "--Inform=General;%Format%;", "/tmp/o.mp4").Output()
			result := strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("mpeg-4"))

			out, _ = exec.Command("mediainfo", "--Inform=Audio;%Codec%;", "/tmp/o.mp4").Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("aac lc"))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Codec%;", "/tmp/o.mp4").Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("avc")) // AVC == H264
		})
	})

})
