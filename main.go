package main

import (
	"fmt"
	attacker "github.com/fenggese/attacker/lib"
	"os"
	"time"
)

// implements attacker.Task
type myTask struct{}

func (t *myTask) Init() {}

var count = 0

func (t *myTask) Run() *attacker.Result {
	rst := attacker.NewResult("流程压测")
	// A
	rstA := rst.SubResult("操作A")
	rstA.StartTiming()
	time.Sleep(1 * time.Millisecond)
	rstA.EndTiming()

	// B
	rstB := rst.SubResult("操作B")
	rstB.StartTiming()
	time.Sleep(1 * time.Millisecond)
	rstB.EndTiming()

	if count == 1 {
		rstB.Code = 500
		rstB.Error = "出错了偶"
		count++
		return rst
	}

	// C
	rstC := rst.SubResult("操作C")
	rstC.StartTiming()
	time.Sleep(1 * time.Millisecond)
	rstC.EndTiming()
	count++

	return rst
}

func (t *myTask) Destroy() {
}

func main() {

	task := &myTask{}

	results := attacker.Run(task, 1, 3*time.Second)
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("开始压测: ", now)
	reporter := attacker.NewTableReporter()

	reporter.Process(results)

	reporter.Report(os.Stdout)

}
