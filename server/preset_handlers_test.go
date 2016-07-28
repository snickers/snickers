package server_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pivotal-golang/lager/lagertest"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/server"
	"github.com/snickers/snickers/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PresetHandlers", func() {
	var (
		httpClient     *http.Client
		snickersServer *server.SnickersServer
		err            error
		log            *lagertest.TestLogger
		tmpDir         string
		dbInstance     db.DatabaseInterface
	)

	BeforeEach(func() {
		log = lagertest.NewTestLogger("snickers-test")
		tmpDir, err = ioutil.TempDir(os.TempDir(), "snickers-server-test")
		socketPath := path.Join(tmpDir, "snickers.sock")

		dbInstance, _ = db.GetDatabase()

		snickersServer = server.New(log, "unix", socketPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(snickersServer.Start()).NotTo(HaveOccurred())
		httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout("unix", socketPath, 2*time.Second)
				},
			},
		}
	})

	AfterEach(func() {
		if tmpDir != "" {
			os.RemoveAll(tmpDir)
		}

		dbInstance.ClearDatabase()
	})

	Describe("ListPresets", func() {
		It("should return application/json on its content type", func() {
			request, _ := http.NewRequest(http.MethodGet,
				fmt.Sprintf("http://server%s", server.Routes[server.ListPresets].Path),
				nil)
			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})

		It("should return stored presets", func() {
			examplePreset1 := types.Preset{Name: "a"}
			examplePreset2 := types.Preset{Name: "b"}
			dbInstance.StorePreset(examplePreset1)
			dbInstance.StorePreset(examplePreset2)

			expected1, _ := json.Marshal(`[{"name":"a","video":{},"audio":{}},{"name":"b","video":{},"audio":{}}]`)
			expected2, _ := json.Marshal(`[{"name":"b","video":{},"audio":{}},{"name":"a","video":{},"audio":{}}]`)

			request, _ := http.NewRequest(http.MethodGet,
				fmt.Sprintf("http://server%s", server.Routes[server.ListPresets].Path),
				nil)
			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			rb, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())
			rb2 := string(rb)
			responseBody, _ := json.Marshal(rb2)

			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(responseBody).To(SatisfyAny(Equal(expected1), Equal(expected2)))
		})
	})

	Describe("GetPresetDetails", func() {
		It("should return the preset with details", func() {
			examplePreset := types.Preset{
				Name:        "examplePreset",
				Description: "This is an example of preset",
				Container:   "mp4",
				RateControl: "vbr",
				Video: types.VideoPreset{
					Width:        "720",
					Height:       "1080",
					Codec:        "h264",
					Bitrate:      "10000",
					GopSize:      "90",
					GopMode:      "fixed",
					Profile:      "high",
					ProfileLevel: "3.1",

					InterlaceMode: "progressive",
				},
				Audio: types.AudioPreset{
					Codec:   "aac",
					Bitrate: "64000",
				},
			}
			dbInstance.StorePreset(examplePreset)
			expected, _ := json.Marshal(examplePreset)

			request, _ := http.NewRequest(http.MethodGet,
				"http://server/presets/examplePreset",
				nil)
			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			rb, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())
			responseBody := string(rb)

			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(responseBody).To(Equal(string(expected)))
		})

		Context("when getting the preset fails", func() {
			It("should return BadRequest", func() {
				request, _ := http.NewRequest(http.MethodGet,
					"http://server/presets/yoyoyo",
					nil)
				response, err := httpClient.Do(request)
				Expect(err).NotTo(HaveOccurred())
				rb, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				rb2 := string(rb)

				expected, _ := json.Marshal(`{"error": "retrieving preset: preset not found"}`)
				responseBody, _ := json.Marshal(rb2)
				Expect(responseBody).To(Equal(expected))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("CreatePreset", func() {
		It("should save a new preset", func() {
			preset := `{"name": "storedPreset", "video": {},"audio": {}}`
			request, _ := http.NewRequest(http.MethodPost,
				"http://server/presets",
				strings.NewReader(preset))
			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			rb, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())
			rb2 := string(rb)

			presets, _ := dbInstance.GetPresets()

			expected := `{"name":"storedPreset","video":{},"audio":{}}`
			Expect(rb2).To(Equal(expected))
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(len(presets)).To(Equal(1))
		})

		Context("when the request is invalid", func() {
			It("should return BadRequest if preset is malformed", func() {
				preset := `{"neime: "badPreset}}`
				request, _ := http.NewRequest(http.MethodPost,
					"http://server/presets",
					strings.NewReader(preset))
				response, err := httpClient.Do(request)
				Expect(err).NotTo(HaveOccurred())

				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("UpdatePreset", func() {
		It("should update an existing preset", func() {
			dbInstance.StorePreset(types.Preset{Name: "examplePreset"})
			preset := `{"name":"examplePreset","Description": "new description","video": {},"audio": {}}`

			request, _ := http.NewRequest(http.MethodPut,
				"http://server/presets",
				strings.NewReader(preset))
			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())

			presets, _ := dbInstance.GetPresets()
			newPreset := presets[0]
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(newPreset.Description).To(Equal("new description"))
		})

		Context("when the request is invalid", func() {
			It("should return BadRequest if preset is malformed", func() {
				dbInstance.StorePreset(types.Preset{Name: "examplePreset"})
				preset := `{"name":"examplePreset","Description: "new description","video": {},"audio": {}}`

				request, _ := http.NewRequest(http.MethodPut,
					"http://server/presets",
					strings.NewReader(preset))
				response, err := httpClient.Do(request)
				Expect(err).NotTo(HaveOccurred())

				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})

			It("should return BadRequest if preset don't exists", func() {
				preset := `{"name":"dont-exists","Description": "new description","video": {},"audio": {}}`

				request, _ := http.NewRequest(http.MethodPut,
					"http://server/presets",
					strings.NewReader(preset))
				response, err := httpClient.Do(request)
				Expect(err).NotTo(HaveOccurred())
				rb, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				responseBody := string(rb)

				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
				expected := `{"error": "retrieving preset: preset not found"}`
				Expect(responseBody).To(Equal(expected))
			})
		})
	})

	Describe("DeletePreset", func() {
		It("should delete the preset", func() {
			examplePreset := types.Preset{
				Name:        "examplePreset",
				Description: "This is an example of preset",
				Container:   "mp4",
				RateControl: "vbr",
				Video: types.VideoPreset{
					Width:        "720",
					Height:       "1080",
					Codec:        "h264",
					Bitrate:      "10000",
					GopSize:      "90",
					GopMode:      "fixed",
					Profile:      "high",
					ProfileLevel: "3.1",

					InterlaceMode: "progressive",
				},
				Audio: types.AudioPreset{
					Codec:   "aac",
					Bitrate: "64000",
				},
			}
			dbInstance.StorePreset(examplePreset)

			request, _ := http.NewRequest(http.MethodDelete,
				"http://server/presets/examplePreset",
				nil)
			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())

			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		Context("when deleting the preset fails", func() {
			It("should return BadRequest", func() {
				request, _ := http.NewRequest(http.MethodDelete,
					"http://server/presets/yoyoyo",
					nil)
				response, err := httpClient.Do(request)
				Expect(err).NotTo(HaveOccurred())
				rb, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				responseBody := string(rb)

				expected := `{"error": "deleting preset: preset not found"}`
				Expect(responseBody).To(Equal(expected))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})
})
