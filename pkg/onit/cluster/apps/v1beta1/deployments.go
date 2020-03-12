package v1beta1

import (
	clustermetav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
)

type DeploymentsClient interface {
	Get(name string) (*Deployment, error)
	List() ([]*Deployment, error)
}

func newDeploymentsClient(objects clustermetav1.ObjectsClient) DeploymentsClient {
	return &deploymentsClient{
		ObjectsClient: objects,
	}
}

type deploymentsClient struct {
	clustermetav1.ObjectsClient
}

func (c *deploymentsClient) Get(name string) (*Deployment, error) {
	object, err := c.ObjectsClient.Get(name, DeploymentResource)
	if err != nil {
		return nil, err
	}
	return newDeployment(object), nil
}

func (c *deploymentsClient) List() ([]*Deployment, error) {
	objects, err := c.ObjectsClient.List(DeploymentResource)
	if err != nil {
		return nil, err
	}
	deployments := make([]*Deployment, len(objects))
	for i, object := range objects {
		deployments[i] = newDeployment(object)
	}
	return deployments, nil
}
