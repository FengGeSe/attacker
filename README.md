# Attacker

attacker是一个支持自定义任务的项目。模仿了[vegeta](https://github.com/tsenart/vegeta)的压测项目，vegeta主要用于http请求的压测。  attacker模仿了vegeta的压测逻辑，支持自定义任务。可以按照指定的压测周期，压测频率执行自定义任务。

### Hello world

```go
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
	cmd.Run(task, os.Stdout)
}
```



使用：

```go run main.go attack -d 1s -r 10```

```bash
开始压测: 2020-02-05 21:01:59 , rate=10/s, duration=1s
+--------------+------+--------+--------+--------+-------+---------+---------+--------+--------+--------+
| Task         | Rate | Ratio  | Mean   | Max    | Total | Success | Failure | P50    | P95    | P99    |
+--------------+------+--------+--------+--------+-------+---------+---------+--------+--------+--------+
| 流程压测     | 9.89 | 90.0%  | 4.04ms | 4.12ms | 10    | 9       | 1       | 4.04ms | 4.12ms | 4.12ms |
|  1.操作A     | -    | 100.0% | 1.32ms | 1.37ms | 10    | 10      | 0       | 1.33ms | 1.37ms | 1.37ms |
|  2.操作B     | -    | 90.0%  | 1.34ms | 1.36ms | 10    | 9       | 1       | 1.33ms | 1.36ms | 1.36ms |
|  3.操作C     | -    | 100.0% | 1.32ms | 1.37ms | 9     | 9       | 0       | 1.32ms | 1.37ms | 1.37ms |
+--------------+------+--------+--------+--------+-------+---------+---------+--------+--------+--------+
Error Set:
[1]	错误: 出错了偶
	数量: 1
```



