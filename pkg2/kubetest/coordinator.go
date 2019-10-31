package kubetest

import (
	"fmt"
	"github.com/dustinkirkland/golang-petname"
	"github.com/onosproject/onos-test/pkg2/util/k8s"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// newTestCoordinator returns a new test coordinator
func newTestCoordinator(test *TestConfig) (Coordinator, error) {
	client, err := k8s.GetClient()
	if err != nil {
		return nil, err
	}
	return &TestCoordinator{
		client: client,
		test:   test,
	}, nil
}

// newBenchmarkCoordinator returns a new benchmark coordinator
func newBenchmarkCoordinator(test *TestConfig) (Coordinator, error) {
	client, err := k8s.GetClient()
	if err != nil {
		return nil, err
	}
	return &BenchmarkCoordinator{
		client: client,
		test:   test,
	}, nil
}

// Coordinator coordinates workers for tests and benchmarks
type Coordinator interface {
	// Run runs the coordinator
	Run() error
}

// TestCoordinator coordinates workers for suites of tests
type TestCoordinator struct {
	client client.Client
	test   *TestConfig
}

// Run runs the tests
func (c *TestCoordinator) Run() error {
	jobs := make([]*TestJob, 0)
	if c.test.Suite == "" {
		for suite := range Registry.tests {
			config := &TestConfig{
				TestID:     newJobID(c.test.TestID),
				Type:       c.test.Type,
				Image:      c.test.Image,
				Suite:      suite,
				Timeout:    c.test.Timeout,
				PullPolicy: c.test.PullPolicy,
			}
			job := &TestJob{
				client: c.client,
				test:   config,
			}
			jobs = append(jobs, job)
		}
	} else {
		config := &TestConfig{
			TestID:     newJobID(c.test.TestID),
			Type:       c.test.Type,
			Image:      c.test.Image,
			Suite:      c.test.Suite,
			Timeout:    c.test.Timeout,
			PullPolicy: c.test.PullPolicy,
		}
		job := &TestJob{
			client: c.client,
			test:   config,
		}
		jobs = append(jobs, job)
	}
	return runJobs(jobs)
}

// BenchmarkCoordinator coordinates workers for suites of benchmarks
type BenchmarkCoordinator struct {
	client client.Client
	test   *TestConfig
}

// Run runs the tests
func (c *BenchmarkCoordinator) Run() error {
	jobs := make([]*TestJob, 0)
	if c.test.Suite == "" {
		for suite := range Registry.benchmarks {
			config := &TestConfig{
				TestID:     newJobID(c.test.TestID),
				Type:       c.test.Type,
				Image:      c.test.Image,
				Suite:      suite,
				Timeout:    c.test.Timeout,
				PullPolicy: c.test.PullPolicy,
			}
			job := &TestJob{
				client: c.client,
				test:   config,
			}
			jobs = append(jobs, job)
		}
	} else {
		config := &TestConfig{
			TestID:     newJobID(c.test.TestID),
			Type:       c.test.Type,
			Image:      c.test.Image,
			Suite:      c.test.Suite,
			Timeout:    c.test.Timeout,
			PullPolicy: c.test.PullPolicy,
		}
		job := &TestJob{
			client: c.client,
			test:   config,
		}
		jobs = append(jobs, job)
	}
	return runJobs(jobs)
}

// runJobs runs the given test jobs
func runJobs(jobs []*TestJob) error {
	for _, job := range jobs {
		if err := job.Start(); err != nil {
			return err
		}
	}

	for _, job := range jobs {
		if err := job.WaitForComplete(); err != nil {
			return err
		}
	}

	exitCode := 0
	for _, job := range jobs {
		output, code, err := job.GetResult()
		if err != nil {
			return err
		}
		_, _ = os.Stdout.WriteString(output)
		if code != 0 {
			exitCode = code
		}
	}

	for _, job := range jobs {
		_ = job.TearDown()
	}
	os.Exit(exitCode)
	return nil
}

// newJobID returns a new unique test job ID
func newJobID(testID string) string {
	return fmt.Sprintf("%s-%s", testID, petname.Generate(2, "-"))
}
