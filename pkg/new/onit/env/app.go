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

package env

import (
	"github.com/onosproject/onos-test/pkg/new/onit/setup"
	corev1 "k8s.io/api/core/v1"
)

// AppEnv provides the environment for an app
type AppEnv interface {
	ServiceEnv
}

var _ AppEnv = &appEnv{}

// appEnv is an implementation of the AppEnv interface
type appEnv struct {
	*serviceEnv
}

var _ setup.AppSetup = &appSetup{}

// appSetup is an implementation of the AppSetup interface
type appSetup struct {
	*testEnv
	name       string
	image      string
	pullPolicy corev1.PullPolicy
}

func (s *appSetup) Name(name string) setup.AppSetup {
	s.name = name
	return s
}

func (s *appSetup) Image(image string) setup.AppSetup {
	s.image = image
	return s
}

func (s *appSetup) PullPolicy(pullPolicy corev1.PullPolicy) setup.AppSetup {
	s.pullPolicy = pullPolicy
	return s
}

func (s *appSetup) Setup() error {
	return nil
}
