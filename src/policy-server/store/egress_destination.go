package store

type EgressDestinationTable struct {}

func (e *EgressDestinationTable)All() ([]EgressDestination, error) {
	return []EgressDestination{}, nil
}