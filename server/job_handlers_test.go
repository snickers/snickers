package server_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"

	"code.cloudfoundry.org/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/db/dbfakes"
	"github.com/snickers/snickers/server"
	"github.com/snickers/snickers/types"
)

var _ = Describe("Job Handlers", func() {

	var (
		client      *http.Client
		fakeStorage *dbfakes.FakeStorage
		logger      *lagertest.TestLogger
		testServer  *httptest.Server
		err         error

		socketPath string
		tmpDir     string
	)

	BeforeEach(func() {
		tmpDir, err = ioutil.TempDir(os.TempDir(), "job-handlers")
		socketPath = path.Join(tmpDir, "snickers.sock")
		logger = lagertest.NewTestLogger("snickers-test")
		fakeStorage = new(dbfakes.FakeStorage)
		snickersServer := server.New(logger, "unix", socketPath, fakeStorage)
		testServer = httptest.NewServer(snickersServer.Handler())

		client = &http.Client{
			Transport: &http.Transport{},
		}
	})

	AfterEach(func() {
		testServer.Close()
		os.RemoveAll(tmpDir)
	})

	createJob := func(body io.Reader) *http.Response {
		req, err := http.NewRequest(http.MethodPost, testServer.URL+"/jobs", body)
		Expect(err).NotTo(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())

		return resp
	}

	listJobs := func() *http.Response {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/jobs", nil)
		Expect(err).NotTo(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())

		return resp
	}

	jobDetails := func(id string) *http.Response {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/jobs/"+id, nil)
		Expect(err).NotTo(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())

		return resp
	}

	Describe("CreateJob", func() {
		var createJobResp *http.Response

		JustBeforeEach(func() {
			createJobResp = createJob(bytes.NewBufferString(`{}`))
		})

		It("creates a new job", func() {
			Expect(fakeStorage.StoreJobCallCount()).To(Equal(1))
			Expect(createJobResp.StatusCode).To(Equal(http.StatusCreated))
		})

		Context("when fails to retrieve the preset", func() {
			BeforeEach(func() {
				fakeStorage.RetrievePresetReturns(types.Preset{}, errors.New("Boom!"))
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(createJobResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(createJobResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(MatchJSON(`{"error": "retrieving preset: Boom!"}`))
			})
		})

		Context("when fails to store the job", func() {
			BeforeEach(func() {
				fakeStorage.StoreJobReturns(types.Job{}, errors.New("Boom!"))
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(createJobResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(createJobResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(MatchJSON(`{"error": "storing job: Boom!"}`))
			})
		})
	})

	Describe("ListJobs", func() {
		var listJobsResp *http.Response

		JustBeforeEach(func() {
			listJobsResp = listJobs()
		})

		It("returns a list of jobs", func() {
			Expect(fakeStorage.GetJobsCallCount()).To(Equal(1))
			Expect(listJobsResp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("when fails to get the jobs", func() {
			BeforeEach(func() {
				fakeStorage.GetJobsReturns(nil, errors.New("Boom!"))
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(listJobsResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(listJobsResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(MatchJSON(`{"error": "getting jobs: Boom!"}`))
			})
		})
	})

	Describe("GetJobDetails", func() {
		var (
			job            types.Job
			jobDetailsResp *http.Response
		)

		BeforeEach(func() {
			job = types.Job{
				ID:               "id",
				Source:           "source",
				Destination:      "destination",
				Preset:           types.Preset{},
				Status:           types.JobCreated,
				Details:          "details",
				LocalSource:      "source",
				LocalDestination: "destionation",
			}
		})

		BeforeEach(func() {
			fakeStorage.RetrieveJobReturns(job, nil)
		})

		JustBeforeEach(func() {
			jobDetailsResp = jobDetails(job.ID)
		})

		It("returns job details", func() {
			body, err := ioutil.ReadAll(jobDetailsResp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobDetailsResp.StatusCode).To(Equal(http.StatusOK))

			jsonJob, err := json.Marshal(job)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(MatchJSON(jsonJob))
		})

		Context("when fails to get the job", func() {
			BeforeEach(func() {
				fakeStorage.RetrieveJobReturns(types.Job{}, errors.New("Boom!"))
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(jobDetailsResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(jobDetailsResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(body)).To(Equal(`{"error": "retrieving job: Boom!"}`))
			})
		})
	})

	Describe("StartJob", func() {
		PIt("starts a job", func() {
			//TODO: Need to create an abstraction for this.
		})
	})
})
