package cluster

import "github.com/onosproject/onos-test/pkg/onit/api"

var instance *api.API

func getAPI() *api.API {
	if instance == nil {
		instance = api.NewFromEnv()
	}
	return instance
}

// Namespace returns the cluster namespace
func Namespace() string {
	return getAPI().Namespace()
}

// Charts returns a list of charts in the cluster
func Charts() []*api.Chart {
	return getAPI().Charts()
}

// Chart returns a chart
func Chart(name string) *api.Chart {
	return getAPI().Chart(name)
}

// Releases returns a list of releases
func Releases() []*api.Release {
	return getAPI().Releases()
}

// Release returns the release with the given name
func Release(name string) *api.Release {
	return getAPI().Release(name)
}
