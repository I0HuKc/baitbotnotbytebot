package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "r",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("BaitbotV2")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
