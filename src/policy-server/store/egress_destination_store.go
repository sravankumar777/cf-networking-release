package store

//go:generate counterfeiter -o fakes/egress_destination_repo.go --fake-name EgressDestinationRepo . egressDestinationRepo
type egressDestinationRepo interface {
	All() ([]EgressDestination, error)
}

type EgressDestinationStore struct {
	EgressDestinationRepo egressDestinationRepo
}

func (e *EgressDestinationStore) All() ([]EgressDestination, error) {
	return e.EgressDestinationRepo.All()
}