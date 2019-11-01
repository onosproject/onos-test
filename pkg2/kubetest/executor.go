package kubetest

import (
	"os"
)

// The executor is the entrypoint for test images. It takes the input and environment and runs
// the image in the appropriate context according to the arguments.

// TestContext is the context in which a test image is running
type TestContext string

const (
	testContextEnv = "TEST_CONTEXT"

	TestContextCoordinator TestContext = "coordinator"
	TestContextWorker      TestContext = "worker"
)

// Main runs a test
func Main() {
	config, err := loadTestConfig()
	if err != nil {
		println("Failed to load configuration")
		os.Exit(1)
	}
	if err := Run(config); err != nil {
		println("Test run failed " + err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

// Run runs a test
func Run(config *TestConfig) error {
	context := getTestContext()
	switch context {
	case TestContextCoordinator:
		return runCoordinator(config)
	case TestContextWorker:
		return runWorker(config)
	}
	return nil
}

// getTestContext returns the current test context
func getTestContext() TestContext {
	return TestContext(os.Getenv(testContextEnv))
}

// runCoordinator runs a test image in the coordinator context
func runCoordinator(test *TestConfig) error {
	var coordinator Coordinator
	var err error
	switch test.Type {
	case TestTypeTest:
		coordinator, err = newTestCoordinator(test)
	case TestTypeBenchmark:
		coordinator, err = newBenchmarkCoordinator(test)
	}
	if err != nil {
		return err
	}
	return coordinator.Run()
}

// runWorker runs a test image in the worker context
func runWorker(test *TestConfig) error {
	var worker Worker
	var err error
	switch test.Type {
	case TestTypeTest:
		worker, err = newTestWorker(test)
	case TestTypeBenchmark:
		worker, err = newBenchmarkWorker(test)
	}
	if err != nil {
		return err
	}
	return worker.Run()
}
