package main

import (
	"os"
	"time"

	attacker "github.com/FengGeSe/attacker/lib"
)

// implements attacker.Task
type myTask struct{}

var _ attacker.Task = &myTask{}

var count = 0

func (t *myTask) Run() *attacker.Result {
	rst := attacker.NewResult("流程压测")
	// A
	{
		subRst := rst.SubResult("操作A")
		subRst.StartTiming()
		// do something
		time.Sleep(1 * time.Millisecond)
		subRst.EndTiming()
	}

	// B
	{
		subRst := rst.SubResult("操作B")
		subRst.StartTiming()
		// do something
		time.Sleep(1 * time.Millisecond)
		subRst.EndTiming()

		// mock error
		if count == 1 {
			subRst.Code = 500
			subRst.Error = "出错了偶"
			count++
			return rst
		}
		count++
	}

	// C
	{
		subRst := rst.SubResult("操作C")
		subRst.StartTiming()
		// do something
		time.Sleep(1 * time.Millisecond)
		subRst.EndTiming()
	}

	return rst
}

func main() {
	task := &myTask{}

	attacker.RunAndReport(task, os.Stdout, 10, 2*time.Second)
}
