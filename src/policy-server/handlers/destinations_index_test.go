package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"policy-server/handlers"
	"policy-server/handlers/fakes"
	storeFakes "policy-server/store/fakes"

	"policy-server/uaa_client"

	apifakes "policy-server/api/fakes"

	"code.cloudfoundry.org/lager"

	"policy-server/store"

	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Destinations index handler", func() {
	var (
		allDestinations      []store.EgressDestination
		expectedResponseBody []byte
		request              *http.Request
		handler              *handlers.DestinationsIndex
		resp                 *httptest.ResponseRecorder
		fakeStore            *storeFakes.EgressDestinationStore
		fakeErrorResponse    *fakes.ErrorResponse
		fakeMapper           *apifakes.EgressDestinationMapper
		logger               *lagertest.TestLogger
		expectedLogger       lager.Logger
		token                uaa_client.CheckTokenResponse
	)

	BeforeEach(func() {
		expectedResponseBody = []byte("some-response")

		var err error
		request, err = http.NewRequest("GET", "/networking/v0/external/destinations", nil)
		Expect(err).NotTo(HaveOccurred())

		fakeStore = &storeFakes.EgressDestinationStore{}
		fakeStore.AllReturns(allDestinations, nil)

		fakeErrorResponse = &fakes.ErrorResponse{}

		fakeMapper = &apifakes.DestinationMapper{}
		fakeMapper.AsBytesReturns(expectedResponseBody, nil)
		logger = lagertest.NewTestLogger("test")
		handler = &handlers.DestinationsIndex{
			ErrorResponse: fakeErrorResponse,
		}

		token = uaa_client.CheckTokenResponse{
			Scope:    []string{"some-scope", "some-other-scope"},
			UserID:   "some-user-id",
			UserName: "some-user",
		}
		resp = httptest.NewRecorder()

		expectedLogger = lager.NewLogger("test").Session("index-destinations")

		testSink := lagertest.NewTestSink()
		expectedLogger.RegisterSink(testSink)
		expectedLogger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))
	})

	It("returns all the destinations", func() {
		MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)

		Expect(fakeStore.AllCallCount()).To(Equal(1))
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.Bytes()).To(Equal(expectedResponseBody))
	})

	//Context("when there are egress policies", func() {
	//	Context("when the user is a network admin", func() {
	//		BeforeEach(func() {
	//			token = uaa_client.CheckTokenResponse{
	//				Scope:    []string{"some-scope", "network.admin"},
	//				UserID:   "some-user-id",
	//				UserName: "some-user",
	//			}
	//		})
	//
	//		It("returns all egress policies", func() {
	//			MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//
	//			Expect(fakeEgressPolicyStore.AllCallCount()).To(Equal(1))
	//			_, egressPolicies := fakeMapper.AsBytesArgsForCall(0)
	//			Expect(egressPolicies).To(Equal(allEgressPolicies))
	//			Expect(resp.Code).To(Equal(http.StatusOK))
	//			Expect(resp.Body.Bytes()).To(Equal(expectedResponseBody))
	//		})
	//
	//		Context("when egressPolicyStore.All returns an error", func() {
	//			BeforeEach(func() {
	//				fakeEgressPolicyStore.AllReturns([]store.EgressPolicy{}, errors.New("I am an error from All"))
	//			})
	//
	//			It("returns a nice error", func() {
	//				MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//
	//				Expect(fakeErrorResponse.InternalServerErrorCallCount()).To(Equal(1))
	//
	//				l, w, err, description := fakeErrorResponse.InternalServerErrorArgsForCall(0)
	//				Expect(l).To(Equal(expectedLogger))
	//				Expect(w).To(Equal(resp))
	//				Expect(err).To(MatchError("I am an error from All"))
	//				Expect(description).To(Equal("getting egress policies failed"))
	//			})
	//		})
	//	})
	//
	//	Context("when the user is not a network admin", func() {
	//		It("does not return any egress policies", func() {
	//			MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//
	//			var emptyEgressPolicies []store.EgressPolicy
	//
	//			Expect(fakeEgressPolicyStore.AllCallCount()).To(Equal(0))
	//			_, egressPolicies := fakeMapper.AsBytesArgsForCall(0)
	//			Expect(egressPolicies).To(Equal(emptyEgressPolicies))
	//			Expect(resp.Code).To(Equal(http.StatusOK))
	//			Expect(resp.Body.Bytes()).To(Equal(expectedResponseBody))
	//		})
	//	})
	//})
	//
	//Context("when the logger isn't on the request context", func() {
	//	It("still works", func() {
	//		MakeRequestWithAuth(handler.ServeHTTP, resp, request, token)
	//
	//		Expect(resp.Code).To(Equal(http.StatusOK))
	//		Expect(resp.Body.Bytes()).To(Equal(expectedResponseBody))
	//	})
	//})
	//
	//Context("when the token isn't on the request context", func() {
	//	It("still works", func() {
	//		MakeRequestWithLogger(handler.ServeHTTP, resp, request, logger)
	//
	//		Expect(fakePolicyFilter.FilterPoliciesCallCount()).To(Equal(1))
	//		_, filterToken := fakePolicyFilter.FilterPoliciesArgsForCall(0)
	//		Expect(filterToken).To(Equal(uaa_client.CheckTokenResponse{}))
	//		Expect(resp.Code).To(Equal(http.StatusOK))
	//		Expect(resp.Body.Bytes()).To(Equal(expectedResponseBody))
	//	})
	//})
	//
	//Context("when rendering the policies as bytes fails", func() {
	//	BeforeEach(func() {
	//		var err error
	//		request, err = http.NewRequest("GET", "/networking/v0/external/policies?id=some-app-guid,yet-another-app-guid", nil)
	//		Expect(err).NotTo(HaveOccurred())
	//
	//		fakeMapper.AsBytesReturns(nil, errors.New("banana"))
	//	})
	//
	//	It("calls the internal server error handler", func() {
	//		MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//
	//		Expect(fakeErrorResponse.InternalServerErrorCallCount()).To(Equal(1))
	//
	//		l, w, err, description := fakeErrorResponse.InternalServerErrorArgsForCall(0)
	//		Expect(l).To(Equal(expectedLogger))
	//		Expect(w).To(Equal(resp))
	//		Expect(err).To(MatchError("banana"))
	//		Expect(description).To(Equal("map policy as bytes failed"))
	//	})
	//})
	//
	//Context("when a list of ids is provided as a query parameter", func() {
	//	BeforeEach(func() {
	//		var err error
	//		request, err = http.NewRequest("GET", "/networking/v0/external/policies?id=some-app-guid,yet-another-app-guid", nil)
	//		Expect(err).NotTo(HaveOccurred())
	//	})
	//
	//	It("filters on only those policies returned by ByGuids", func() {
	//		MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//
	//		Expect(fakeStore.ByGuidsCallCount()).To(Equal(1))
	//		srcGuids, destGuids, inSourceAndDest := fakeStore.ByGuidsArgsForCall(0)
	//		Expect(srcGuids).To(ConsistOf([]string{"some-app-guid", "yet-another-app-guid"}))
	//		Expect(destGuids).To(ConsistOf([]string{"some-app-guid", "yet-another-app-guid"}))
	//		Expect(inSourceAndDest).To(BeFalse())
	//		Expect(fakePolicyFilter.FilterPoliciesCallCount()).To(Equal(1))
	//		policies, userToken := fakePolicyFilter.FilterPoliciesArgsForCall(0)
	//		Expect(policies).To(Equal(byGuidsAPIPolicies))
	//		Expect(userToken).To(Equal(token))
	//		Expect(resp.Code).To(Equal(http.StatusOK))
	//	})
	//
	//	Context("when the id list is empty", func() {
	//		It("filters on only those policies returned by ByGuids", func() {
	//			var err error
	//			request, err = http.NewRequest("GET", "/networking/v0/external/policies?id=", nil)
	//			Expect(err).NotTo(HaveOccurred())
	//
	//			MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//			Expect(fakeStore.ByGuidsCallCount()).To(Equal(1))
	//			srcGuids, destGuids, inSourceAndDest := fakeStore.ByGuidsArgsForCall(0)
	//			Expect(srcGuids).To(Equal([]string{""}))
	//			Expect(destGuids).To(Equal([]string{""}))
	//			Expect(inSourceAndDest).To(BeFalse())
	//			Expect(fakePolicyFilter.FilterPoliciesCallCount()).To(Equal(1))
	//			policies, userToken := fakePolicyFilter.FilterPoliciesArgsForCall(0)
	//			Expect(policies).To(Equal(byGuidsAPIPolicies))
	//			Expect(userToken).To(Equal(token))
	//
	//			Expect(resp.Code).To(Equal(http.StatusOK))
	//		})
	//	})
	//})
	//
	//Context("when dest_id is provided as a query parameter", func() {
	//	BeforeEach(func() {
	//		var err error
	//		request, err = http.NewRequest("GET", "/networking/v0/external/policies?dest_id=not-a-real-app-guid,some-other-app-guid", nil)
	//		Expect(err).NotTo(HaveOccurred())
	//	})
	//
	//	It("filters on those policies with provided dest_id", func() {
	//		MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//
	//		Expect(fakeStore.ByGuidsCallCount()).To(Equal(1))
	//		srcGuids, destGuids, inSourceAndDest := fakeStore.ByGuidsArgsForCall(0)
	//		Expect(srcGuids).To(ConsistOf([]string{}))
	//		Expect(destGuids).To(ConsistOf([]string{"not-a-real-app-guid", "some-other-app-guid"}))
	//		Expect(inSourceAndDest).To(BeFalse())
	//		Expect(fakePolicyFilter.FilterPoliciesCallCount()).To(Equal(1))
	//		policies, userToken := fakePolicyFilter.FilterPoliciesArgsForCall(0)
	//		Expect(policies).To(Equal(byGuidsAPIPolicies))
	//		Expect(userToken).To(Equal(token))
	//		Expect(resp.Code).To(Equal(http.StatusOK))
	//	})
	//
	//	Context("when the dest_id list is empty", func() {
	//		It("filters on only those policies returned by ByGuids", func() {
	//			var err error
	//			request, err = http.NewRequest("GET", "/networking/v0/external/policies?dest_id=", nil)
	//			Expect(err).NotTo(HaveOccurred())
	//
	//			MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//			Expect(fakeStore.ByGuidsCallCount()).To(Equal(1))
	//			srcGuids, destGuids, inSourceAndDest := fakeStore.ByGuidsArgsForCall(0)
	//			Expect(srcGuids).To(Equal([]string{}))
	//			Expect(destGuids).To(Equal([]string{""}))
	//			Expect(inSourceAndDest).To(BeFalse())
	//			Expect(fakePolicyFilter.FilterPoliciesCallCount()).To(Equal(1))
	//			policies, userToken := fakePolicyFilter.FilterPoliciesArgsForCall(0)
	//			Expect(policies).To(Equal(byGuidsAPIPolicies))
	//			Expect(userToken).To(Equal(token))
	//			Expect(resp.Code).To(Equal(http.StatusOK))
	//		})
	//	})
	//})
	//
	//Context("when source_id is provided as a query parameter", func() {
	//	BeforeEach(func() {
	//		var err error
	//		request, err = http.NewRequest("GET", "/networking/v0/external/policies?source_id=some-app-guid,yet-another-app-guid,some-other-app-guid", nil)
	//		Expect(err).NotTo(HaveOccurred())
	//	})
	//
	//	It("filters on those policies with provided source_id", func() {
	//		MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//
	//		Expect(fakeStore.ByGuidsCallCount()).To(Equal(1))
	//		srcGuids, destGuids, inSourceAndDest := fakeStore.ByGuidsArgsForCall(0)
	//		Expect(srcGuids).To(ConsistOf([]string{"some-app-guid", "yet-another-app-guid", "some-other-app-guid"}))
	//		Expect(destGuids).To(ConsistOf([]string{}))
	//		Expect(inSourceAndDest).To(BeFalse())
	//		Expect(fakePolicyFilter.FilterPoliciesCallCount()).To(Equal(1))
	//		policies, userToken := fakePolicyFilter.FilterPoliciesArgsForCall(0)
	//		Expect(policies).To(Equal(byGuidsAPIPolicies))
	//		Expect(userToken).To(Equal(token))
	//		Expect(resp.Code).To(Equal(http.StatusOK))
	//	})
	//
	//	Context("when the source_id list is empty", func() {
	//		It("filters on only those policies returned by ByGuids", func() {
	//			var err error
	//			request, err = http.NewRequest("GET", "/networking/v0/external/policies?source_id=", nil)
	//			Expect(err).NotTo(HaveOccurred())
	//
	//			MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//			Expect(fakeStore.ByGuidsCallCount()).To(Equal(1))
	//			srcGuids, destGuids, inSourceAndDest := fakeStore.ByGuidsArgsForCall(0)
	//			Expect(srcGuids).To(Equal([]string{""}))
	//			Expect(destGuids).To(Equal([]string{}))
	//			Expect(inSourceAndDest).To(BeFalse())
	//			Expect(fakePolicyFilter.FilterPoliciesCallCount()).To(Equal(1))
	//			policies, userToken := fakePolicyFilter.FilterPoliciesArgsForCall(0)
	//			Expect(policies).To(Equal(byGuidsAPIPolicies))
	//			Expect(userToken).To(Equal(token))
	//
	//			Expect(resp.Code).To(Equal(http.StatusOK))
	//		})
	//	})
	//})
	//
	//Context("when dest_id and source_id is provided as a query parameter", func() {
	//	BeforeEach(func() {
	//		var err error
	//		request, err = http.NewRequest("GET", "/networking/v0/external/policies?source_id=some-app-guid,meow&dest_id=not-a-real-app-guid,some-other-app-guid", nil)
	//		Expect(err).NotTo(HaveOccurred())
	//	})
	//
	//	It("filters on those policies with provided source_id and dest_id", func() {
	//		MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//
	//		Expect(fakeStore.ByGuidsCallCount()).To(Equal(1))
	//		srcGuids, destGuids, inSourceAndDest := fakeStore.ByGuidsArgsForCall(0)
	//		Expect(srcGuids).To(ConsistOf([]string{"some-app-guid", "meow"}))
	//		Expect(destGuids).To(ConsistOf([]string{"not-a-real-app-guid", "some-other-app-guid"}))
	//		Expect(inSourceAndDest).To(BeTrue())
	//		Expect(fakePolicyFilter.FilterPoliciesCallCount()).To(Equal(1))
	//		policies, userToken := fakePolicyFilter.FilterPoliciesArgsForCall(0)
	//		Expect(policies).To(Equal(byGuidsAPIPolicies))
	//		Expect(userToken).To(Equal(token))
	//		Expect(resp.Code).To(Equal(http.StatusOK))
	//	})
	//})
	//
	//Context("when the store throws an error", func() {
	//	BeforeEach(func() {
	//		fakeStore.AllReturns(nil, errors.New("banana"))
	//	})
	//
	//	It("calls the internal server error handler", func() {
	//		MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//
	//		Expect(fakeErrorResponse.InternalServerErrorCallCount()).To(Equal(1))
	//
	//		l, w, err, description := fakeErrorResponse.InternalServerErrorArgsForCall(0)
	//		Expect(l).To(Equal(expectedLogger))
	//		Expect(w).To(Equal(resp))
	//		Expect(err).To(MatchError("banana"))
	//		Expect(description).To(Equal("database read failed"))
	//	})
	//})
	//
	//Context("when the policy filter throws an error", func() {
	//	BeforeEach(func() {
	//		fakePolicyFilter.FilterPoliciesReturns(nil, errors.New("banana"))
	//	})
	//
	//	It("calls the internal server error handler", func() {
	//		MakeRequestWithLoggerAndAuth(handler.ServeHTTP, resp, request, logger, token)
	//
	//		Expect(fakeErrorResponse.InternalServerErrorCallCount()).To(Equal(1))
	//
	//		l, w, err, description := fakeErrorResponse.InternalServerErrorArgsForCall(0)
	//		Expect(l).To(Equal(expectedLogger))
	//		Expect(w).To(Equal(resp))
	//		Expect(err).To(MatchError("banana"))
	//		Expect(description).To(Equal("filter policies failed"))
	//	})
	//})
})
