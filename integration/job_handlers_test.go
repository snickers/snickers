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
	"github.com/snickers/snickers/db/memory"
	"github.com/snickers/snickers/server"
	"github.com/snickers/snickers/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Job Handlers", func() {
	var (
		httpClient     *http.Client
		snickersServer *server.SnickersServer
		err            error
		log            *lagertest.TestLogger
		tmpDir         string
		dbInstance     db.Storage
	)

	BeforeEach(func() {
		log = lagertest.NewTestLogger("snickers-test")
		tmpDir, err = ioutil.TempDir(os.TempDir(), "snickers-server-test")
		socketPath := path.Join(tmpDir, "snickers.sock")

		dbInstance, _ = memory.GetDatabase()
		currentDir, _ := os.Getwd()
		configPath := currentDir + "/../fixtures/config.json"
		snickersServer = server.New(log, configPath, "unix", socketPath, dbInstance)
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

	Describe("CreateJob", func() {
		var sendRequest = func(body string) *http.Response {
			request, _ := http.NewRequest(http.MethodPost,
				"http://server/jobs",
				strings.NewReader(body))

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			return response
		}

		It("creates a new job", func() {
			dbInstance.StorePreset(types.Preset{Name: "presetName"})

			response := sendRequest(`{"source": "http://flv.io/src.mp4", "destination": "s3://l@p:google.com", "preset": "presetName"}`)

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).To(MatchRegexp(`{"id":".*","source":"http://flv.io/src.mp4","destination":"s3://l@p:google.com","preset":{"name":"presetName","video":{},"audio":{}},"status":"created","progress":""}`))
			Expect(response.StatusCode).To(Equal(http.StatusCreated))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})

		Context("when the request is invalid", func() {
			It("should return BadRequest if preset is not set", func() {
				response := sendRequest(`{"source": "http://flv.io/src.mp4", "destination": "s3://l@p:google.com", "preset": "presetName"}`)

				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
					"error": "retrieving preset: preset not found"
				}`))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})

			It("should return BadRequest if preset is malformed", func() {
				response := sendRequest(`{"source: "http://flv.io/src.mp4", "destinat}`)

				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
				"error": "unpacking job: invalid character 'h' after object key"
				}`))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("ListJobs", func() {
		var sendRequest = func() *http.Response {
			request, _ := http.NewRequest(http.MethodGet,
				"http://server/jobs",
				nil)

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			return response
		}

		It("returns stored jobs", func() {
			job := types.Job{ID: "123"}
			dbInstance.StoreJob(job)

			response := sendRequest()
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).To(MatchJSON(`[{
				"destination": "",
				"id": "123",
				"preset": {
						"audio": {},
						"video": {}
				},
				"progress": "",
				"source": "",
				"status": ""
			}]`))
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})
	})

	Describe("GetJobDetails", func() {
		var sendRequest = func(jobID string) *http.Response {
			request, _ := http.NewRequest(http.MethodGet,
				fmt.Sprintf("http://server/jobs/%s", jobID),
				nil)

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			return response
		}

		It("should return the job with details", func() {
			job := types.Job{
				ID:          "123-123-123",
				Source:      "http://source.here.mp4",
				Destination: "s3://ae@ae.com",
				Preset:      types.Preset{},
				Status:      types.JobCreated,
				Details:     "0%",
			}
			dbInstance.StoreJob(job)
			expected, _ := json.Marshal(&job)

			response := sendRequest("123-123-123")
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).To(Equal(expected))
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})

		Context("when getting the job fails", func() {
			It("should return BadRequest", func() {
				response := sendRequest("123-456-9292")
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
					"error": "retrieving job: job not found"
				}`))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("StartJob", func() {
		var sendRequest = func(jobID string) *http.Response {
			request, _ := http.NewRequest(http.MethodPost,
				fmt.Sprintf("http://server/jobs/%s/start", jobID),
				nil)

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			return response
		}

		It("should return status OK", func() {
			job := types.Job{
				ID:          "123-123-123",
				Source:      "http://source.here.mp4",
				Destination: "s3://ae@ae.com",
				Preset:      types.Preset{},
				Status:      types.JobCreated,
				Details:     "0%",
			}
			dbInstance.StoreJob(job)

			response := sendRequest("123-123-123")
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})

		Context("when starting a job fails", func() {
			It("should return BadRequest", func() {
				response := sendRequest("123-456-9292")
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
					"error": "retrieving job: job not found"
				}`))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})
})
