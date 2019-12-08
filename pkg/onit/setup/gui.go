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

package setup

import (
	"github.com/onosproject/onos-test/pkg/onit/cluster"

	corev1 "k8s.io/api/core/v1"
)

// GuiSetup is an interface for setting up gui nodes
type GuiSetup interface {

	// SetEnabled enables the Gui
	SetEnabled() GuiSetup

	// SetReplicas sets the number of onos-gui replicas to deploy
	SetReplicas(replicas int) GuiSetup

	// SetImage sets the onos-gui image to deploy
	SetImage(image string) GuiSetup

	// SetPullPolicy sets the image pull policy
	SetPullPolicy(pullPolicy corev1.PullPolicy) GuiSetup
}

var _ GuiSetup = &clusterGuiSetup{}

// clusterGuiSetup is an implementation of the Gui interface
type clusterGuiSetup struct {
	gui *cluster.Gui
}

func (s *clusterGuiSetup) SetEnabled() GuiSetup {
	s.gui.SetEnabled(true)
	return s
}

func (s *clusterGuiSetup) SetReplicas(replicas int) GuiSetup {
	s.gui.SetReplicas(replicas)
	return s
}

func (s *clusterGuiSetup) SetImage(image string) GuiSetup {
	s.gui.SetImage(image)
	return s
}

func (s *clusterGuiSetup) SetPullPolicy(pullPolicy corev1.PullPolicy) GuiSetup {
	s.gui.SetPullPolicy(pullPolicy)
	return s
}

func (s *clusterGuiSetup) setup() error {
	if s.gui.Enabled() {
		err := s.gui.OnosConfig.Setup()
		if err != nil {
			return err
		}
		err = s.gui.OnosTopo.Setup()
		if err != nil {
			return err
		}
		err = s.gui.Setup()
		if err != nil {
			return err
		}
	}
	return nil
}
