package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "attacker",
	Short: "Attacker压测工具",
	Long:  `Attacker压测工具`,
}
