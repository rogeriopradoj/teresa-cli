package cmd

import "github.com/spf13/cobra"

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add users to a team",
}

func init() {
	RootCmd.AddCommand(addCmd)
}
