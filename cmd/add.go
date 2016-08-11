package cmd

import "github.com/spf13/cobra"

var addCmd = &cobra.Command{
	Use: "add",
}

func init() {
	RootCmd.AddCommand(addCmd)
}
