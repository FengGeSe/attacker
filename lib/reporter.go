package lib

import (
	"fmt"
	"github.com/modood/table"
	"io"
	"time"
)

type TableReporter struct {
	Metrics *Metrics
}

// 处理结果输出到表格
func NewTableReporter() *TableReporter {
	return &TableReporter{
		Metrics: &Metrics{},
	}
}

func (this *TableReporter) Process(results <-chan *Result) {
	for rst := range results {
		this.Metrics.Add(rst)
	}
}

// 输出的格式
type Table struct {
	Task                    string
	Rate                    string
	Ratio                   string
	Mean                    string
	Max                     string
	Total, Success, Failure uint64
	P50, P95, P99           string
}

// 生成报告
func (this *TableReporter) Report(w io.Writer) error {
	this.Metrics.Close()

	data := MetricsToTable(this.Metrics)

	// 	s := table.Table(data)
	s := table.AsciiTable(data)
	_, err := w.Write([]byte(s + "\n"))

	return err
}

// metrics to tables
func MetricsToTable(m *Metrics) []Table {
	tables := []Table{}
	tb := Table{
		Task:    m.Name,
		Rate:    To2Float(m.Rate),
		Ratio:   ToPercent(m.Ratio),
		Mean:    FmtLatency(m.Latencies.Mean),
		Max:     FmtLatency(m.Latencies.Max),
		Total:   m.Total,
		Success: m.Success,
		Failure: m.Total - m.Success,
		P50:     FmtLatency(m.Latencies.P50),
		P95:     FmtLatency(m.Latencies.P95),
		P99:     FmtLatency(m.Latencies.P99),
	}

	tables = append(tables, tb)

	for _, v := range m.SubMetrics {
		tbs := MetricsToTable(v)
		tables = append(tables, tbs...)
	}
	return tables
}

func FmtLatency(latency time.Duration) string {
	return fmt.Sprintf("%.2fms", float64(latency)/1e6)
}

func To2Float(v float64) string {
	if v == float64(0) {
		return "-"
	}
	return fmt.Sprintf("%.2f", v)
}

func ToPercent(v float64) string {
	v = v * 100
	return fmt.Sprintf("%.1f%%", v)
}
