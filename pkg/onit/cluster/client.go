package cluster

import (
    apps "github.com/onosproject/onos-test/pkg/onit/cluster/apps"
    batch "github.com/onosproject/onos-test/pkg/onit/cluster/batch"
    core "github.com/onosproject/onos-test/pkg/onit/cluster/core"
    extensions "github.com/onosproject/onos-test/pkg/onit/cluster/extensions"
    networking "github.com/onosproject/onos-test/pkg/onit/cluster/networking"
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

func newClient(objects metav1.ObjectsClient) *client {
	return &client{
		ObjectsClient:    objects,
        appsClient: apps.NewClient(objects),
        batchClient: batch.NewClient(objects),
        coreClient: core.NewClient(objects),
        extensionsClient: extensions.NewClient(objects),
        networkingClient: networking.NewClient(objects),
	}
}

type client struct {
	metav1.ObjectsClient
    appsClient
    batchClient
    coreClient
    extensionsClient
    networkingClient
}
