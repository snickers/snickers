package encoders

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"code.cloudfoundry.org/lager/lagertest"

	"github.com/dchest/uniuri"
	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

var _ = Describe("FFmpeg Encoder", func() {
	var (
		logger     *lagertest.TestLogger
		dbInstance db.Storage
		cfg        gonfig.Gonfig
	)

	BeforeEach(func() {
		currentDir, _ := os.Getwd()
		cfg, _ = gonfig.FromJsonFile(currentDir + "/../fixtures/config.json")
		dbInstance, _ = db.GetDatabase(cfg)
		dbInstance.ClearDatabase()
		logger = lagertest.NewTestLogger("ffmpeg-encoder")
	})

	Context("when calling", func() {

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

			err := FFMPEGEncode(logger, dbInstance, exampleJob.ID)
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
				LocalSource:      projectPath + "/../fixtures/videos/comingsoon.mov",
				LocalDestination: "/nowhere",
			}

			dbInstance.StoreJob(exampleJob)

			err := FFMPEGEncode(logger, dbInstance, exampleJob.ID)
			Expect(err.Error()).To(Equal("output format is not initialized. Unable to allocate context"))
		})

		It("Should change job status and Progress when encoding", func() {
			projectPath, _ := os.Getwd()
			swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://source.here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset: types.Preset{
					Name:        "presetHere",
					Container:   "mp4",
					RateControl: "vbr",
					Video: types.VideoPreset{
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
					Audio: types.AudioPreset{
						Codec:   "aac",
						Bitrate: "64000",
					},
				},
				Status:           types.JobCreated,
				Details:          "",
				LocalSource:      projectPath + "/../fixtures/videos/nyt.mp4",
				LocalDestination: swapDir + "/output.mp4",
			}

			dbInstance.StoreJob(exampleJob)

			FFMPEGEncode(logger, dbInstance, exampleJob.ID)
			changedJob, _ := dbInstance.RetrieveJob("123")

			Expect(changedJob.Progress).To(Equal("100%"))
			Expect(changedJob.Status).To(Equal(types.JobEncoding))
		})
	})

	Context("Regarding the application of presets", func() {
		It("should create h264/mp4 output", func() {
			currentDir, _ := os.Getwd()
			destinationFile := "/tmp/" + uniuri.New() + ".mp4"

			job := types.Job{
				ID: "123",
				Preset: types.Preset{
					Container:   "mp4", // OK
					RateControl: "vbr", // NOK
					Video: types.VideoPreset{
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
					Audio: types.AudioPreset{
						Codec:   "aac",   // OK
						Bitrate: "64000", // OK
					},
				},
				Status:           types.JobCreated,
				Progress:         "0%",
				LocalSource:      currentDir + "/../fixtures/videos/nyt.mp4",
				LocalDestination: destinationFile,
			}

			dbInstance.StoreJob(job)
			FFMPEGEncode(logger, dbInstance, job.ID)

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

			job := types.Job{
				ID: "123",
				Preset: types.Preset{
					Container:   "webm",
					RateControl: "vbr",
					Video: types.VideoPreset{
						Height:  "360",
						Width:   "640",
						Codec:   "vp8",
						Bitrate: "800000",
						GopSize: "90",
						GopMode: "fixed",
					},
					Audio: types.AudioPreset{
						Codec:   "vorbis",
						Bitrate: "64000",
					},
				},
				Status:           types.JobCreated,
				Progress:         "0%",
				LocalSource:      currentDir + "/../fixtures/videos/nyt.mp4",
				LocalDestination: destinationFile,
			}

			dbInstance.StoreJob(job)
			FFMPEGEncode(logger, dbInstance, job.ID)

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

			job := types.Job{
				ID: "123",
				Preset: types.Preset{
					Container:   "webm",
					RateControl: "vbr",
					Video: types.VideoPreset{
						Height:  "360",
						Width:   "640",
						Codec:   "vp9",
						Bitrate: "200000",
						GopSize: "90",
						GopMode: "fixed",
					},
					Audio: types.AudioPreset{
						Codec:   "vorbis",
						Bitrate: "64000",
					},
				},
				Status:           types.JobCreated,
				Progress:         "0%",
				LocalSource:      currentDir + "/../fixtures/videos/nyt.mp4",
				LocalDestination: destinationFile,
			}

			dbInstance.StoreJob(job)
			FFMPEGEncode(logger, dbInstance, job.ID)

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

			job := types.Job{
				ID: "123",
				Preset: types.Preset{
					Container:   "ogg",
					RateControl: "vbr",
					Video: types.VideoPreset{
						Height:  "360",
						Width:   "640",
						Codec:   "theora",
						Bitrate: "200000",
						GopSize: "90",
						GopMode: "fixed",
					},
					Audio: types.AudioPreset{
						Codec:   "vorbis",
						Bitrate: "64000",
					},
				},
				Status:           types.JobCreated,
				Progress:         "0%",
				LocalSource:      currentDir + "/../fixtures/videos/nyt.mp4",
				LocalDestination: destinationFile,
			}

			dbInstance.StoreJob(job)
			FFMPEGEncode(logger, dbInstance, job.ID)

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
			job := types.Job{
				Preset: types.Preset{
					Video: types.VideoPreset{
						Height: "360",
						Width:  "1000",
					},
					Audio: types.AudioPreset{},
				},
			}

			resultWidth, resultHeight := getResolution(job, 1280, 720)
			Expect(resultWidth).To(Equal(1000))
			Expect(resultHeight).To(Equal(360))
		})

		It("should maintain source's aspect ratio if one of the values on job.Preset is missing", func() {
			job1 := types.Job{
				Preset: types.Preset{
					Video: types.VideoPreset{
						Height: "360",
						Width:  "",
					},
					Audio: types.AudioPreset{},
				},
			}
			resultWidth, resultHeight := getResolution(job1, 1280, 720)
			Expect(resultWidth).To(Equal(640))
			Expect(resultHeight).To(Equal(360))

			job2 := types.Job{
				Preset: types.Preset{
					Video: types.VideoPreset{
						Height: "",
						Width:  "640",
					},
					Audio: types.AudioPreset{},
				},
			}
			resultWidth, resultHeight = getResolution(job2, 1280, 720)
			Expect(resultWidth).To(Equal(640))
			Expect(resultHeight).To(Equal(360))

			job3 := types.Job{
				Preset: types.Preset{
					Video: types.VideoPreset{
						Height: "",
						Width:  "",
					},
					Audio: types.AudioPreset{},
				},
			}
			resultWidth, resultHeight = getResolution(job3, 1280, 720)
			Expect(resultWidth).To(Equal(1280))
			Expect(resultHeight).To(Equal(720))
		})
	})
})
