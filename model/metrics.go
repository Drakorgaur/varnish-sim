package model

// CacheMetric is a struct for cache hit/miss metrics
type CacheMetric struct {
	hit  int
	miss int
}

// Hit increments the hit counter
func (m *CacheMetric) Hit() {
	m.hit++
}

// Miss increments the miss counter
func (m *CacheMetric) Miss() {
	m.miss++
}

// ExportType returns a map of cache hit/miss for exporting
// 	as these fields are private
func (m *CacheMetric) ExportType() map[string]float64 {
	total := m.hit + m.miss
	if total == 0 {
		total = 1
	}
	return map[string]float64{"hit": float64(m.hit), "miss": float64(m.miss), "total": float64(total), "hit_ratio": float64(m.hit) / float64(total)}
}

// RoutingMetric is a map showing traffic info that was routed to each backend
// Generic type V may be numeric type. Maybe either (M/G/..)Byte count or request count.
type RoutingMetric[T Numeric] map[WebInterface]T

// ExportType returns a map of routing metrics
func (m RoutingMetric[T]) ExportType() map[string]T {
	var result = make(map[string]T)
	for k, v := range m {
		result[k.String()] = v
	}

	return result
}
