package main

import (
	"os"
	"time"

	attacker "github.com/FengGeSe/attacker/lib"

	cmd "github.com/FengGeSe/attacker/cmd"
)

// implements attacker.Task
type myTask struct{}

var _ attacker.Task = &myTask{}

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

func main() {

	task := &myTask{}

	cmd.Run(task, os.Stdout)
}
