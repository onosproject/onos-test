package test2

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

type SuiteOne struct {

}

func (s *SuiteOne) SetUp(client client.Client) {
	panic("implement me")
}

func (s *SuiteOne) Run(t *testing.T) {
	panic("implement me")
}

type SuiteTwo struct {

}

func (s *SuiteTwo) SetUp(client client.Client) {
	panic("implement me")
}

func (s *SuiteTwo) Run(t *testing.T) {
	panic("implement me")
}

type SuiteThree struct {

}

func (s *SuiteThree) SetUp(client client.Client) {
	panic("implement me")
}

func (s *SuiteThree) Run(t *testing.T) {
	panic("implement me")
}
