package lib

import ()

type Task interface {
	Init()
	Run() *Result
	Destroy()
}
