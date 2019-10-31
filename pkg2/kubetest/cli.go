package kubetest

import (
	"fmt"
	"github.com/dustinkirkland/golang-petname"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

// GetCommand returns the kubetest command
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubetest",
		Short: "Start and manage Kubernetes tests",
	}
	cmd.AddCommand(getRunCommand())
	return cmd
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
	cmd.Flags().String("image-pull-policy", string(corev1.PullAlways), "the image pull policy to use")
	cmd.Flags().Duration("timeout", 0, "the test timeout")
	return cmd
}

// runRunCommand runs the run command
func runRunCommand(cmd *cobra.Command, args []string) error {
	typeName, _ := cmd.Flags().GetString("type")
	image, _ := cmd.Flags().GetString("image")
	pullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	suite, _ := cmd.Flags().GetString("suite")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	var testType TestType
	switch typeName {
	case string(TestTypeTest):
		testType = TestTypeTest
	case string(TestTypeBenchmark):
		testType = TestTypeBenchmark
	default:
		return fmt.Errorf("unknown test type %s", typeName)
	}

	test := &TestConfig{
		TestID:     newTestID(),
		Type:       testType,
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

// newTestID returns a new test ID
func newTestID() string {
	return petname.Generate(2, "-")
}
