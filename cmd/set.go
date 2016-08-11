package cmd

import "github.com/spf13/cobra"

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set properties for an app.",
}

func init() {
	RootCmd.AddCommand(setCmd)
}
