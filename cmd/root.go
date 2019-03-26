package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "attacker",
	Short: "压测工具",
	Long:  `压测工具`,
}
