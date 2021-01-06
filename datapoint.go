package tsdbh

import (
	"fmt"
)

// DataPoint struct represents one data point for tsdb
type DataPoint struct {
	Metric   string  `json:"metric" `
	Unixtime int     `json:"timestamp"`
	Value    float64 `json:"value"`
	Tags     Tags    `json:"tags"`
}

// NewDataPoint takes data point params and return a data point struct
func NewDataPoint(metric string, unixtime int, value float64, tags []Tag) DataPoint {
	return DataPoint{metric, unixtime, value, tags}
}

// String would return a string with data point values
func (d DataPoint) String() string {
	s := fmt.Sprintf("%s %d %f", d.Metric, d.Unixtime, d.Value)
	for _, t := range d.Tags {
		s += " " + t.String()
	}
	return s
}
