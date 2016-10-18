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
		dbInstance db.Storage
		cfg        gonfig.Gonfig
		sn         *SnickersServer
	)

	BeforeEach(func() {
		currentDir, _ := os.Getwd()
		cfg, _ = gonfig.FromJsonFile(currentDir + "/../fixtures/config.json")
		dbInstance, _ = db.GetDatabase(cfg)
		dbInstance.ClearDatabase()
		sn = New(lagertest.NewTestLogger("job-handler"), cfg, "tcp", ":8000", dbInstance)
	})

	Context("Create job", func() {
		It("should create a job in the db instance", func() {
			input := types.JobInput{
				Source:      "http://s3.example.com/videos/video1.mov",
				Destination: "s3://example-bucket/future/video1.mp4",
				PresetName:  "mp4_1080p",
			}
			preset := types.Preset{
				Name: "mp4_1080p",
			}
			presetRecorder := httptest.NewRecorder()
			presetData, _ := json.Marshal(preset)
			reqPreset, _ := http.NewRequest(http.MethodPost, "/presets", bytes.NewReader(presetData))
			sn.Handler().ServeHTTP(presetRecorder, reqPreset)

			recorder := httptest.NewRecorder()
			payloadData, _ := json.Marshal(input)
			req, _ := http.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(payloadData))

			sn.Handler().ServeHTTP(recorder, req)
			Expect(recorder.Code).To(BeIdenticalTo(http.StatusCreated))
			var respBody map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &respBody)
			Expect(err).NotTo(HaveOccurred())
			jobID, ok := respBody["id"].(string)

			Expect(ok).To(BeIdenticalTo(true))
			job, err := dbInstance.RetrieveJob(jobID)
			Expect(err).NotTo(HaveOccurred())
			Expect(job.Destination).To(BeIdenticalTo(input.Destination))
		})
	})
})
