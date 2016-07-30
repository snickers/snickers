package snickers_test

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/core"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers"
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
			exampleJob := snickers.Job{
				ID:               "123",
				Source:           "http://source.here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           snickers.Preset{Name: "presetHere", Container: "mp4"},
				Status:           snickers.JobCreated,
				Details:          "",
				LocalSource:      "notfound.mp4",
				LocalDestination: "anywhere",
			}
			dbInstance.StoreJob(exampleJob)

			err := core.FFMPEGEncode(exampleJob.ID)
			Expect(err.Error()).To(Equal("Error opening input 'notfound.mp4': No such file or directory"))
		})

		It("should return error if output path doesn't exists", func() {
			projectPath, _ := os.Getwd()
			exampleJob := snickers.Job{
				ID:               "123",
				Source:           "http://source.here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           snickers.Preset{Name: "presetHere", Container: "mp4"},
				Status:           snickers.JobCreated,
				Details:          "",
				LocalSource:      projectPath + "/videos/comingsoon.mov",
				LocalDestination: "/nowhere",
			}

			dbInstance.StoreJob(exampleJob)

			err := core.FFMPEGEncode(exampleJob.ID)
			Expect(err.Error()).To(Equal("output format is not initialized. Unable to allocate context"))
		})

		It("Should change job status and details when encoding", func() {
			projectPath, _ := os.Getwd()
			swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")
			exampleJob := snickers.Job{
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset: snickers.Preset{
					Name:        "presetHere",
					Container:   "mp4",
					RateControl: "vbr",
					Video: snickers.VideoPreset{
						Height:        "240",
						Width:         "426",
						Codec:         "h264",
						Bitrate:       "1000000",
						GopSize:       "90",
						GopMode:       "fixed",
						Profile:       "main",
						ProfileLevel:  "3.1",
						InterlaceMode: "progressive",
					},
					Audio: snickers.AudioPreset{
						Codec:   "aac",
						Bitrate: "64000",
					},
				},
				Status:           snickers.JobCreated,
				Details:          "",
				LocalSource:      projectPath + "/videos/nyt.mp4",
				LocalDestination: swapDir + "/output.mp4",
			}

			dbInstance.StoreJob(exampleJob)

			core.FFMPEGEncode(exampleJob.ID)
			changedJob, _ := dbInstance.RetrieveJob("123")

			Expect(changedJob.Details).To(Equal("100%"))
			Expect(changedJob.Status).To(Equal(snickers.JobEncoding))
		})
	})

	Context("Regarding the application of presets", func() {
		It("should create h264/mp4 output", func() {
			currentDir, _ := os.Getwd()
			destinationFile := "/tmp/" + uniuri.New() + ".mp4"

			job := snickers.Job{
				ID: "123",
				Preset: snickers.Preset{
					Container:   "mp4", // OK
					RateControl: "vbr", // NOK
					Video: snickers.VideoPreset{
						Height:       "240",    // OK
						Width:        "426",    // OK
						Codec:        "h264",   // OK
						Bitrate:      "400000", // OK
						GopSize:      "90",     // NOK
						GopMode:      "fixed",  // NOK
						Profile:      "main",   // OK
						ProfileLevel: "3.1",    // NOK

						InterlaceMode: "progressive", // NOK
					},
					Audio: snickers.AudioPreset{
						Codec:   "aac",   // OK
						Bitrate: "64000", // OK
					},
				},
				Status:           snickers.JobCreated,
				Details:          "0%",
				LocalSource:      currentDir + "/videos/nyt.mp4",
				LocalDestination: destinationFile,
			}

			dbInstance, _ := db.GetDatabase()
			dbInstance.StoreJob(job)
			core.FFMPEGEncode(job.ID)

			out, _ := exec.Command("mediainfo", "--Inform=General;%Format%;", destinationFile).Output()
			result := strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("mpeg-4"))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Codec%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("avc")) // AVC == H264

			out, _ = exec.Command("mediainfo", "--Inform=Video;%ScanType%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(ContainSubstring(job.Preset.Video.InterlaceMode))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Format_Profile%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(ContainSubstring(job.Preset.Video.Profile))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Width%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal(job.Preset.Video.Width))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Height%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal(job.Preset.Video.Height))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%BitRate_Nominal%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal(job.Preset.Video.Bitrate))

			out, _ = exec.Command("mediainfo", "--Inform=Audio;%Codec%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("aac lc"))

			out, _ = exec.Command("mediainfo", "--Inform=Audio;%BitRate%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			resultInt, _ := strconv.Atoi(result)
			Expect(resultInt).To(SatisfyAll(BeNumerically(">", 50000), BeNumerically("<", 70000)))
		})

		It("should create webm/vp8 output", func() {
			currentDir, _ := os.Getwd()
			destinationFile := "/tmp/" + uniuri.New() + ".webm"

			job := snickers.Job{
				ID: "123",
				Preset: snickers.Preset{
					Container:   "webm",
					RateControl: "vbr",
					Video: snickers.VideoPreset{
						Height:  "360",
						Width:   "640",
						Codec:   "vp8",
						Bitrate: "800000",
						GopSize: "90",
						GopMode: "fixed",
					},
					Audio: snickers.AudioPreset{
						Codec:   "vorbis",
						Bitrate: "64000",
					},
				},
				Status:           snickers.JobCreated,
				Details:          "0%",
				LocalSource:      currentDir + "/videos/nyt.mp4",
				LocalDestination: destinationFile,
			}

			dbInstance, _ := db.GetDatabase()
			dbInstance.StoreJob(job)
			core.FFMPEGEncode(job.ID)

			out, _ := exec.Command("mediainfo", "--Inform=General;%Format%;", destinationFile).Output()
			result := strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("webm"))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Codec%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("v_vp8"))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Width%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal(job.Preset.Video.Width))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Height%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal(job.Preset.Video.Height))

			out, _ = exec.Command("mediainfo", "--Inform=Audio;%Codec%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("vorbis"))

			out, _ = exec.Command("mediainfo", "--Inform=General;%BitRate%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			resultInt, _ := strconv.Atoi(result)
			Expect(resultInt).To(SatisfyAll(BeNumerically(">", 700000), BeNumerically("<", 900000)))
		})

		It("should create webm/vp9 output", func() {
			currentDir, _ := os.Getwd()
			destinationFile := "/tmp/" + uniuri.New() + ".webm"

			job := snickers.Job{
				ID: "123",
				Preset: snickers.Preset{
					Container:   "webm",
					RateControl: "vbr",
					Video: snickers.VideoPreset{
						Height:  "360",
						Width:   "640",
						Codec:   "vp9",
						Bitrate: "200000",
						GopSize: "90",
						GopMode: "fixed",
					},
					Audio: snickers.AudioPreset{
						Codec:   "vorbis",
						Bitrate: "64000",
					},
				},
				Status:           snickers.JobCreated,
				Details:          "0%",
				LocalSource:      currentDir + "/videos/nyt.mp4",
				LocalDestination: destinationFile,
			}

			dbInstance, _ := db.GetDatabase()
			dbInstance.StoreJob(job)
			core.FFMPEGEncode(job.ID)

			out, _ := exec.Command("mediainfo", "--Inform=General;%Format%;", destinationFile).Output()
			result := strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("webm"))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Codec%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("v_vp9"))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Width%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal(job.Preset.Video.Width))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Height%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal(job.Preset.Video.Height))

			out, _ = exec.Command("mediainfo", "--Inform=Audio;%Codec%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("vorbis"))

			out, _ = exec.Command("mediainfo", "--Inform=General;%BitRate%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			resultInt, _ := strconv.Atoi(result)
			Expect(resultInt).To(SatisfyAll(BeNumerically(">", 100000), BeNumerically("<", 300000)))
		})

		It("should create ogg/theora output", func() {
			currentDir, _ := os.Getwd()
			destinationFile := "/tmp/" + uniuri.New() + ".ogg"

			job := snickers.Job{
				ID: "123",
				Preset: snickers.Preset{
					Container:   "ogg",
					RateControl: "vbr",
					Video: snickers.VideoPreset{
						Height:  "360",
						Width:   "640",
						Codec:   "theora",
						Bitrate: "200000",
						GopSize: "90",
						GopMode: "fixed",
					},
					Audio: snickers.AudioPreset{
						Codec:   "vorbis",
						Bitrate: "64000",
					},
				},
				Status:           snickers.JobCreated,
				Details:          "0%",
				LocalSource:      currentDir + "/videos/nyt.mp4",
				LocalDestination: destinationFile,
			}

			dbInstance, _ := db.GetDatabase()
			dbInstance.StoreJob(job)
			core.FFMPEGEncode(job.ID)

			out, _ := exec.Command("mediainfo", "--Inform=General;%Format%;", destinationFile).Output()
			result := strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("ogg"))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Codec%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("theora"))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Width%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal(job.Preset.Video.Width))

			out, _ = exec.Command("mediainfo", "--Inform=Video;%Height%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal(job.Preset.Video.Height))

			out, _ = exec.Command("mediainfo", "--Inform=Audio;%Codec%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			Expect(result).To(Equal("vorbis"))

			out, _ = exec.Command("mediainfo", "--Inform=General;%BitRate%;", destinationFile).Output()
			result = strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			resultInt, _ := strconv.Atoi(result)
			Expect(resultInt).To(SatisfyAll(BeNumerically(">", 100000), BeNumerically("<", 400000)))
		})
	})
	Context("Regarding the definition of output resolution", func() {
		It("should return width and height of job.Preset", func() {
			job := snickers.Job{
				Preset: snickers.Preset{
					Video: snickers.VideoPreset{
						Height: "360",
						Width:  "1000",
					},
					Audio: snickers.AudioPreset{},
				},
			}

			resultWidth, resultHeight := core.GetResolution(job, 1280, 720)
			Expect(resultWidth).To(Equal(1000))
			Expect(resultHeight).To(Equal(360))
		})

		It("should maintain source's aspect ratio if one of the values on job.Preset is missing", func() {
			job1 := snickers.Job{
				Preset: snickers.Preset{
					Video: snickers.VideoPreset{
						Height: "360",
						Width:  "",
					},
					Audio: snickers.AudioPreset{},
				},
			}
			resultWidth, resultHeight := core.GetResolution(job1, 1280, 720)
			Expect(resultWidth).To(Equal(640))
			Expect(resultHeight).To(Equal(360))

			job2 := snickers.Job{
				Preset: snickers.Preset{
					Video: snickers.VideoPreset{
						Height: "",
						Width:  "640",
					},
					Audio: snickers.AudioPreset{},
				},
			}
			resultWidth, resultHeight = core.GetResolution(job2, 1280, 720)
			Expect(resultWidth).To(Equal(640))
			Expect(resultHeight).To(Equal(360))

			job3 := snickers.Job{
				Preset: snickers.Preset{
					Video: snickers.VideoPreset{
						Height: "",
						Width:  "",
					},
					Audio: snickers.AudioPreset{},
				},
			}
			resultWidth, resultHeight = core.GetResolution(job3, 1280, 720)
			Expect(resultWidth).To(Equal(1280))
			Expect(resultHeight).To(Equal(720))
		})
	})
})
