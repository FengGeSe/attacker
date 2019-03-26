package lib

import (
	"io"
	"sync"
	"time"
)

// 发起进攻
// task 任务
// rate 速率每秒
// du 压测时长

var workers = uint64(10)

func Run(task Task, rate int, du time.Duration) <-chan *Result {

	var wg sync.WaitGroup
	ticks := make(chan uint64)
	results := make(chan *Result)
	for i := uint64(0); i < workers; i++ {
		wg.Add(1)
		go attack(task, &wg, ticks, results)
	}

	go func() {
		defer close(results)
		defer wg.Wait()
		defer close(ticks)
		interval := uint64(time.Second.Nanoseconds() / int64(rate))
		hits := uint64(du) / interval
		began, count := time.Now(), uint64(0)
		for {
			now, next := time.Now(), began.Add(time.Duration(count*interval))
			time.Sleep(next.Sub(now))
			select {
			case ticks <- count:
				if count++; count == hits {
					return
				}
			default: // all workers are blocked. start one more and try again
				wg.Add(1)
				go attack(task, &wg, ticks, results)
			}
		}

	}()
	return results
}

func RunAndReport(task Task, w io.Writer, rate int, du time.Duration) {
	reporter := NewTableReporter()
	results := Run(task, rate, du)
	reporter.Process(results)
	reporter.Report(w)
}

func attack(task Task, wg *sync.WaitGroup, ticks <-chan uint64, results chan<- *Result) {
	defer wg.Done()
	for t := range ticks {
		task.Init()

		startTime := time.Now()
		rst := task.Run()
		rst.StartTime = startTime
		rst.EndTiming()

		rst.Id = t

		results <- rst
		task.Destroy()
	}
}
