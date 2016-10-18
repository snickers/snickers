package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"

	"code.cloudfoundry.org/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/server"
	"github.com/snickers/snickers/types"
)

var _ = Describe("Preset Handlers", func() {

	var (
		client     *http.Client
		dbInstance db.Storage
		logger     *lagertest.TestLogger
		testServer *httptest.Server
		err        error

		socketPath string
		tmpDir     string
	)

	BeforeEach(func() {
		currentDir, _ := os.Getwd()
		cfg, _ := gonfig.FromJsonFile(currentDir + "/../fixtures/config.json")
		tmpDir, err = ioutil.TempDir(os.TempDir(), "preset-handlers")
		socketPath = path.Join(tmpDir, "snickers.sock")
		logger = lagertest.NewTestLogger("preset-handlers")

		dbInstance, _ = db.GetDatabase(cfg)
		snickersServer := server.New(logger, cfg, "unix", socketPath, dbInstance)
		testServer = httptest.NewServer(snickersServer.Handler())

		client = &http.Client{
			Transport: &http.Transport{},
		}
	})

	AfterEach(func() {
		testServer.Close()
		os.RemoveAll(tmpDir)
	})

	createPreset := func(body io.Reader) *http.Response {
		req, err := http.NewRequest(http.MethodPost, testServer.URL+"/presets", body)
		Expect(err).NotTo(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())

		return resp
	}

	updatePreset := func(body io.Reader) *http.Response {
		req, err := http.NewRequest(http.MethodPut, testServer.URL+"/presets", body)
		Expect(err).NotTo(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())

		return resp
	}

	listPresets := func() *http.Response {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/presets", nil)
		Expect(err).NotTo(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())

		return resp
	}

	presetDetails := func(name string) *http.Response {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/presets/"+name, nil)
		Expect(err).NotTo(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())

		return resp
	}

	deletePreset := func(name string) *http.Response {
		req, err := http.NewRequest(http.MethodDelete, testServer.URL+"/presets/"+name, nil)
		Expect(err).NotTo(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())

		return resp
	}

	Describe("CreatePreset", func() {
		var (
			preset           io.Reader
			createPresetResp *http.Response
		)

		It("creates a new preset", func() {
			preset = bytes.NewBufferString(`{"name":"foobar"}`)
			createPresetResp = createPreset(preset)
			Expect(createPresetResp.StatusCode).To(Equal(http.StatusCreated))
		})

	})

	Describe("UpdatePreset", func() {
		var (
			preset           io.Reader
			updatePresetResp *http.Response
		)

		It("updates an existing new preset based on its name", func() {
			preset = bytes.NewBufferString(`{"name":"foobar"}`)
			updatePresetResp = updatePreset(preset)
			Expect(updatePresetResp.StatusCode).To(Equal(http.StatusOK))
		})

	})

	Describe("ListPresets", func() {
		var listPresetsResp *http.Response

		It("returns a list of presets", func() {
			listPresetsResp = listPresets()
			Expect(listPresetsResp.StatusCode).To(Equal(http.StatusOK))
		})

	})

	Describe("GetPresetDetails", func() {
		var (
			preset            types.Preset
			presetDetailsResp *http.Response
		)

		It("returns preset details", func() {
			preset = types.Preset{
				Name: "foobar",
			}
			presetDetailsResp = presetDetails(preset.Name)
			body, err := ioutil.ReadAll(presetDetailsResp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(presetDetailsResp.StatusCode).To(Equal(http.StatusOK))

			jsonPreset, err := json.Marshal(preset)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(MatchJSON(jsonPreset))
		})
	})

	Describe("DeletePreset", func() {
		It("deletes the preset", func() {
			presetDetailsResp := deletePreset("foobar")
			Expect(presetDetailsResp.StatusCode).To(Equal(http.StatusOK))
		})
	})

})
