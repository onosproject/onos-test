// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package onostest

const (
	AtomixChartRepo                 = "https://charts.atomix.io"
	OnosChartRepo                   = "https://charts.onosproject.org"
	SdranChartRepo                  = "https://sdrancharts.onosproject.org"
	AtomixControllerPort            = "5679"
	SecretsName                     = "helmit-secrets"
	ControllerChartName             = "atomix-controller"
	RaftStorageControllerChartName  = "raft-storage-controller"
	CacheStorageControllerChartName = "cache-storage-controller"
)

func AtomixName(testName string, componentName string) string {
	return testName + "-" + componentName + "-atomix"
}

func AtomixControllerName(testName string, componentName string) string {
	return AtomixName(testName, componentName) + "-atomix-controller"
}

func AtomixController(testName string, componentName string) string {
	return AtomixControllerName(testName, componentName) + ":" + AtomixControllerPort
}

func RaftReleaseName(componentName string) string {
	return componentName + "-raft"
}

func CacheReleaseName(componentName string) string {
	return componentName + "-cache"
}
