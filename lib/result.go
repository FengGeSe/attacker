package lib

import (
	// "sort"
	"encoding/json"
	"fmt"
	"time"
)

type Result struct {
	Name       string             `json:"name" desc:"任务或操作名称"`
	Id         uint64             `json:"Id" desc:"每次任务的ID"`
	StartTime  time.Time          `json:"start_time" desc:"任务开始执行时间"`
	Latency    time.Duration      `json:"latency" desc:"耗时"`
	BytesOut   uint64             `json:"bytes_out"`
	BytesIn    uint64             `json:"bytes_in"`
	Code       uint16             `json:"code" desc:"返回码"`
	Error      string             `json:"error"  desc:"错误信息"`
	SubResults map[string]*Result `json:"sub_results" desc:"子操作的结果"`
}

func NewResult(name string) *Result {
	return &Result{
		Name: name,
	}
}

func (r *Result) String() string {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Sprintf("json序列化出错！Error: %v", err)
	}

	return string(data)
}

func (r *Result) SubResult(name string) *Result {
	subRst := &Result{
		Id:   r.Id,
		Name: name,
	}
	if r.SubResults == nil {
		r.SubResults = make(map[string]*Result)
	}
	r.SubResults[name] = subRst
	return subRst
}

func (r *Result) StartTiming() {
	r.StartTime = time.Now()
}
func (r *Result) EndTiming() {
	r.Latency = time.Now().Sub(r.StartTime)
}

//type Results []Result
//
//func (rs *Results) Add(r *Result) { *rs = append(*rs, *r) }
//
//// Close implements the Close method of the Report interface by sorting the Results
//func (rs *Results) Close() { sort.Sort(rs) }
//
//// The following methods implement sort.Interfac
//func (rs Results) Len() int           { return len(rs) }
//func (rs Results) Less(i, j int) bool { return rs[i].StartTime.Before(rs[j].EndTime) }
//func (rs Results) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
