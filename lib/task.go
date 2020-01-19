package lib

import ()

type Task interface {
	Run() *Result // 发压的逻辑
}
