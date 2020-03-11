package setup

import "github.com/onosproject/onos-test/pkg/onit/cluster"

// ChartSetup is an interface for setting up a chart
type ChartSetup interface {
	// SetRelease sets the chart's release name
	SetRelease(name string) ChartSetup

	// SetRepository sets the chart's repository
	SetRepository(url string) ChartSetup

	// SetValue sets the value of a path in the chart's configuration
	SetValue(path string, value interface{}) ChartSetup

	// AddValue adds a value to a list in the chart's configuration
	AddValue(path string, value interface{}) ChartSetup
}

func (s *clusterSetup) Chart(name string) ChartSetup {
	chart := s.getChart(name, func() *cluster.Chart {
		chart := s.cluster.Charts().New()
		chart.SetChart(name)
		return chart
	})
	return &clusterChartSetup{
		chart: chart,
	}
}

// clusterChartSetup is an implementation of the ChartSetup interface
type clusterChartSetup struct {
	chart *cluster.Chart
}

func (s *clusterChartSetup) SetRelease(name string) ChartSetup {
	s.chart.SetRelease(name)
	return s
}

func (s *clusterChartSetup) SetRepository(url string) ChartSetup {
	s.chart.SetRepository(url)
	return s
}

func (s *clusterChartSetup) SetValue(path string, value interface{}) ChartSetup {
	s.chart.SetValue(path, value)
	return s
}

func (s *clusterChartSetup) AddValue(path string, value interface{}) ChartSetup {
	s.chart.AddValue(path, value)
	return s
}

func (s *clusterChartSetup) setup() error {
	return s.chart.Setup()
}
