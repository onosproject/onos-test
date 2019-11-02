package onit

import "github.com/onosproject/onos-test/pkg/new/kubetest"

// Tests is the base type for ONIT test suites
type Tests struct {
	*kubetest.Tests
}

// SetupTestSuite sets up the ONOS cluster
func (t *Tests) SetupTestSuite() {
	setupONOSTest(t)
}

// TestSuite is an ONIT test suite
type TestSuite interface {
	kubetest.TestSuite
}

// SetupONOSTestSuite is an interface for setting up an ONOS test
type SetupONOSTestSuite interface {
	SetupONOSTestSuite(setup Setup)
}

// setupONOSTest sets up the ONOS cluster for the given benchmark suite
func setupONOSTest(t TestSuite) {
	if setupONOS, ok := t.(SetupONOSTestSuite); ok {
		setupONOS.SetupONOSTestSuite(NewSetup(t.KubeAPI()))
	}
}
