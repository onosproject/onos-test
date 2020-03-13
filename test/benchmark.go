package test

import (
	"github.com/onosproject/onos-test/pkg/benchmark"
	"time"
)

type ChartBenchmarkSuite struct {
	benchmark.Suite
}

func (s *ChartBenchmarkSuite) BenchmarkTest(b *benchmark.Benchmark) {
	time.Sleep(time.Second)
}
