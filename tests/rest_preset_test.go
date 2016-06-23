package snickers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/rest"
	"github.com/flavioribeiro/snickers/types"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rest API", func() {
	Context("/presets location", func() {
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

		It("GET should return application/json on its content type", func() {
			request, _ := http.NewRequest("GET", "/presets", nil)
			server.ServeHTTP(response, request)
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		It("GET should return stored presets", func() {
			Skip("we should sort the arrays before compare")
			examplePreset1 := types.Preset{Name: "a"}
			examplePreset2 := types.Preset{Name: "b"}
			dbInstance.StorePreset(examplePreset1)
			dbInstance.StorePreset(examplePreset2)

			expected, _ := json.Marshal(`[{"name":"a","video":{},"audio":{}},{"name":"b","video":{},"audio":{}}]`)

			request, _ := http.NewRequest("GET", "/presets", nil)
			server.ServeHTTP(response, request)
			responseBody, _ := json.Marshal(response.Body.String())

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(responseBody).To(Equal(expected))
		})

		It("POST should save a new preset", func() {
			preset := []byte(`{"name": "storedPreset", "video": {},"audio": {}}`)
			request, _ := http.NewRequest("POST", "/presets", bytes.NewBuffer(preset))
			server.ServeHTTP(response, request)

			presets, _ := dbInstance.GetPresets()
			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			Expect(len(presets)).To(Equal(1))
		})

		It("POST with malformed preset should return bad request", func() {
			preset := []byte(`{"neime: "badPreset}}`)
			request, _ := http.NewRequest("POST", "/presets", bytes.NewBuffer(preset))
			server.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		It("PUT with a new preset should update the preset", func() {
			dbInstance.StorePreset(types.Preset{Name: "examplePreset"})
			preset := []byte(`{"name":"examplePreset","Description": "new description","video": {},"audio": {}}`)

			request, _ := http.NewRequest("PUT", "/presets", bytes.NewBuffer(preset))
			server.ServeHTTP(response, request)

			presets, _ := dbInstance.GetPresets()
			newPreset := presets[0]
			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			Expect(newPreset.Description).To(Equal("new description"))
		})

		It("PUT with malformed preset should return bad request", func() {
			dbInstance.StorePreset(types.Preset{Name: "examplePreset"})
			preset := []byte(`{"name":"examplePreset","Description: "new description","video": {},"audio": {}}`)

			request, _ := http.NewRequest("PUT", "/presets", bytes.NewBuffer(preset))
			server.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		It("GET for a given preset should return preset details", func() {
			examplePreset := types.Preset{
				Name:         "examplePreset",
				Description:  "This is an example of preset",
				Container:    "mp4",
				Profile:      "high",
				ProfileLevel: "3.1",
				RateControl:  "VBR",
				Video: types.VideoPreset{
					Width:         "720",
					Height:        "1080",
					Codec:         "h264",
					Bitrate:       "10000",
					GopSize:       "90",
					GopMode:       "fixed",
					InterlaceMode: "progressive",
				},
				Audio: types.AudioPreset{
					Codec:   "aac",
					Bitrate: "64000",
				},
			}
			dbInstance.StorePreset(examplePreset)
			expected, _ := json.Marshal(examplePreset)

			request, _ := http.NewRequest("GET", "/presets/examplePreset", nil)
			server.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			Expect(response.Body.String()).To(Equal(string(expected)))
		})
	})
})
