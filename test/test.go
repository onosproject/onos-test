// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/helm"
	"github.com/onosproject/onos-test/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ChartTestSuite struct {
	test.Suite
}

func (s *ChartTestSuite) TestLocalInstall(t *testing.T) {
	atomix := helm.Helm().
		Chart("/etc/charts/atomix-controller").
		Release("atomix-controller").
		Set("namespace", helm.Namespace())
	err := atomix.Install(true)
	assert.NoError(t, err)

	topo := helm.Helm().
		Chart("/etc/charts/onos-topo").
		Release("onos-topo").
		Set("store.controller", fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", helm.Namespace()))
	err = topo.Install(true)
	assert.NoError(t, err)

	pods, err := topo.Core().V1().Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 1)

	deployment, err := topo.Apps().
		V1().
		Deployments().
		Get("onos-topo")
	assert.NoError(t, err)

	pods, err = deployment.Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 1)
}

func (s *ChartTestSuite) TestRemoteInstall(t *testing.T) {
	kafka := helm.Helm().
		Chart("kafka").
		SetRepository("http://storage.googleapis.com/kubernetes-charts-incubator").
		Release("kafka").
		Set("replicas", 1).
		Set("zookeeper.replicaCount", 1)
	err := kafka.Install(true)
	assert.NoError(t, err)
}
