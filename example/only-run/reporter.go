package main

import (
	"fmt"
	"os"

	reporter "github.com/FengGeSe/attacker/lib"
)

func main() {

	path := "./result.out"
	reporter := reporter.NewTableReporter()
	err := reporter.ReadAndReport(path, os.Stdout)
	if err != nil {
		fmt.Printf("生成报告错误！%v \n", err)
	}
}
