package cluster

import (
	"bytes"
	"fmt"
	"github.com/onosproject/onos-test/pkg/onit/cluster/apps"
	"github.com/onosproject/onos-test/pkg/onit/cluster/batch"
	"github.com/onosproject/onos-test/pkg/onit/cluster/core"
	"github.com/onosproject/onos-test/pkg/onit/cluster/extensions"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	"github.com/onosproject/onos-test/pkg/onit/cluster/networking"
	"helm.sh/helm/v3/internal/completion"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	helm "helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/releaseutil"
	"log"
	"os"
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
	filter metav1.ObjectFilter
}

// getConfig gets the Helm configuration
func (c *client) getConfig() (*action.Configuration, error) {
	config := &action.Configuration{}
	if err := config.Init(settings.RESTClientGetter(), c.ObjectsClient.Namespace(), os.Getenv(helmDriverEnv), log.Printf); err != nil {
		return nil, err
	}
	return config, nil
}

// getReleases returns a list of releases
func (c *client) getReleases() ([]*Release, error) {
	config, err := c.getConfig()
	if err != nil {
		return nil, err
	}

	client := action.NewList(config)
	client.SetStateMask()

	releases, err := client.Run()
	if err != nil {
		return nil, err
	}

	results := make([]*Release, len(releases))
	for i, release := range releases {
		release.
		results[i] = newRelease(release.Name, newCha)
	}
	return results, nil
}

// getResources returns a list of resources for the given release
func (c *client) getResources(name string) (helm.ResourceList, error) {
	config, err := c.getConfig()
	if err != nil {
		return nil, err
	}
	releases, err := config.Releases.History(name)
	if err != nil {
		return nil, err
	}
	if len(releases) < 1 {
		return nil, nil
	}

	releaseutil.SortByRevision(releases)
	release := releases[0]

	resources, err := config.KubeClient.Build(bytes.NewBufferString(release.Manifest), true)
	if err != nil {
		return nil, err
	}
	return resources, nil
}
