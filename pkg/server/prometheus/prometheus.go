// Package prometheus is based on: https://github.com/banzaicloud/go-gin-prometheus
package prometheus

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

// Standard default metrics
// counter, counter_vec, gauge, gauge_vec,
// histogram, histogram_vec, summary, summary_vec
var reqCnt = &Metric{
	ID:          "reqCnt",
	Name:        "requests_total",
	Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	Type:        "counter_vec",
	Args:        []string{"code", "method", "host", "url"},
}

var reqDur = &Metric{
	ID:          "reqDur",
	Name:        "request_duration_seconds",
	Description: "The HTTP request latencies in seconds.",
	Type:        "summary",
}

var resSz = &Metric{
	ID:          "resSz",
	Name:        "response_size_bytes",
	Description: "The HTTP response sizes in bytes.",
	Type:        "summary",
}

var reqSz = &Metric{
	ID:          "reqSz",
	Name:        "request_size_bytes",
	Description: "The HTTP request sizes in bytes.",
	Type:        "summary",
}

var standardMetrics = []*Metric{
	reqCnt,
	reqDur,
	resSz,
	reqSz,
}

/*
RequestCounterURLLabelMappingFn is a function which can be supplied to the middleware to control
the cardinality of the request counter's "url" label, which might be required in some contexts.
For instance, if for a "/customer/:name" route you don't want to generate a time series for every
possible customer name, you could use this function:
func(c *gin.Context) string {
	url := c.Request.URL.String()
	for _, p := range c.Params {
		if s.Key == "name" {
			url = strings.Replace(url, s.Value, ":name", 1)
			break
		}
	}
	return url
}
which would map "/customer/alice" and "/customer/bob" to their template "/customer/:name".
*/
type RequestCounterURLLabelMappingFn func(c *gin.Context) string

// Metric is a definition for the name, description, type, ID, and
// prometheus.Collector type (i.e. CounterVec, Summary, etc) of each metric.
type Metric struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
}

// Server represents the implementation of Prometheus server object.
// It contains the metrics gathered by the instance and its path.
type Server struct {
	reqCnt               *prometheus.CounterVec
	reqDur, reqSz, resSz prometheus.Summary
	PushGateway          PushGateway

	MetricsList []*Metric

	ReqCntURLLabelMappingFn RequestCounterURLLabelMappingFn

	http   *http.Server
	Router *gin.Engine
	Port   string
}

// PushGateway contains the configuration for pushing to a Prometheus pushgateway (optional).
type PushGateway struct {
	// Pushgateway interval in seconds.
	IntervalSeconds time.Duration

	// Pushgateway URL in format http://domain:port
	// where JOBNAME can be any string of your choice.
	GatewayURL string

	// Local metrics URL where metrics are fetched from, this could be ommited in the future
	// if implemented using prometheus common/expfmt instead.
	MetricsURL string

	// Pushgateway job name, defaults to "gin".
	Job string
}

// NewServer creates new Prometheus server. It generates a new set of metrics
// with a certain subsystem name.
func NewServer(port, subsystem string, expandedParams []string) *Server {
	var metricsList []*Metric

	for _, metric := range standardMetrics {
		metricsList = append(metricsList, metric)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	expandedParams = append(expandedParams, "code")

	s := &Server{
		MetricsList: metricsList,
		ReqCntURLLabelMappingFn: func(c *gin.Context) string {
			url := c.Request.URL.EscapedPath() // i.e. by default do nothing, i.e. return URL as is
			for _, p := range c.Params {
				if contains(expandedParams, p.Key) {
					continue
				}

				// Overcome wildcard (*path) matching issue, which takes the beginning slash as well.
				value := strings.TrimPrefix(p.Value, "/")

				url = strings.Replace(url, value, ":"+p.Key, 1)
			}
			return url
		},
		http: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: router,
		},
		Router: router,
		Port:   port,
	}

	s.registerMetrics("subsystem")

	return s
}

func contains(slice []string, s string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}
	return false
}

// Start starts the Prometheus server.
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("Start Prometheus server on port %s", s.Port))
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops the Prometheus server.
func (s *Server) Stop(ctx context.Context) error {
	log.Info(fmt.Sprintf("Stop Prometheus server on port %s", s.Port))
	if err := s.http.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// SetPushGateway sends metrics to a remote pushgateway exposed on pushGatewayURL
// every pushIntervalSeconds. Metrics are fetched from metricsURL.
func (s *Server) SetPushGateway(pushGatewayURL, metricsURL string, pushIntervalSeconds time.Duration) {
	s.PushGateway.GatewayURL = pushGatewayURL
	s.PushGateway.MetricsURL = metricsURL
	s.PushGateway.IntervalSeconds = pushIntervalSeconds
	s.startPushTicker()
}

// SetPushGatewayJob job name, defaults to "gin"
func (s *Server) SetPushGatewayJob(j string) {
	s.PushGateway.Job = j
}

func (s *Server) setMetricsPath(e *gin.Engine, metricsPath string) {
	s.Router.GET(metricsPath, prometheusHandler())
}

func (s *Server) setMetricsPathWithAuth(e *gin.Engine, accounts gin.Accounts, metricsPath string) {
	s.Router.GET(metricsPath, gin.BasicAuth(accounts), prometheusHandler())
}

func (s *Server) getMetrics() []byte {
	response, err := http.Get(s.PushGateway.MetricsURL)
	if err != nil {
		log.Error(err, "Failed getting Prometheus metrics")
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Error(err, "Failed closing Prometheus response body")
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err, "Failed reading Prometheus response body")
	}

	return body
}

func (s *Server) getPushGatewayURL() string {
	h, _ := os.Hostname()
	if s.PushGateway.Job == "" {
		s.PushGateway.Job = "gin"
	}
	return s.PushGateway.GatewayURL + "/metrics/job/" + s.PushGateway.Job + "/instance/" + h
}

func (s *Server) sendMetricsToPushGateway(metrics []byte) {
	req, err := http.NewRequest("POST", s.getPushGatewayURL(), bytes.NewBuffer(metrics))
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Error(err, "Failed sending Prometheus metrics to Pushgateway")
	}
}

func (s *Server) startPushTicker() {
	ticker := time.NewTicker(time.Second * s.PushGateway.IntervalSeconds)
	go func() {
		for range ticker.C {
			s.sendMetricsToPushGateway(s.getMetrics())
		}
	}()
}

// NewMetric associates prometheus.Collector based on Metric.Type.
func NewMetric(m *Metric, subsystem string) prometheus.Collector {
	var metric prometheus.Collector
	switch m.Type {
	case "counter_vec":
		metric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "counter":
		metric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "gauge_vec":
		metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "gauge":
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "histogram_vec":
		metric = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "histogram":
		metric = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "summary_vec":
		metric = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "summary":
		metric = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	}
	return metric
}

func (s *Server) registerMetrics(subsystem string) {
	for _, metricDef := range s.MetricsList {
		metric := NewMetric(metricDef, subsystem)
		if err := prometheus.Register(metric); err != nil {
			log.Error(err, fmt.Sprintf("%s could not be registered in Prometheus", metricDef.Name))
		}
		switch metricDef {
		case reqCnt:
			s.reqCnt = metric.(*prometheus.CounterVec)
		case reqDur:
			s.reqDur = metric.(prometheus.Summary)
		case resSz:
			s.resSz = metric.(prometheus.Summary)
		case reqSz:
			s.reqSz = metric.(prometheus.Summary)
		}
		metricDef.MetricCollector = metric
	}
}

// Use adds the middleware to a gin engine.
func (s *Server) Use(e *gin.Engine, metricsPath string) {
	e.Use(s.handlerFunc(metricsPath))
	s.setMetricsPath(e, metricsPath)
}

// UseWithCustomMetrics adds the middleware to a gin engine with custom metrics.
func (s *Server) UseWithCustomMetrics(e *gin.Engine, gatherer prometheus.Gatherers, metricsPath string) {
	s.setMetricsPathWithCustomMetrics(e, gatherer, metricsPath)
}

func (s *Server) setMetricsPathWithCustomMetrics(e *gin.Engine, gatherer prometheus.Gatherers, metricsPath string) {
	s.Router.GET(metricsPath, prometheusHandlerFor(gatherer))
}

// UseWithAuth adds the middleware to a gin engine with BasicAuth.
func (s *Server) UseWithAuth(e *gin.Engine, accounts gin.Accounts, metricsPath string) {
	e.Use(s.handlerFunc(metricsPath))
	s.setMetricsPathWithAuth(e, accounts, metricsPath)
}

func (s *Server) handlerFunc(metricsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.String() == metricsPath {
			c.Next()
			return
		}

		start := time.Now()
		reqSz := computeApproximateRequestSize(c.Request)

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		resSz := float64(c.Writer.Size())

		s.reqDur.Observe(elapsed)
		url := s.ReqCntURLLabelMappingFn(c)
		s.reqCnt.WithLabelValues(status, c.Request.Method, c.Request.Host, url).Inc()
		s.reqSz.Observe(float64(reqSz))
		s.resSz.Observe(resSz)
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func prometheusHandlerFor(gatherer prometheus.Gatherers) gin.HandlerFunc {
	h := promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
