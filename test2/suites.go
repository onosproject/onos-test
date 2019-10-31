package test2

import (
	"github.com/atomix/atomix-go-client/pkg/client"
	"github.com/onosproject/onos-test/pkg2/kubetest"
)

type TestsOne struct {
	*kubetest.Tests
}

func (t *TestsOne) SetupTestSuite(client client.Client) {

}

type TestsTwo struct {
	*kubetest.Tests
}

type TestsThree struct {
	*kubetest.Tests
}

type BenchmarksOne struct {
	*kubetest.Benchmarks
}
