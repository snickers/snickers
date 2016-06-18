package snickers_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/db/memory"
	"github.com/flavioribeiro/snickers/rest"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rest API", func() {
	Context("/presets location", func() {
		var (
			response   *httptest.ResponseRecorder
			server     *mux.Router
			dbInstance *memory.Database
		)

		BeforeEach(func() {
			response = httptest.NewRecorder()
			server = rest.NewRouter()
			dbInstance, _ = memory.GetDatabase()
			memory.ClearDatabase()

		})

		It("GET should return stored presets", func() {
			examplePreset := db.Preset{
				Name: "examplePreset",
			}
			dbInstance.StorePreset(examplePreset)
			expected := `[{"name":"examplePreset","video":{},"audio":{}}]`

			request, _ := http.NewRequest("GET", "/presets", nil)
			server.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(200))
			Expect(string(response.Body.String())).To(Equal(expected))
		})
	})
})
