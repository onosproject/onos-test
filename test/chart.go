package test

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ChartTestSuite struct {
	test.Suite
}

func (s *ChartTestSuite) TestLocalInstall(t *testing.T) {
	atomix := cluster.Chart("/etc/charts/atomix").
		Release("atomix-controller")
	atomix.Values().Set("namespace", cluster.Namespace())
	err := atomix.Install(true)
	assert.NoError(t, err)

	topo := cluster.Chart("/etc/charts/onos-topo").
		Release("onos-topo")
	topo.Values().Set("store.controller", fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", cluster.Namespace()))
	err = topo.Install(true)
	assert.NoError(t, err)

	deployment, err := topo.Apps().V1().Deployments().Get("onos-topo")
	assert.NoError(t, err)

	pods, err := deployment.Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 1)
}

func (s *ChartTestSuite) TestRemoteInstall(t *testing.T) {
	kafka := cluster.Chart("kafka").
		SetRepository("http://storage.googleapis.com/kubernetes-charts-incubator").
		Release("device-simulator-test")
	kafka.Values().Set("replicas", 1)
	kafka.Values().Set("zookeeper.replicaCount", 1)
	err := kafka.Install(true)
	assert.NoError(t, err)
}
