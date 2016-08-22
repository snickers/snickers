package server_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/db/dbfakes"
	"github.com/snickers/snickers/server"
)

var _ = Describe("Snickers Server", func() {
	var (
		logger         *lagertest.TestLogger
		snickersServer *server.SnickersServer
		fakeStorage    *dbfakes.FakeStorage
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("snickers-test")
		fakeStorage = new(dbfakes.FakeStorage)
	})

	Context("when passed a socket", func() {
		var (
			socketPath string
			tmpDir     string
		)

		JustBeforeEach(func() {
			var err error
			tmpDir, err = ioutil.TempDir(os.TempDir(), "snickers-server-test")
			socketPath = path.Join(tmpDir, "snickers.sock")
			snickersServer = server.New(logger, "unix", socketPath, fakeStorage)
			Expect(err).NotTo(HaveOccurred())

			err = snickersServer.Start(false)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(tmpDir)
		})

		Describe("Start", func() {
			It("listens on the socket provided", func() {
				info, err := os.Stat(socketPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(info).NotTo(BeNil())
			})
		})

		Describe("Stop", func() {
			JustBeforeEach(func() {
				info, err := os.Stat(socketPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(info).NotTo(BeNil())
			})

			It("removes the existing socket", func() {
				Expect(snickersServer.Stop()).To(Succeed())

				info, err := os.Stat(socketPath)
				Expect(err).To(HaveOccurred())
				Expect(info).To(BeNil())
			})

			Context("when fails to stop the server because it's already stopped", func() {
				JustBeforeEach(func() {
					Expect(snickersServer.Stop()).NotTo(HaveOccurred())
				})

				It("returns an error", func() {
					Expect(snickersServer.Stop()).To(HaveOccurred())
				})
			})
		})
	})

	Context("when passed a tcp addr", func() {
		var httpClient *http.Client

		JustBeforeEach(func() {
			var err error
			port := fmt.Sprintf(":%d", 8000+config.GinkgoConfig.ParallelNode)
			snickersServer = server.New(logger, "tcp", port, fakeStorage)

			err = snickersServer.Start(false)
			Expect(err).NotTo(HaveOccurred())

			httpClient = &http.Client{
				Transport: &http.Transport{
					Dial: func(network, addr string) (net.Conn, error) {
						return net.DialTimeout("tcp", port, 2*time.Second)
					},
				},
			}
		})

		It("listens on the address provided", func() {
			request, err := http.NewRequest(http.MethodGet,
				fmt.Sprintf("http://server%s", server.Routes[server.Ping].Path), nil)

			Expect(err).NotTo(HaveOccurred())
			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
	})
})
