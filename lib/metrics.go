package lib

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/influxdata/tdigest"
)

type Metrics struct {
	Name       string              `json:"name" desc:"任务或操作名称"`
	Rate       float64             `json:"rate" desc:"请求速率每秒"`
	Ratio      float64             `json:"ratio" desc:"成功率"`
	Total      uint64              `json:"total" desc:"总共请求或操作数"`
	StartTime  time.Time           `json:"startTime" desc:"开始时间"`
	EndTime    time.Time           `json:"endTime" desc:"结束时间"`
	Success    uint64              `json:"success" desc:"成功的请求或操作数"`
	SubMetrics map[string]*Metrics `json:"sub-metrics" desc:"子操作的统计"`
	Errors     map[string]uint64   `json:"errors" desc:"错误集合"`

	Latencies LatencyMetrics `json:"latencies"`
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Add(r *Result) {
	// 错误统计
	Add(m, r)
}

// 递归处理result
func Add(m *Metrics, r *Result) string {
	// 公用处理
	m.Total++

	if m.Name == "" {
		m.Name = r.Name
	}

	// 统计子操作
	for k, v := range r.SubResults {
		if m.SubMetrics == nil {
			m.SubMetrics = make(map[string]*Metrics)
		}
		var subm *Metrics
		if v, ok := m.SubMetrics[k]; ok {
			subm = v
		} else {
			subm = m.SubMetric(k)
		}
		r.Error += Add(subm, v)
	}

	if r.Error != "" {
		if m.Errors == nil {
			m.Errors = make(map[string]uint64)
		}
		if _, ok := m.Errors[r.Error]; ok {
			m.Errors[r.Error]++
		} else {
			m.Errors[r.Error] = 1
		}
	} else {
		m.Success++
		m.Latencies.Add(r.Latency)
	}

	return r.Error
}

func AddError(m *Metrics, err string) {
}

func (m *Metrics) SubMetric(name string) *Metrics {
	if m.SubMetrics == nil {
		m.SubMetrics = make(map[string]*Metrics)
	}
	subm := &Metrics{}
	m.SubMetrics[name] = subm
	return subm
}

func (m *Metrics) String() string {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Sprintf("json序列化出错！Error: %v", err)
	}

	return string(data)
}

func (m *Metrics) Close() {
	// 请求速率，只有父任务的统计才有意义, 这里不考虑子任务的rate
	total := m.EndTime.Sub(m.StartTime).Seconds()
	// 精确到0.01, 四舍五入
	total = precision(total, 2, true)
	if total == 0 {
		total = 1
	}
	m.Rate = float64(m.Success) / total
	// 其他统计信息递归处理
	Close(m)
}

func precision(f float64, prec int, round bool) float64 {
	pow10_n := math.Pow10(prec)
	if round {
		return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
	}
	return math.Trunc((f)*pow10_n) / pow10_n
}

func Close(m *Metrics) {

	m.Latencies.Mean = time.Duration(float64(m.Latencies.Total) / float64(m.Success))
	m.Latencies.P50 = m.Latencies.Quantile(0.50)
	m.Latencies.P95 = m.Latencies.Quantile(0.95)
	m.Latencies.P99 = m.Latencies.Quantile(0.99)

	m.Ratio = float64((m.Success)) / float64(m.Total)
	for _, v := range m.SubMetrics {
		Close(v)
	}
}

// 用于判定返回结果是不错误的方法
// type Judge func(r *Result) bool

// LatencyMetrics holds computed request latency metrics.
type LatencyMetrics struct {
	// Total is the total latency sum of all requests in an attack.
	Total time.Duration `json:"total"`
	// Mean is the mean request latency.
	Mean time.Duration `json:"mean"`
	// P50 is the 50th percentile request latency.
	P50 time.Duration `json:"50th"`
	// P95 is the 95th percentile request latency.
	P95 time.Duration `json:"95th"`
	// P99 is the 99th percentile request latency.
	P99 time.Duration `json:"99th"`
	// Max is the maximum observed request latency.
	Max time.Duration `json:"max"`

	estimator estimator
}

// Add adds the given latency to the latency metrics.
func (l *LatencyMetrics) Add(latency time.Duration) {
	l.init()
	if l.Total += latency; latency > l.Max {
		l.Max = latency
	}
	l.estimator.Add(float64(latency))
}

// Quantile returns the nth quantile from the latency summary.
func (l LatencyMetrics) Quantile(nth float64) time.Duration {
	l.init()
	return time.Duration(l.estimator.Get(nth))
}

func (l *LatencyMetrics) init() {
	if l.estimator == nil {
		// This compression parameter value is the recommended value
		// for normal uses as per http://javadox.com/com.tdunning/t-digest/3.0/com/tdunning/math/stats/TDigest.html
		l.estimator = newTdigestEstimator(100)
	}
}

type estimator interface {
	Add(sample float64)
	Get(quantile float64) float64
}

type tdigestEstimator struct{ *tdigest.TDigest }

func newTdigestEstimator(compression float64) *tdigestEstimator {
	return &tdigestEstimator{TDigest: tdigest.NewWithCompression(compression)}
}

func (e *tdigestEstimator) Add(s float64) { e.TDigest.Add(s, 1) }
func (e *tdigestEstimator) Get(q float64) float64 {
	return e.TDigest.Quantile(q)
}
