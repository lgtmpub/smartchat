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
	counterM = make(map[metric]counterItem)
	counterL sync.Mutex
)

type metric struct {
	Subsystem string
	Name      string
	Desc      string
}

type counterItem struct {
	counter *counter
	labels  string
}

type counter struct {
	vec *prometheus.CounterVec
}

// NewCounter creates a new counter metric
func NewCounter(subsystem, name, desc string, labels []string) (*counter, error) {
	if subsystem == "" || name == "" {
		return nil, fmt.Errorf("subsystem and name are required")
	}
	if len(labels) == 0 {
		return nil, fmt.Errorf("labels are required")
	}
	slices.Sort(labels)
	labStr := strings.Join(labels, ",")

	if desc == "" {
		desc = fmt.Sprintf("%s %s Counter", subsystem, name)
	}
	m := metric{Subsystem: subsystem, Name: name, Desc: desc}
	counterL.Lock()
	defer counterL.Unlock()
	item, ok := counterM[m]
	if ok {
		if item.labels == labStr {
			return item.counter, nil
		}
		return nil, errors.New("metric already exists with different labels")
	}

	item = counterItem{
		counter: &counter{vec: prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      name,
			Help:      desc,
		}, labels)},
		labels: labStr,
	}

	prometheus.MustRegister(item.counter.vec)
	counterM[m] = item
	return item.counter, nil
}

// Inc increments the counter of the operation
func (c *counter) Inc(labels map[string]string) {
	c.vec.With(labels).Inc()
}

// Add adds the value to the counter of the operation
func (c *counter) Add(labels map[string]string, value float64) {
	c.vec.With(labels).Add(value)
}
