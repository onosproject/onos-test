package kubetest

import (
	"fmt"
	"github.com/dustinkirkland/golang-petname"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"math/rand"
	"time"
)

// GetCommand returns the kubetest command
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubetest",
		Short: "Start and manage Kubernetes tests",
	}
	cmd.AddCommand(getRunCommand())
	cmd.AddCommand(getTestCommand())
	cmd.AddCommand(getBenchCommand())
	return cmd
}

// getTestCommand returns the test command
func getTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run a test",
		RunE:  runTestCommand,
	}
	cmd.Flags().StringP("image", "i", "", "the test image to run")
	cmd.Flags().StringP("suite", "s", "", "the name of a suite to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the image pull policy to use")
	cmd.Flags().Duration("timeout", 0, "the test timeout")
	return cmd
}

// runTestCommand runs the test command
func runTestCommand(cmd *cobra.Command, _ []string) error {
	image, _ := cmd.Flags().GetString("image")
	pullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	suite, _ := cmd.Flags().GetString("suite")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	test := &TestConfig{
		TestID:     newTestID(),
		Type:       TestTypeTest,
		Image:      image,
		Suite:      suite,
		Timeout:    timeout,
		PullPolicy: corev1.PullPolicy(pullPolicy),
	}

	runner, err := newTestRunner(test)
	if err != nil {
		return err
	}
	return runner.Run()
}

// getBenchCommand returns the bench command
func getBenchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bench",
		Short: "Run a benchmark",
		RunE:  runBenchCommand,
	}
	cmd.Flags().StringP("image", "i", "", "the test image to run")
	cmd.Flags().StringP("suite", "s", "", "the name of a suite to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the image pull policy to use")
	cmd.Flags().Duration("timeout", 0, "the test timeout")
	return cmd
}

// runBenchCommand runs the bench command
func runBenchCommand(cmd *cobra.Command, _ []string) error {
	image, _ := cmd.Flags().GetString("image")
	pullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	suite, _ := cmd.Flags().GetString("suite")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	test := &TestConfig{
		TestID:     newTestID(),
		Type:       TestTypeBenchmark,
		Image:      image,
		Suite:      suite,
		Timeout:    timeout,
		PullPolicy: corev1.PullPolicy(pullPolicy),
	}

	runner, err := newTestRunner(test)
	if err != nil {
		return err
	}
	return runner.Run()
}

// getRunCommand returns the run command
func getRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a test",
		RunE:  runRunCommand,
	}
	cmd.Flags().StringP("type", "t", "", "the type of test to run")
	cmd.Flags().StringP("image", "i", "", "the test image to run")
	cmd.Flags().StringP("suite", "s", "", "the name of a suite to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the image pull policy to use")
	cmd.Flags().Duration("timeout", 0, "the test timeout")
	return cmd
}

// runRunCommand runs the run command
func runRunCommand(cmd *cobra.Command, args []string) error {
	typeName, _ := cmd.Flags().GetString("type")
	switch typeName {
	case string(TestTypeTest):
		return runTestCommand(cmd, []string{})
	case string(TestTypeBenchmark):
		return runBenchCommand(cmd, []string{})
	default:
		return fmt.Errorf("unknown test type %s", typeName)
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// newTestID returns a new test ID
func newTestID() string {
	return petname.Generate(2, "-")
}
