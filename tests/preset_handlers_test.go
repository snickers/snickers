package snickers_test

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

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/server"
	"github.com/snickers/snickers/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Preset Handlers", func() {
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

		Expect(snickersServer.Start(false)).NotTo(HaveOccurred())
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
		var sendRequest = func() *http.Response {
			request, err := http.NewRequest(http.MethodGet,
				"http://server/presets",
				nil)
			Expect(err).NotTo(HaveOccurred())

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			return response
		}

		It("returns application/json on its content type", func() {
			response := sendRequest()
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})

		It("returns stored presets", func() {
			preset := types.Preset{Name: "a"}
			dbInstance.StorePreset(preset)

			response := sendRequest()
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(body).To(MatchJSON(`[{
						"audio": {},
						"name": "a",
						"video": {}
	  	}]`))
		})
	})

	Describe("GetPresetDetails", func() {
		var sendRequest = func(presetName string) *http.Response {
			request, err := http.NewRequest(http.MethodGet,
				fmt.Sprintf("http://server/presets/%s", presetName),
				nil)
			Expect(err).NotTo(HaveOccurred())

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			return response
		}

		It("should return the preset with details", func() {
			preset := types.Preset{
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
			dbInstance.StorePreset(preset)
			expected, _ := json.Marshal(&preset)

			response := sendRequest(preset.Name)
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).To(MatchJSON(expected))
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})

		Context("when getting the preset fails", func() {
			It("should return BadRequest", func() {
				response := sendRequest("yoyoyo")
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
					"error": "retrieving preset: preset not found"
				}`))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("CreatePreset", func() {
		var sendRequest = func(payload string) *http.Response {
			request, err := http.NewRequest(http.MethodPost,
				"http://server/presets",
				strings.NewReader(payload))
			Expect(err).NotTo(HaveOccurred())

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			return response
		}

		It("should save a new preset", func() {
			response := sendRequest(`{"name": "storedPreset", "video": {},"audio": {}}`)

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).To(MatchJSON(`{
				"name":"storedPreset",
				"video":{},
				"audio":{}
			}`))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(response.StatusCode).To(Equal(http.StatusCreated))
		})

		Context("when the request is invalid", func() {
			It("should return BadRequest if preset is malformed", func() {
				response := sendRequest(`{"neime: "badPreset}}`)

				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("UpdatePreset", func() {
		var sendRequest = func(payload string) *http.Response {
			request, err := http.NewRequest(http.MethodPut,
				"http://server/presets",
				strings.NewReader(payload))
			Expect(err).NotTo(HaveOccurred())

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			return response
		}

		It("updates an existing preset", func() {
			dbInstance.StorePreset(types.Preset{Name: "examplePreset"})
			response := sendRequest(`{"name":"examplePreset","Description": "new description","video": {},"audio": {}}`)

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).To(MatchJSON(`{
				"name":"examplePreset",
				"description":"new description",
				"video":{},
				"audio":{}
			}`))
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})

		Context("when the request is invalid", func() {
			It("should return BadRequest if preset is malformed", func() {
				response := sendRequest(`{"name":"examplePreset","Description: "new description","video": {},"audio": {}}`)

				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})

			It("should return BadRequest if preset don't exists", func() {
				response := sendRequest(`{"name":"dont-exists","Description": "new description","video": {},"audio": {}}`)

				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
					"error": "retrieving preset: preset not found"
				}`))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("DeletePreset", func() {
		var sendRequest = func(presetName string) *http.Response {
			request, err := http.NewRequest(http.MethodDelete,
				fmt.Sprintf("http://server/presets/%s", presetName),
				nil)
			Expect(err).NotTo(HaveOccurred())

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			return response
		}

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

			response := sendRequest("examplePreset")
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		Context("when deleting the preset fails", func() {
			It("should return BadRequest", func() {
				response := sendRequest("yoyoyo")
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
					"error": "deleting preset: preset not found"
				}`))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})
})
