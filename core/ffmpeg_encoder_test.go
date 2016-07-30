package core_test

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
	"github.com/snickers/snickers/types"
)

var _ = Describe("FFmpeg Encoder", func() {
	Context("when encoding", func() {
		var (
			dbInstance db.DatabaseInterface
			cfg        gonfig.Gonfig
			err        error
			job        types.Job
		)

		BeforeEach(func() {
			dbInstance, err = db.GetDatabase()
			Expect(err).NotTo(HaveOccurred())

			currentDir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			cfg, err = gonfig.FromJsonFile(currentDir + "/config.json")
			Expect(err).NotTo(HaveOccurred())

			projectPath, _ := os.Getwd()
			swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")
			job = types.Job{
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
				LocalSource:      projectPath + "/videos/nyt.mp4",
				LocalDestination: swapDir + "/output.mp4",
			}
		})

		JustBeforeEach(func() {
			dbInstance.StoreJob(job)
		})

		AfterEach(func() {
			dbInstance.ClearDatabase()
		})

		It("changes the job status and details", func() {
			core.FFMPEGEncode(job.ID)
			changedJob, _ := dbInstance.RetrieveJob("123")

			Expect(changedJob.Details).To(Equal("100%"))
			Expect(changedJob.Status).To(Equal(types.JobEncoding))
		})

		Context("when the input is not found", func() {
			BeforeEach(func() {
				job = types.Job{
					ID:               "123",
					Source:           "http://source.here.mp4",
					Destination:      "s3://user@pass:/bucket/",
					Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
					Status:           types.JobCreated,
					Details:          "",
					LocalSource:      "notfound.mp4",
					LocalDestination: "anywhere",
				}
			})

			It("returns an error", func() {
				err := core.FFMPEGEncode(job.ID)
				Expect(err.Error()).To(Equal("Error opening input 'notfound.mp4': No such file or directory"))
			})
		})

		Context("output path doesn't exists", func() {
			BeforeEach(func() {
				projectPath, _ := os.Getwd()
				job = types.Job{
					ID:               "123",
					Source:           "http://source.here.mp4",
					Destination:      "s3://user@pass:/bucket/",
					Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
					Status:           types.JobCreated,
					Details:          "",
					LocalSource:      projectPath + "/videos/comingsoon.mov",
					LocalDestination: "/nowhere",
				}
			})

			It("returns an error", func() {
				err := core.FFMPEGEncode(job.ID)
				Expect(err.Error()).To(Equal("output format is not initialized. Unable to allocate context"))
			})
		})
	})

	Context("Regarding the application of presets", func() {
		var (
			job             types.Job
			dbInstance      db.DatabaseInterface
			destinationFile string
			err             error
		)

		JustBeforeEach(func() {
			dbInstance, err = db.GetDatabase()
			Expect(err).NotTo(HaveOccurred())
			dbInstance.StoreJob(job)
		})

		AfterEach(func() {
			dbInstance.ClearDatabase()
		})

		checkMediaInfo := func(flag string, expect func(result string)) {
			out, _ := exec.Command("mediainfo", flag, destinationFile).Output()
			result := strings.Replace(strings.ToLower(string(out[:])), "\n", "", -1)
			expect(result)
		}

		Context("when h264/mp4", func() {
			BeforeEach(func() {
				currentDir, _ := os.Getwd()
				destinationFile = "/tmp/" + uniuri.New() + ".mp4"

				job = types.Job{
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
					Details:          "0%",
					LocalSource:      currentDir + "/videos/nyt.mp4",
					LocalDestination: destinationFile,
				}
			})

			It("creates h264/mp4 output", func() {
				core.FFMPEGEncode(job.ID)

				checkMediaInfo("--Inform=General;%Format%;", func(result string) {
					Expect(result).To(Equal("mpeg-4"))
				})

				checkMediaInfo("--Inform=Video;%Codec%;", func(result string) {
					Expect(result).To(Equal("avc")) // AVC == H264
				})

				checkMediaInfo("--Inform=Video;%ScanType%;", func(result string) {
					Expect(result).To(ContainSubstring(job.Preset.Video.InterlaceMode))
				})

				checkMediaInfo("--Inform=Video;%Format_Profile%;", func(result string) {
					Expect(result).To(ContainSubstring(job.Preset.Video.Profile))
				})

				checkMediaInfo("--Inform=Video;%Width%;", func(result string) {
					Expect(result).To(Equal(job.Preset.Video.Width))
				})

				checkMediaInfo("--Inform=Video;%Height%;", func(result string) {
					Expect(result).To(Equal(job.Preset.Video.Height))
				})

				checkMediaInfo("--Inform=Video;%BitRate_Nominal%;", func(result string) {
					Expect(result).To(Equal(job.Preset.Video.Bitrate))
				})

				checkMediaInfo("--Inform=Audio;%Codec%;", func(result string) {
					Expect(result).To(Equal("aac lc"))
				})

				checkMediaInfo("--Inform=Audio;%BitRate%;", func(result string) {
					resultInt, _ := strconv.Atoi(result)
					Expect(resultInt).To(SatisfyAll(BeNumerically(">", 50000), BeNumerically("<", 70000)))
				})
			})
		})

		Context("when webm/vp8", func() {
			BeforeEach(func() {
				currentDir, _ := os.Getwd()
				destinationFile = "/tmp/" + uniuri.New() + ".webm"

				job = types.Job{
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
					Details:          "0%",
					LocalSource:      currentDir + "/videos/nyt.mp4",
					LocalDestination: destinationFile,
				}
			})

			It("creates webm/vp8 output", func() {
				core.FFMPEGEncode(job.ID)

				checkMediaInfo("--Inform=General;%Format%;", func(result string) {
					Expect(result).To(Equal("webm"))
				})

				checkMediaInfo("--Inform=Video;%Codec%;", func(result string) {
					Expect(result).To(Equal("v_vp8"))
				})

				checkMediaInfo("--Inform=Video;%Width%;", func(result string) {
					Expect(result).To(Equal(job.Preset.Video.Width))
				})

				checkMediaInfo("--Inform=Video;%Height%;", func(result string) {
					Expect(result).To(Equal(job.Preset.Video.Height))
				})

				checkMediaInfo("--Inform=Audio;%Codec%;", func(result string) {
					Expect(result).To(Equal("vorbis"))
				})

				checkMediaInfo("--Inform=General;%BitRate%;", func(result string) {
					resultInt, _ := strconv.Atoi(result)
					Expect(resultInt).To(SatisfyAll(BeNumerically(">", 700000), BeNumerically("<", 900000)))
				})
			})
		})

		Context("when webm/vp9", func() {
			BeforeEach(func() {
				currentDir, _ := os.Getwd()
				destinationFile = "/tmp/" + uniuri.New() + ".webm"

				job = types.Job{
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
					Details:          "0%",
					LocalSource:      currentDir + "/videos/nyt.mp4",
					LocalDestination: destinationFile,
				}
			})

			It("creates webm/vp9 output", func() {
				core.FFMPEGEncode(job.ID)

				checkMediaInfo("--Inform=General;%Format%;", func(result string) {
					Expect(result).To(Equal("webm"))
				})

				checkMediaInfo("--Inform=Video;%Codec%;", func(result string) {
					Expect(result).To(Equal("v_vp9"))
				})

				checkMediaInfo("--Inform=Video;%Width%;", func(result string) {
					Expect(result).To(Equal(job.Preset.Video.Width))
				})

				checkMediaInfo("--Inform=Video;%Height%;", func(result string) {
					Expect(result).To(Equal(job.Preset.Video.Height))
				})

				checkMediaInfo("--Inform=Audio;%Codec%;", func(result string) {
					Expect(result).To(Equal("vorbis"))
				})

				checkMediaInfo("--Inform=General;%BitRate%;", func(result string) {
					resultInt, _ := strconv.Atoi(result)
					Expect(resultInt).To(SatisfyAll(BeNumerically(">", 100000), BeNumerically("<", 300000)))
				})
			})
		})

		Context("when ogg/theora", func() {
			BeforeEach(func() {
				currentDir, _ := os.Getwd()
				destinationFile = "/tmp/" + uniuri.New() + ".ogg"

				job = types.Job{
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
					Details:          "0%",
					LocalSource:      currentDir + "/videos/nyt.mp4",
					LocalDestination: destinationFile,
				}
			})

			It("creates ogg/theora output", func() {
				core.FFMPEGEncode(job.ID)

				checkMediaInfo("--Inform=General;%Format%;", func(result string) {
					Expect(result).To(Equal("ogg"))
				})

				checkMediaInfo("--Inform=Video;%Codec%;", func(result string) {
					Expect(result).To(Equal("theora"))
				})

				checkMediaInfo("--Inform=Video;%Width%;", func(result string) {
					Expect(result).To(Equal(job.Preset.Video.Width))
				})

				checkMediaInfo("--Inform=Video;%Height%;", func(result string) {
					Expect(result).To(Equal(job.Preset.Video.Height))
				})

				checkMediaInfo("--Inform=Audio;%Codec%;", func(result string) {
					Expect(result).To(Equal("vorbis"))
				})

				checkMediaInfo("--Inform=General;%BitRate%;", func(result string) {
					resultInt, _ := strconv.Atoi(result)
					Expect(resultInt).To(SatisfyAll(BeNumerically(">", 100000), BeNumerically("<", 400000)))
				})
			})
		})
	})

	Context("Regarding the definition of output resolution", func() {
		var (
			job types.Job
		)

		BeforeEach(func() {
			job = types.Job{
				Preset: types.Preset{
					Video: types.VideoPreset{
						Height: "360",
						Width:  "1000",
					},
					Audio: types.AudioPreset{},
				},
			}
		})

		It("returns width and height of job.Preset", func() {
			resultWidth, resultHeight := core.GetResolution(job, 1280, 720)
			Expect(resultWidth).To(Equal(1000))
			Expect(resultHeight).To(Equal(360))
		})

		Context("when one preset property is missing", func() {
			var (
				job types.Job
			)

			Context("when width is missing", func() {
				BeforeEach(func() {
					job = types.Job{
						Preset: types.Preset{
							Video: types.VideoPreset{
								Height: "360",
								Width:  "",
							},
							Audio: types.AudioPreset{},
						},
					}
				})

				It("maintains source's aspect ratio ", func() {
					resultWidth, resultHeight := core.GetResolution(job, 1280, 720)
					Expect(resultWidth).To(Equal(640))
					Expect(resultHeight).To(Equal(360))
				})
			})

			Context("when height is missing", func() {
				BeforeEach(func() {
					job = types.Job{
						Preset: types.Preset{
							Video: types.VideoPreset{
								Height: "",
								Width:  "640",
							},
							Audio: types.AudioPreset{},
						},
					}
				})

				It("maintains source's aspect ratio ", func() {
					resultWidth, resultHeight := core.GetResolution(job, 1280, 720)
					Expect(resultWidth).To(Equal(640))
					Expect(resultHeight).To(Equal(360))
				})
			})

			Context("when height and with are missing", func() {
				BeforeEach(func() {
					job = types.Job{
						Preset: types.Preset{
							Video: types.VideoPreset{
								Height: "",
								Width:  "",
							},
							Audio: types.AudioPreset{},
						},
					}

					It("maintains source's aspect ratio ", func() {
						resultWidth, resultHeight := core.GetResolution(job, 1280, 720)
						Expect(resultWidth).To(Equal(1280))
						Expect(resultHeight).To(Equal(720))
					})
				})
			})
		})
	})
})
