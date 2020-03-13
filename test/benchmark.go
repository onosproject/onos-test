package test

import (
	"github.com/onosproject/onos-test/pkg/benchmark"
	"github.com/onosproject/onos-test/pkg/input"
	"time"
)

type ChartBenchmarkSuite struct {
	benchmark.Suite
	value input.Source
}

func (s *ChartBenchmarkSuite) SetupWorker(b *benchmark.Context) {
	s.value = input.RandomString(8)
}

func (s *ChartBenchmarkSuite) BenchmarkTest(b *benchmark.Benchmark) error {
	println(s.value.Next().String())
	time.Sleep(time.Second)
	return nil
}
