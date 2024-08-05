package store

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	histogramM = make(map[metric]histogramItem)
	histogramL sync.Mutex
)

type histogramItem struct {
	histogram *histogram
	labels    string
}

type histogram struct {
	vec *prometheus.HistogramVec
}

// NewHistogram creates a new histogram metric
func NewHistogram(subsystem, name, desc string, labels []string) (*histogram, error) {
	if subsystem == "" || name == "" {
		return nil, fmt.Errorf("subsystem and name are required")
	}
	if len(labels) == 0 {
		return nil, fmt.Errorf("labels are required")
	}
	slices.Sort(labels)
	labStr := strings.Join(labels, ",")

	if desc == "" {
		desc = fmt.Sprintf("%s %s Histogram", subsystem, name)
	}
	m := metric{Subsystem: subsystem, Name: name, Desc: desc}
	histogramL.Lock()
	defer histogramL.Unlock()
	item, ok := histogramM[m]
	if ok {
		if item.labels == labStr {
			return item.histogram, nil
		}
		return nil, errors.New("metric already exists with different labels")
	}

	item = histogramItem{
		histogram: &histogram{
			vec: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Subsystem: subsystem,
					Name:      name,
					Help:      desc,
				},
				labels)},
		labels: labStr,
	}
	histogramM[m] = item
	return item.histogram, nil
}

// Observe records a value to the histogram
func (h *histogram) Observe(labels map[string]string, value float64) {
	h.vec.With(labels).Observe(value)
}
