package test

import (
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ChartTestSuite struct {
	test.Suite
}

func (s *ChartTestSuite) TestLocalInstall(t *testing.T) {
	chart := cluster.Chart("/etc/charts/onos-topo")
	release := chart.Release("onos-topo")
	err := release.Install(true)
	assert.NoError(t, err)

	deployment, err := release.AppsV1().Deployments().Get("onos-topo")
	assert.NoError(t, err)

	pods, err := deployment.Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 1)
}

func (s *ChartTestSuite) TestRemoteInstall(t *testing.T) {
	chart := cluster.Chart("kafka").
		SetRepository("http://storage.googleapis.com/kubernetes-charts-incubator")
	release := chart.Release("device-simulator-test")
	release.Values().Set("replicas", 1)
	release.Values().Set("zookeeper.replicaCount", 1)
	err := release.Install(true)
	assert.NoError(t, err)
}
