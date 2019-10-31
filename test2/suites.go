package test2

import (
	"github.com/onosproject/onos-test/pkg2/kubetest"
)

type TestsOne struct {
	*kubetest.Tests
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
