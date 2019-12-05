// Copyright 2019-present Open Networking Foundation.
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
	"github.com/onosproject/onos-test/pkg/kube"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

// newWorker returns a new test worker
func newWorker(config *Config) (*Worker, error) {
	kubeAPI, err := kube.GetAPI(config.ID)
	if err != nil {
		return nil, err
	}
	return &Worker{
		client: kubeAPI.Client(),
		config: config,
	}, nil
}

// Worker runs a test job
type Worker struct {
	client client.Client
	config *Config
}

// Run runs a test
func (w *Worker) Run() error {
	test, ok := Registry.tests[w.config.Suite]
	if !ok {
		return fmt.Errorf("unknown test suite %s", w.config.Suite)
	}

	tests := []testing.InternalTest{
		{
			Name: w.config.Suite,
			F: func(t *testing.T) {
				RunTests(t, test, w.config)
			},
		},
	}

	// Hack to enable verbose testing.
	os.Args = []string{
		os.Args[0],
		"-test.v",
	}

	testing.Main(func(_, _ string) (bool, error) { return true, nil }, tests, nil, nil)
	return nil
}
