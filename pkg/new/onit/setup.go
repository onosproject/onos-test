package onit

import "k8s.io/client-go/rest"

// NewSetup returns a new onit Setup
func NewSetup(config *rest.Config) Setup {
	return &testSetup{}
}

// Setup is an interface for setting up ONOS clusters
type Setup interface {
}

// testSetup is an implementation of the Setup interface
type testSetup struct {
}
