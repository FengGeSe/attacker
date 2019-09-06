package lib

import ()

type Task interface {
	Init()        // 开始发压前只执行一次
	Run() *Result // 发压的逻辑
	Destroy()     // 结束发压后只执行一次
}
