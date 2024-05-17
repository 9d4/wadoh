package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(tmpCmd)
}

var tmpCmd = &cobra.Command{
	Use:   "tmp",
	Short: "temp command for experiment",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
