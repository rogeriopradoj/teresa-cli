package cmd

import "github.com/spf13/cobra"

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset properties for an app.",
}

func init() {
	RootCmd.AddCommand(unsetCmd)
}
