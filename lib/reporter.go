package lib

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/modood/table"
)

type TableReporter struct {
	Metrics *Metrics
}

// 处理结果输出到表格
func NewTableReporter() *TableReporter {
	return &TableReporter{
		Metrics: NewMetrics(),
	}
}

func (this *TableReporter) Process(results <-chan *Result) {
	this.Metrics.StartTime = time.Now()
	for rst := range results {
		this.Metrics.Add(rst)
	}
	this.Metrics.EndTime = time.Now()
}

// 写入文件
func (this *TableReporter) Write(path string, results <-chan *Result) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		panic(err)
	}
	defer func() {
		f.Close()
	}()
	encoder := gob.NewEncoder(f)
	for rst := range results {
		encoder.Encode(rst)
	}
}

// 读取结果文件并产生报告
func (this *TableReporter) ReadAndReport(path string, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
	}()

	decoder := gob.NewDecoder(f)

	for {
		var result Result
		err := decoder.Decode(&result)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		this.Metrics.Add(&result)
	}

	this.Report(w)
	return nil
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

	// 结果表格化
	// 	s := table.Table(data)
	s := table.AsciiTable(data)
	_, err := w.Write([]byte(s + "\n"))

	// 错误结果
	var errStr = "Error Set: \n"
	var count uint64
	for k, v := range this.Metrics.Errors {
		count++
		errStr += fmt.Sprintf("[%d]\t错误: %s\n\t数量: %d\n", count, k, v)
	}
	w.Write([]byte(errStr))
	w.Write([]byte("\n"))

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

	if len(m.SubMetrics) == 0 {
		return tables
	}

	// 子任务排序
	subMetricsNames := sort.StringSlice{}
	for k, _ := range m.SubMetrics {
		subMetricsNames = append(subMetricsNames, k)
	}
	subMetricsNames.Sort()

	for _, n := range subMetricsNames {
		tbs := MetricsToTable(m.SubMetrics[n])
		tables = append(tables, tbs...)
	}
	return tables
}

func FmtLatency(latency time.Duration) string {
	if latency <= 0 {
		return fmt.Sprintf("0 ms")
	}
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
