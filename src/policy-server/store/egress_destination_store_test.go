package store_test

import (
	//dbfakes "policy-server/db/fakes"
	"policy-server/store"
	"policy-server/store/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EgressDestinationStore", func() {
	var (
		egressDestinationsStore *store.EgressDestinationStore
		egressDestinationRepo   *fakes.EgressDestinationRepo
		//mockDb                  *fakes.Db
		//
		//tx                 *dbfakes.Transaction
		egressDestinations []store.EgressDestination
	)

	BeforeEach(func() {
		egressDestinationRepo = &fakes.EgressDestinationRepo{}
		egressDestinations = []store.EgressDestination{
			{},
		}
		egressDestinationsStore = &store.EgressDestinationStore{
			EgressDestinationRepo: egressDestinationRepo,
		}
	})

	// TODO make this use a real db and table to test the EgressDestinationTable

	FDescribe("All", func() {
		Context("when there are policies created", func() {
			BeforeEach(func() {
				egressDestinationRepo.AllReturns(egressDestinations, nil)
			})

			It("should return a list of all policies", func() {
				destinations, err := egressDestinationsStore.All()
				Expect(err).NotTo(HaveOccurred())
				Expect(destinations).To(Equal(egressDestinations))
			})
		})
	})
})
