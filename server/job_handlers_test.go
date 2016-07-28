package server_test

import (
	"encoding/json"
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
	"github.com/snickers/snickers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobHandlers", func() {
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

	Describe("CreateJob", func() {
		It("should create a new job", func() {
			dbInstance.StorePreset(snickers.Preset{Name: "presetName"})
			jobJSON := `{"source": "http://flv.io/src.mp4", "destination": "s3://l@p:google.com", "preset": "presetName"}`
			request, _ := http.NewRequest(http.MethodPost,
				"http://server/jobs",
				strings.NewReader(jobJSON))

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())

			rb, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())
			responseBody := string(rb)

			jobs, _ := dbInstance.GetJobs()
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			Expect(len(jobs)).To(Equal(1))
			Expect(strings.Contains(responseBody, "id")).To(BeTrue())
			Expect(strings.Contains(responseBody, `"status":"created"`)).To(BeTrue())
			Expect(strings.Contains(responseBody, `"progress":""`)).To(BeTrue())

			job := jobs[0]
			Expect(job.Source).To(Equal("http://flv.io/src.mp4"))
			Expect(job.Destination).To(Equal("s3://l@p:google.com"))
			Expect(job.Preset.Name).To(Equal("presetName"))
		})

		Context("when the request is invalid", func() {
			It("should return BadRequest if preset is not set", func() {
				jobJSON := `{"source": "http://flv.io/src.mp4", "destination": "s3://l@p:google.com", "preset": "presetName"}`
				request, err := http.NewRequest(http.MethodPost,
					"http://server/jobs",
					strings.NewReader(jobJSON))
				Expect(err).NotTo(HaveOccurred())

				response, err := httpClient.Do(request)
				Expect(err).NotTo(HaveOccurred())
				rb, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				responseBody := string(rb)

				expected := `{"error": "retrieving preset: preset not found"}`
				Expect(responseBody).To(Equal(expected))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})

			It("should return BadRequest if preset is malformed", func() {
				jobJSON := `{"source: "http://flv.io/src.mp4", "destinat}`
				request, err := http.NewRequest(http.MethodPost,
					"http://server/jobs",
					strings.NewReader(jobJSON))
				Expect(err).NotTo(HaveOccurred())

				response, err := httpClient.Do(request)
				Expect(err).NotTo(HaveOccurred())
				rb, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				responseBody := string(rb)

				expected := `{"error": "unpacking job: invalid character 'h' after object key"}`
				Expect(responseBody).To(Equal(expected))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("ListJobs", func() {
		It("should return application/json on its content type", func() {
			request, err := http.NewRequest(http.MethodGet,
				"http://server/jobs",
				nil)
			Expect(err).NotTo(HaveOccurred())
			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})

		It("should return stored jobs", func() {
			exampleJob1 := snickers.Job{ID: "123"}
			exampleJob2 := snickers.Job{ID: "321"}
			dbInstance.StoreJob(exampleJob1)
			dbInstance.StoreJob(exampleJob2)

			expected1, _ := json.Marshal(`[{"id":"123","source":"","destination":"","preset":{"video":{},"audio":{}},"status":"","progress":""},{"id":"321","source":"","destination":"","preset":{"video":{},"audio":{}},"status":"","progress":""}]`)
			expected2, _ := json.Marshal(`[{"id":"321","source":"","destination":"","preset":{"video":{},"audio":{}},"status":"","progress":""},{"id":"123","source":"","destination":"","preset":{"video":{},"audio":{}},"status":"","progress":""}]`)

			request, _ := http.NewRequest(http.MethodGet,
				"http://server/jobs",
				nil)

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			rb, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())
			responseBody, _ := json.Marshal(string(rb))

			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(responseBody).To(SatisfyAny(Equal(expected1), Equal(expected2)))
		})
	})

	Describe("GetJobDetails", func() {
		It("should return the job with details", func() {
			job := snickers.Job{
				ID:          "123-123-123",
				Source:      "http://source.here.mp4",
				Destination: "s3://ae@ae.com",
				Preset:      snickers.Preset{},
				Status:      snickers.JobCreated,
				Details:     "0%",
			}
			dbInstance.StoreJob(job)
			expected, _ := json.Marshal(&job)

			request, _ := http.NewRequest(http.MethodGet,
				"http://server/jobs/123-123-123",
				nil)

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			responseBody, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(responseBody).To(Equal(expected))
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})

		Context("when getting the job fails", func() {
			It("should return BadRequest", func() {
				request, _ := http.NewRequest(http.MethodGet,
					"http://server/jobs/123-456-9292",
					nil)

				response, err := httpClient.Do(request)
				Expect(err).NotTo(HaveOccurred())
				rb, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				expected := `{"error": "retrieving job: job not found"}`
				Expect(string(rb)).To(Equal(expected))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("StartJob", func() {
		It("should return status OK", func() {
			job := snickers.Job{
				ID:          "123-123-123",
				Source:      "http://source.here.mp4",
				Destination: "s3://ae@ae.com",
				Preset:      snickers.Preset{},
				Status:      snickers.JobCreated,
				Details:     "0%",
			}
			dbInstance.StoreJob(job)

			request, _ := http.NewRequest(http.MethodPost,
				"http://server/jobs/123-123-123/start",
				nil)
			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
		})

		Context("when starting a job fails", func() {
			It("should return BadRequest", func() {
				request, _ := http.NewRequest(http.MethodPost,
					"http://server/jobs/123-456-9292/start",
					nil)
				response, err := httpClient.Do(request)
				Expect(err).NotTo(HaveOccurred())
				expected, _ := json.Marshal(`{"error": "retrieving job: job not found"}`)
				rb, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				responseBody, _ := json.Marshal(string(rb))
				Expect(responseBody).To(Equal(expected))
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})
})
