package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	attacker "github.com/FengGeSe/attacker/lib"
)

var AttackCmd = &cobra.Command{
	Use:   "attack",
	Short: "start attack",
	Long:  `start attack`,
}

func init() {
	RootCmd.AddCommand(AttackCmd)

	AttackCmd.Flags().IntP("rate", "r", 1, "攻击频率, 每秒")
	AttackCmd.Flags().DurationP("duration", "d", time.Second, "攻击时长, 10s=10秒钟 5m=5分钟")
	AttackCmd.Flags().StringP("file", "f", "./result.bin", "攻击结果输出指定文件")
	AttackCmd.Flags().BoolP("split", "s", false, "是否发出攻击和统计结果分开")
}

func Run(task attacker.Task, w io.Writer) {
	AttackCmd.Run = func(cmd *cobra.Command, args []string) {
		// 1. 读取参数
		duration, err := cmd.Flags().GetDuration("duration")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		rate, err := cmd.Flags().GetInt("rate")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		split, err := cmd.Flags().GetBool("split")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}

		// 3. 开始攻击
		now := time.Now().Format("2006-01-02 15:04:05")
		cmd.Printf("开始压测: %s , rate=%d/s, duration=%v\n", now, rate, duration)
		if split {
			attacker.RunOnly(task, w, rate, duration, file)
		} else {
			attacker.RunAndReport(task, w, rate, duration)
		}
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
