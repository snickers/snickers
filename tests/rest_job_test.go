package snickers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/rest"
	"github.com/snickers/snickers/types"
)

var _ = Describe("Rest API", func() {
	var (
		response   *httptest.ResponseRecorder
		server     *mux.Router
		dbInstance db.DatabaseInterface
	)

	BeforeEach(func() {
		response = httptest.NewRecorder()
		server = rest.NewRouter()
		dbInstance, _ = db.GetDatabase()
		dbInstance.ClearDatabase()
	})

	Describe("GET /jobs", func() {
		It("should return application/json on its content type", func() {
			request, _ := http.NewRequest("GET", "/jobs", nil)
			server.ServeHTTP(response, request)
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		It("should return stored jobs", func() {
			exampleJob1 := types.Job{ID: "123"}
			exampleJob2 := types.Job{ID: "321"}
			dbInstance.StoreJob(exampleJob1)
			dbInstance.StoreJob(exampleJob2)

			expected1, _ := json.Marshal(`[{"id":"123","source":"","destination":"","preset":{"video":{},"audio":{}},"status":"","progress":""},{"id":"321","source":"","destination":"","preset":{"video":{},"audio":{}},"status":"","progress":""}]`)
			expected2, _ := json.Marshal(`[{"id":"321","source":"","destination":"","preset":{"video":{},"audio":{}},"status":"","progress":""},{"id":"123","source":"","destination":"","preset":{"video":{},"audio":{}},"status":"","progress":""}]`)

			request, _ := http.NewRequest("GET", "/jobs", nil)
			server.ServeHTTP(response, request)
			responseBody, _ := json.Marshal(response.Body.String())

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(responseBody).To(SatisfyAny(Equal(expected1), Equal(expected2)))
		})
	})

	Describe("POST /jobs", func() {
		It("should create a new job", func() {
			dbInstance.StorePreset(types.Preset{Name: "presetName"})
			jobJSON := []byte(`{"source": "http://flv.io/src.mp4", "destination": "s3://l@p:google.com", "preset": "presetName"}`)
			request, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jobJSON))
			server.ServeHTTP(response, request)
			responseBody := response.Body.String()

			jobs, _ := dbInstance.GetJobs()
			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
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
				jobJSON := []byte(`{"source": "http://flv.io/src.mp4", "destination": "s3://l@p:google.com", "preset": "presetName"}`)
				request, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jobJSON))
				server.ServeHTTP(response, request)
				responseBody, _ := json.Marshal(response.Body.String())

				expected, _ := json.Marshal(`{"error": "retrieving preset: preset not found"}`)
				Expect(responseBody).To(Equal(expected))
				Expect(response.Code).To(Equal(http.StatusBadRequest))
				Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			})

			It("should return BadRequest if preset is malformed", func() {
				jobJSON := []byte(`{"source: "http://flv.io/src.mp4", "destinat}`)
				request, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jobJSON))
				server.ServeHTTP(response, request)
				responseBody, _ := json.Marshal(response.Body.String())

				expected, _ := json.Marshal(`{"error": "unpacking job: invalid character 'h' after object key"}`)
				Expect(responseBody).To(Equal(expected))
				Expect(response.Code).To(Equal(http.StatusBadRequest))
				Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("GET /job/:id", func() {
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

			request, _ := http.NewRequest("GET", "/jobs/123-123-123", nil)
			server.ServeHTTP(response, request)
			Expect(response.Body.String()).To(Equal(string(expected)))
			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		Context("when getting the job fails", func() {
			It("should return BadRequest", func() {
				request, _ := http.NewRequest("GET", "/jobs/123-456-9292", nil)
				server.ServeHTTP(response, request)
				expected, _ := json.Marshal(`{"error": "retrieving job: job not found"}`)
				responseBody, _ := json.Marshal(response.Body.String())
				Expect(responseBody).To(Equal(expected))
				Expect(response.Code).To(Equal(http.StatusBadRequest))
				Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})

	Describe("POST /jobs/:id/start", func() {
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

			request, _ := http.NewRequest("POST", "/jobs/123-123-123/start", nil)
			server.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		Context("when starting a job fails", func() {
			It("should return BadRequest", func() {
				request, _ := http.NewRequest("POST", "/jobs/123-456-9292/start", nil)
				server.ServeHTTP(response, request)
				expected, _ := json.Marshal(`{"error": "retrieving job: job not found"}`)
				responseBody, _ := json.Marshal(response.Body.String())
				Expect(responseBody).To(Equal(expected))
				Expect(response.Code).To(Equal(http.StatusBadRequest))
				Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			})
		})
	})
})
