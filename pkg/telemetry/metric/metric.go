package metric

type metadata struct {
	Name        string
	Description string
}

var (
	// Latency in seconds buckets for histogram metric (5ms to 10s).
	durationBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	// Size in bytes buckets for histogram metrics (100B to 1GB).
	sizeBuckets = []float64{100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000}
)
