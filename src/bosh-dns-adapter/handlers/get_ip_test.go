package handlers_test

import (
	"bosh-dns-adapter/handlers"
	"bosh-dns-adapter/handlers/fakes"
	"errors"
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetIP", func() {
	var (
		getIP handlers.GetIP

		fakeSDCClient *fakes.SDCClient

		resp    *httptest.ResponseRecorder
		request *http.Request
		logger  *lagertest.TestLogger
	)

	BeforeEach(func() {
		fakeSDCClient = &fakes.SDCClient{}
		fakeSDCClient.IPsReturns([]string{"192.168.0.1"}, nil)
		logger = lagertest.NewTestLogger("get ip handler test logger")
		getIP = handlers.GetIP{
			SDCClient: fakeSDCClient,
			Logger:    logger,
		}

		resp = httptest.NewRecorder()
	})

	Context("when the user requests an A record for a given hostname", func() {
		BeforeEach(func() {
			var err error
			request, err = http.NewRequest("GET", "?type=1&name=app.example.com.", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("return an A record response with a list of ips", func() {
			getIP.ServeHTTP(resp, request)

			Expect(fakeSDCClient.IPsCallCount()).To(Equal(1))
			Expect(fakeSDCClient.IPsArgsForCall(0)).To(Equal("app.example.com."))

			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Body.String()).To(MatchJSON(`{
					"Status": 0,
					"TC": false,
					"RD": false,
					"RA": false,
					"AD": false,
					"CD": false,
					"Question":
					[
						{
							"name": "app.example.com.",
							"type": 1
						}
					],
					"Answer":
					[
						{
							"name": "app.example.com.",
							"type": 1,
							"TTL":  0,
							"data": "192.168.0.1"
						}
					],
					"Additional": [ ],
					"edns_client_subnet": "0.0.0.0/0"
				}`))
		})
	})

	Context("when the user provides only a hostname", func() {
		BeforeEach(func() {
			var err error
			request, err = http.NewRequest("GET", "?name=app.example.com.", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an A record with a list of ips", func() {
			getIP.ServeHTTP(resp, request)

			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Body.String()).To(MatchJSON(`{
					"Status": 0,
					"TC": false,
					"RD": false,
					"RA": false,
					"AD": false,
					"CD": false,
					"Question":
					[
						{
							"name": "app.example.com.",
							"type": 1
						}
					],
					"Answer":
					[
						{
							"name": "app.example.com.",
							"type": 1,
							"TTL":  0,
							"data": "192.168.0.1"
						}
					],
					"Additional": [ ],
					"edns_client_subnet": "0.0.0.0/0"
				}`))
		})
	})

	Context("when the user makes a requests without providing the hostname", func() {
		BeforeEach(func() {
			var err error
			request, err = http.NewRequest("GET", "?type=1", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a http 400 status", func() {
			getIP.ServeHTTP(resp, request)

			Expect(resp.Code).To(Equal(http.StatusBadRequest))
			Expect(resp.Body.String()).To(MatchJSON(`{
				"Status": 2,
				"TC": false,
				"RD": false,
				"RA": false,
				"AD": false,
				"CD": false,
				"Question":
				[
					{
						"name": "",
						"type": 1
					}
				],
				"Answer": [ ],
				"Additional": [ ],
				"edns_client_subnet": "0.0.0.0/0"
			}`))
		})
	})

	Context("when requesting anything but an A record", func() {
		It("should return a successful response with no answers", func() {
			request, err := http.NewRequest("GET", "?type=16&name=app-id.internal.local.", nil)
			Expect(err).ToNot(HaveOccurred())

			getIP.ServeHTTP(resp, request)

			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Body.String()).To(MatchJSON(`{
					"Status": 0,
					"TC": false,
					"RD": false,
					"RA": false,
					"AD": false,
					"CD": false,
					"Question":
					[
						{
							"name": "app-id.internal.local.",
							"type": 16
						}
					],
					"Answer": [ ],
					"Additional": [ ],
					"edns_client_subnet": "0.0.0.0/0"
				}`))
		})
	})

	Context("when the sdc client returns an error", func() {
		BeforeEach(func() {
			fakeSDCClient.IPsReturns(nil, errors.New("failed to get ips"))
		})

		It("returns a http 500 status", func() {
			request, err := http.NewRequest("GET", "?type=1&name=app-id.internal.local.", nil)
			Expect(err).To(Succeed())

			getIP.ServeHTTP(resp, request)

			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
