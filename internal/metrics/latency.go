package metrics

import (
	"sort"
	"sync"
)

type LatencyStore struct {
	mu   sync.RWMutex
	data map[string][]int64
}

type EndpointStats struct {
	Path         string `json:"path"`
	Count        int    `json:"count"`
	P50Millis    int64  `json:"p50_ms"`
	P95Millis    int64  `json:"p95_ms"`
	P99Millis    int64  `json:"p99_ms"`
	MaxMillis    int64  `json:"max_ms"`
	RecentSample int64  `json:"recent_ms"`
}

func NewLatencyStore() *LatencyStore {
	return &LatencyStore{data: map[string][]int64{}}
}

func (s *LatencyStore) Observe(path string, millis int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[path] = append(s.data[path], millis)
}

func (s *LatencyStore) Snapshot() []EndpointStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]EndpointStats, 0, len(s.data))
	for path, samples := range s.data {
		if len(samples) == 0 {
			continue
		}
		cp := append([]int64(nil), samples...)
		sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })
		out = append(out, EndpointStats{
			Path:         path,
			Count:        len(cp),
			P50Millis:    percentile(cp, 0.50),
			P95Millis:    percentile(cp, 0.95),
			P99Millis:    percentile(cp, 0.99),
			MaxMillis:    cp[len(cp)-1],
			RecentSample: samples[len(samples)-1],
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

func percentile(sorted []int64, p float64) int64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * p)
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
