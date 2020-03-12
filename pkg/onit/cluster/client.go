package cluster

import (
	"github.com/onosproject/onos-test/pkg/onit/cluster/apps"
	"github.com/onosproject/onos-test/pkg/onit/cluster/batch"
	"github.com/onosproject/onos-test/pkg/onit/cluster/core"
	"github.com/onosproject/onos-test/pkg/onit/cluster/extensions"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	"github.com/onosproject/onos-test/pkg/onit/cluster/networking"
)

type appsClient apps.Client
type batchClient batch.Client
type coreClient core.Client
type extensionsClient extensions.Client
type networkingClient networking.Client

type Client interface {
	appsClient
	batchClient
	coreClient
	extensionsClient
	networkingClient
}

// newClient creates a new cluster client
func newClient(objects metav1.ObjectsClient) *client {
	return &client{
		ObjectsClient:    objects,
		appsClient:       apps.NewClient(objects),
		batchClient:      batch.NewClient(objects),
		coreClient:       core.NewClient(objects),
		extensionsClient: extensions.NewClient(objects),
		networkingClient: networking.NewClient(objects),
	}
}

// client is an implementation of the Client interface
type client struct {
	metav1.ObjectsClient
	appsClient
	batchClient
	coreClient
	extensionsClient
	networkingClient
}
