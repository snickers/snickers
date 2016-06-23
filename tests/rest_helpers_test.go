package snickers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/flavioribeiro/snickers/rest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rest API", func() {
	Context("helper functions", func() {
		It("should write the error as json", func() {
			w := httptest.NewRecorder()
			rest.HTTPError(w, http.StatusOK, "error here", errors.New("database broken"))

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal(`{"error": "error here: database broken"}`))
		})
	})
})
