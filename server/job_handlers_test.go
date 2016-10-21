package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

var _ = Describe("Job handler", func() {
	var (
		dbInstance       db.Storage
		cfg              gonfig.Gonfig
		sn               *SnickersServer
		input            types.JobInput
		jobRecorder      *httptest.ResponseRecorder
		respJobInputBody map[string]interface{}
	)

	BeforeEach(func() {
		currentDir, _ := os.Getwd()
		cfg, _ = gonfig.FromJsonFile(currentDir + "/../fixtures/config.json")
		dbInstance, _ = db.GetDatabase(cfg)
		dbInstance.ClearDatabase()
		sn = New(lagertest.NewTestLogger("job-handler"), cfg, "tcp", ":8000", dbInstance)

		preset := types.Preset{
			Name: "mp4_1080p",
		}
		presetRecorder := httptest.NewRecorder()
		presetData, _ := json.Marshal(preset)
		reqPreset, _ := http.NewRequest(http.MethodPost, "/presets", bytes.NewReader(presetData))
		sn.Handler().ServeHTTP(presetRecorder, reqPreset)

		input = types.JobInput{
			Source:      "http://s3.example.com/videos/video1.mov",
			Destination: "s3://example-bucket/future/video1.mp4",
			PresetName:  "mp4_1080p",
		}

		jobRecorder = httptest.NewRecorder()
		payloadData, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(payloadData))

		sn.Handler().ServeHTTP(jobRecorder, req)
		json.Unmarshal(jobRecorder.Body.Bytes(), &respJobInputBody)
	})

	Context("Create job", func() {
		It("should create a job in the db instance", func() {
			Expect(jobRecorder.Code).To(BeIdenticalTo(http.StatusCreated))
			jobID, ok := respJobInputBody["id"].(string)
			Expect(ok).To(BeIdenticalTo(true))
			job, err := dbInstance.RetrieveJob(jobID)
			Expect(err).NotTo(HaveOccurred())
			Expect(job.Destination).To(BeIdenticalTo(input.Destination))
		})

		It("should get a given job details", func() {
			var jobBody map[string]interface{}

			recorder := httptest.NewRecorder()
			reqJobDetail, _ := http.NewRequest(http.MethodGet, "/jobs/"+respJobInputBody["id"].(string), nil)
			sn.Handler().ServeHTTP(recorder, reqJobDetail)
			json.Unmarshal(recorder.Body.Bytes(), &jobBody)
			Expect(recorder.Code).To(BeIdenticalTo(http.StatusOK))
			Expect(jobBody["id"]).To(BeIdenticalTo(respJobInputBody["id"]))
			Expect(jobBody["status"]).To(BeIdenticalTo(respJobInputBody["status"]))
		})

		It("should list all jobs", func() {
			secondInput := types.JobInput{
				Source:      "http://s3.example.com/videos/video2.mov",
				Destination: "s3://example-bucket/future/video2.mp4",
				PresetName:  "mp4_1080p",
			}

			recorder := httptest.NewRecorder()
			data, _ := json.Marshal(secondInput)
			req, _ := http.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(data))
			sn.Handler().ServeHTTP(recorder, req)

			listRecorder := httptest.NewRecorder()
			var jobListBody []map[string]interface{}

			listJobsRequest, _ := http.NewRequest(http.MethodGet, "/jobs", nil)
			sn.Handler().ServeHTTP(listRecorder, listJobsRequest)
			json.Unmarshal(listRecorder.Body.Bytes(), &jobListBody)
			Expect(listRecorder.Code).To(BeIdenticalTo(http.StatusOK))
			Expect(len(jobListBody)).To(Equal(2))
		})
	})
})
