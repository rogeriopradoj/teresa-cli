package cmd

import "github.com/spf13/cobra"

// createCmd represents the create command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get app details",
	Long: `Get application details or get teams.

To get details about an application:

	$ teresa get app --app my_app_name --team my_team

To get the teams you belong to:

	$ teresa get teams
	`,
}

func init() {
	RootCmd.AddCommand(getCmd)
}
