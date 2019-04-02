package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	reporter "github.com/FengGeSe/attacker/lib"
)

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "从结果文件中生成报告",
	Long:  `从结果文件中生成报告`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 读取参数
		path, err := cmd.Flags().GetString("file")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		// 2. 开始生成报告
		now := time.Now().Format("2006-01-02 15:04:05")
		cmd.Printf("开始生成报告: %s \n", now)

		reporter := reporter.NewTableReporter()
		err = reporter.ReadAndReport(path, cmd.OutOrStdout())
		if err != nil {
			cmd.Printf("生成报告错误！%v", err)
		}

	},
}

func init() {
	RootCmd.AddCommand(ReportCmd)

	ReportCmd.Flags().StringP("file", "f", "./result.bin", "压测结果文件")
}
