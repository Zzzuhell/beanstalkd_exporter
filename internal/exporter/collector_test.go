package exporter

import (
	"bytes"
	"fmt"
	"log/slog"
	"reflect"
	"testing"

	"github.com/davidtannock/beanstalkd_exporter/v2/internal/beanstalkd"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func readCounter(m prometheus.Counter) float64 {
	// TODO: Revisit this once client_golang offers better testing tools.
	pb := &dto.Metric{}
	err := m.Write(pb)
	if err != nil {
		return -1
	}
	return pb.GetCounter().GetValue()
}

func readGauge(m prometheus.Gauge) float64 {
	// TODO: Revisit this once client_golang offers better testing tools.
	pb := &dto.Metric{}
	err := m.Write(pb)
	if err != nil {
		return -1
	}
	return pb.GetGauge().GetValue()
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		opts          CollectorOpts
		expectedError string
	}{
		// We expect an error when there's an unknown system metric.
		{
			opts:          CollectorOpts{SystemMetrics: []string{"does_not_exist"}},
			expectedError: "unknown system metric: does_not_exist",
		},
		// We expect an error when there's an unknown tube metric.
		{
			opts:          CollectorOpts{TubeMetrics: []string{"tube_no_exist"}},
			expectedError: "unknown tube metric: tube_no_exist",
		},
		// If specific tubes metrics are requested, we must have tubes.
		{
			opts:          CollectorOpts{Tubes: []string{}, TubeMetrics: []string{"tube_current_jobs_ready_count"}},
			expectedError: "tube metrics without tubes is not supported",
		},
		{
			opts:          CollectorOpts{AllTubes: false, TubeMetrics: []string{"tube_current_jobs_ready_count"}},
			expectedError: "tube metrics without tubes is not supported",
		},
	}

	for _, tt := range tests {
		actualError := tt.opts.validate()
		if actualError.Error() != tt.expectedError {
			t.Errorf("expected error %v, actual %v", tt.expectedError, actualError.Error())
		}
	}
}

func TestValidateDefaultFetchAllSystemMetrics(t *testing.T) {
	opts := CollectorOpts{}
	err := opts.validate()
	if err != nil {
		t.Errorf("expected nil error, actual %v", err)
	}
	if len(opts.SystemMetrics) != len(descSystemMetrics) {
		t.Errorf(
			"expected system metrics length to be %v, actual %v",
			len(descSystemMetrics),
			len(opts.SystemMetrics),
		)
	}
	for _, m := range opts.SystemMetrics {
		if _, found := descSystemMetrics[m]; !found {
			t.Errorf("unexpected system metric: %v", m)
		}
	}
}

func TestValidateDefaultFetchAllTubeMetrics(t *testing.T) {
	tests := []CollectorOpts{
		{Tubes: []string{"default"}},
		{AllTubes: true},
	}
	for _, opts := range tests {
		err := opts.validate()
		if err != nil {
			t.Errorf("expected nil error, actual %v", err)
		}
		if len(opts.TubeMetrics) != len(descTubeMetrics) {
			t.Errorf(
				"expected tube metrics length to be %v, actual %v",
				len(descTubeMetrics),
				len(opts.TubeMetrics),
			)
		}
		for _, m := range opts.TubeMetrics {
			if _, found := descTubeMetrics[m]; !found {
				t.Errorf("unexpected tube metric: %v", m)
			}
		}
	}
}

func TestNewBeanstalkdCollector(t *testing.T) {
	logger := mockLogger()
	beanstalkdServer, _ := beanstalkd.NewServer("localhost:11300", 10, 10)

	tests := []struct {
		beanstalkd                  *beanstalkd.Server
		opts                        CollectorOpts
		expectedError               error
		expectedSystemMetricsLength int
		expectedTubeMetricsLength   int
	}{
		// We expect validation to return errors.
		{
			beanstalkd:                  beanstalkdServer,
			opts:                        CollectorOpts{SystemMetrics: []string{"does_not_exist"}},
			expectedError:               fmt.Errorf("unknown system metric: does_not_exist"),
			expectedSystemMetricsLength: 0,
			expectedTubeMetricsLength:   0,
		},
		// We expect an initialised collector when there are no errors.
		{
			beanstalkd:                  beanstalkdServer,
			opts:                        CollectorOpts{},
			expectedError:               nil,
			expectedSystemMetricsLength: len(descSystemMetrics),
			expectedTubeMetricsLength:   0,
		},
		{
			beanstalkd:                  beanstalkdServer,
			opts:                        CollectorOpts{Tubes: []string{"default"}},
			expectedError:               nil,
			expectedSystemMetricsLength: len(descSystemMetrics),
			expectedTubeMetricsLength:   len(descTubeMetrics),
		},
		{
			beanstalkd:                  beanstalkdServer,
			opts:                        CollectorOpts{AllTubes: true},
			expectedError:               nil,
			expectedSystemMetricsLength: len(descSystemMetrics),
			expectedTubeMetricsLength:   len(descTubeMetrics),
		},
	}

	for _, tt := range tests {
		c, err := NewBeanstalkdCollector(tt.beanstalkd, tt.opts, logger)
		if !reflect.DeepEqual(tt.expectedError, err) {
			t.Errorf("expected error %v, actual %v", tt.expectedError, err)
		}
		if err != nil && c != nil {
			t.Error("expected nil collector because of error")
		}
		if tt.expectedError == nil && tt.expectedSystemMetricsLength != len(c.systemMetrics) {
			t.Errorf(
				"expected system metrics length %v, actual %v",
				tt.expectedSystemMetricsLength,
				len(c.systemMetrics),
			)
		}
		if tt.expectedError == nil && tt.expectedTubeMetricsLength != len(c.tubesMetrics) {
			t.Errorf(
				"expected tube metrics length %v, actual %v",
				tt.expectedTubeMetricsLength,
				len(c.tubesMetrics),
			)
		}
	}
}

func TestHealthyBeanstalkdServer(t *testing.T) {
	tests := []struct {
		allTubes           bool
		tubes              []string
		expectedNumMetrics int
	}{
		{
			allTubes:           false,
			tubes:              []string{"anotherTube"},
			expectedNumMetrics: 4, // 2 system metrics, 2 tube metrics (1 label)
		},
		{
			allTubes:           true,
			tubes:              nil,
			expectedNumMetrics: 6, // 2 system metrics, 4 tube metrics (2 + 2 labels)
		},
	}

	for _, tt := range tests {
		logger := mockLogger()
		collector, err := NewBeanstalkdCollector(
			mockHealthyBeanstalkd(),
			CollectorOpts{
				SystemMetrics: []string{"current_jobs_urgent_count", "current_jobs_ready_count"},
				AllTubes:      tt.allTubes,
				Tubes:         tt.tubes,
				TubeMetrics:   []string{"tube_current_jobs_urgent_count", "tube_current_jobs_ready_count"},
			},
			logger,
		)
		if err != nil {
			t.Errorf("expected nil error, actual %v", err)
		}

		ch := make(chan prometheus.Metric)

		go func() {
			defer close(ch)
			collector.Collect(ch)
		}()

		// "up" gauge
		if expected, actual := 1., readGauge((<-ch).(prometheus.Gauge)); expected != actual {
			t.Errorf("expected 'up' value %v, actual %v", expected, actual)
		}

		// "total scrapes" counter
		if expected, actual := 1., readCounter((<-ch).(prometheus.Counter)); expected != actual {
			t.Errorf("expected 'totalScrapes' value %v, actual %v", expected, actual)
		}

		// system metrics & tube metrics gauges
		actualTotal := 0
		for range ch {
			actualTotal++
		}
		if tt.expectedNumMetrics != actualTotal {
			t.Errorf("expected %d metrics, actual %d", tt.expectedNumMetrics, actualTotal)
		}
	}
}

func TestGetTubesToScrape(t *testing.T) {
	tests := []struct {
		opts          CollectorOpts
		beanstalkd    BeanstalkdServer
		expectedTubes []string
		expectedError error
	}{
		// We expect a healthy beanstalkd server to return all tubes
		// when we're fetching all tubes.
		{
			opts: CollectorOpts{
				AllTubes: true,
			},
			beanstalkd:    mockHealthyBeanstalkd(),
			expectedTubes: []string{"default", "anotherTube"},
			expectedError: nil,
		},
		// We expect only specific tubes when we're only interested
		// in specific tubes.
		{
			opts: CollectorOpts{
				Tubes: []string{"anotherTube"},
			},
			beanstalkd:    nil, // Not needed as we don't fetch the tubes.
			expectedTubes: []string{"anotherTube"},
			expectedError: nil,
		},
		// We expect an error when the beanstalkd server is not healthy
		{
			opts: CollectorOpts{
				AllTubes: true,
			},
			beanstalkd:    mockUnhealthyBeanstalkd(),
			expectedTubes: nil,
			expectedError: fmt.Errorf("list tubes error"),
		},
	}

	for _, tt := range tests {
		collector := BeanstalkdCollector{
			opts:       tt.opts,
			beanstalkd: tt.beanstalkd,
		}
		actualTubes, _ := collector.getTubesToScrape()
		if !reflect.DeepEqual(tt.expectedTubes, actualTubes) {
			t.Errorf("expected %v tubes, actual %v", tt.expectedTubes, actualTubes)
		}
	}
}

/********************     MOCKS     ********************/

type mockBeanstalkdServer struct {
	stats           beanstalkd.ServerStats
	statsError      error
	tubesStats      beanstalkd.ManyTubeStats
	tubesStatsError error
	listTubes       []string
	listTubesError  error
}

func (m *mockBeanstalkdServer) ListTubes() ([]string, error) {
	return m.listTubes, m.listTubesError
}

func (m *mockBeanstalkdServer) FetchStats() (beanstalkd.ServerStats, error) {
	return m.stats, m.statsError
}

func (m *mockBeanstalkdServer) FetchTubesStats(tubes map[string]bool) (beanstalkd.ManyTubeStats, error) {
	tubesStats := make(beanstalkd.ManyTubeStats, len(tubes))
	for tubeName := range tubes {
		tubesStats[tubeName] = m.tubesStats[tubeName]
	}
	return tubesStats, m.tubesStatsError
}

func mockHealthyBeanstalkd() *mockBeanstalkdServer {
	return &mockBeanstalkdServer{
		listTubes:      []string{"default", "anotherTube"},
		listTubesError: nil,
		stats: beanstalkd.ServerStats{
			"current-jobs-urgent": "10",
			"current-jobs-ready":  "20",
		},
		statsError: nil,
		tubesStats: beanstalkd.ManyTubeStats{
			"default": beanstalkd.TubeStatsOrError{
				Stats: beanstalkd.TubeStats{
					"current-jobs-urgent": "5",
					"current-jobs-ready":  "10",
				},
				Err: nil,
			},
			"anotherTube": beanstalkd.TubeStatsOrError{
				Stats: beanstalkd.TubeStats{
					"current-jobs-urgent": "1",
					"current-jobs-ready":  "2",
				},
				Err: nil,
			},
		},
		tubesStatsError: nil,
	}
}

func mockUnhealthyBeanstalkd() *mockBeanstalkdServer {
	return &mockBeanstalkdServer{
		listTubesError:  fmt.Errorf("list tubes error"),
		statsError:      fmt.Errorf("stats error"),
		tubesStatsError: fmt.Errorf("tubes stats error"),
	}
}

func mockLogger() *slog.Logger {
	var buff bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buff, nil))
	return logger
}
