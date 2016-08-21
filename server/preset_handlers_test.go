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

var _ = Describe("Preset Handlers", func() {

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
		tmpDir, err = ioutil.TempDir(os.TempDir(), "preset-handlers")
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

		BeforeEach(func() {
			preset = bytes.NewBufferString(`{"name":"foobar"}`)
		})

		JustBeforeEach(func() {
			createPresetResp = createPreset(preset)
		})

		It("creates a new preset", func() {
			Expect(fakeStorage.StorePresetCallCount()).To(Equal(1))
			Expect(createPresetResp.StatusCode).To(Equal(http.StatusCreated))
		})

		Context("when fails to parse the preset", func() {
			BeforeEach(func() {
				preset = bytes.NewBufferString("invalid")
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(createPresetResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(createPresetResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(ContainSubstring("unpacking preset: invalid character"))
			})
		})

		Context("when fails to store the preset", func() {
			BeforeEach(func() {
				fakeStorage.StorePresetReturns(types.Preset{}, errors.New("Boom!"))
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(createPresetResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(createPresetResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(MatchJSON(`{"error": "storing preset: Boom!"}`))
			})
		})
	})

	Describe("UpdatePreset", func() {
		var (
			preset           io.Reader
			updatePresetResp *http.Response
		)

		BeforeEach(func() {
			preset = bytes.NewBufferString(`{"name":"foobar"}`)
		})

		JustBeforeEach(func() {
			updatePresetResp = updatePreset(preset)
		})

		It("updates an existing new preset based on its name", func() {
			presetName := fakeStorage.RetrievePresetArgsForCall(0)
			Expect(presetName).To(Equal("foobar"))
			Expect(updatePresetResp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("when fails to parse the preset", func() {
			BeforeEach(func() {
				preset = bytes.NewBufferString("invalid")
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(updatePresetResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(updatePresetResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(ContainSubstring("unpacking preset: invalid character"))
			})
		})

		Context("when fails to retrieve the preset", func() {
			BeforeEach(func() {
				fakeStorage.RetrievePresetReturns(types.Preset{}, errors.New("Boom!"))
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(updatePresetResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(updatePresetResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(MatchJSON(`{"error": "retrieving preset: Boom!"}`))
			})
		})

		Context("when fails to update the preset", func() {
			BeforeEach(func() {
				fakeStorage.UpdatePresetReturns(types.Preset{}, errors.New("Boom!"))
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(updatePresetResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(updatePresetResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(MatchJSON(`{"error": "updating preset: Boom!"}`))
			})
		})
	})

	Describe("ListPresets", func() {
		var listPresetsResp *http.Response

		JustBeforeEach(func() {
			listPresetsResp = listPresets()
		})

		It("returns a list of presets", func() {
			Expect(fakeStorage.GetPresetsCallCount()).To(Equal(1))
			Expect(listPresetsResp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("when fails to get the presets", func() {
			BeforeEach(func() {
				fakeStorage.GetPresetsReturns(nil, errors.New("Boom!"))
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(listPresetsResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(listPresetsResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(MatchJSON(`{"error": "getting presets: Boom!"}`))
			})
		})
	})

	Describe("GetPresetDetails", func() {
		var (
			preset            types.Preset
			presetDetailsResp *http.Response
		)

		BeforeEach(func() {
			preset = types.Preset{
				Name: "foobar",
			}
		})

		BeforeEach(func() {
			fakeStorage.RetrievePresetReturns(preset, nil)
		})

		JustBeforeEach(func() {
			presetDetailsResp = presetDetails(preset.Name)
		})

		It("returns preset details", func() {
			body, err := ioutil.ReadAll(presetDetailsResp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(presetDetailsResp.StatusCode).To(Equal(http.StatusOK))

			jsonPreset, err := json.Marshal(preset)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(MatchJSON(jsonPreset))
		})

		Context("when fails to get the preset", func() {
			BeforeEach(func() {
				fakeStorage.RetrievePresetReturns(types.Preset{}, errors.New("Boom!"))
			})

			It("returns bad request", func() {
				body, err := ioutil.ReadAll(presetDetailsResp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(presetDetailsResp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(body)).To(Equal(`{"error": "retrieving preset: Boom!"}`))
			})
		})
	})

	Describe("DeletePreset", func() {
		It("deletes the preset", func() {
			presetDetailsResp := deletePreset("foobar")
			Expect(presetDetailsResp.StatusCode).To(Equal(http.StatusOK))

			deleteArgs := fakeStorage.DeletePresetArgsForCall(0)
			Expect(deleteArgs).To(Equal("foobar"))
		})

		Context("when fails to delete a preset", func() {
			It("returns bad request", func() {
				fakeStorage.DeletePresetReturns(types.Preset{}, errors.New("Boom!"))
				presetDetailsResp := deletePreset("foobar")
				Expect(presetDetailsResp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})

})
