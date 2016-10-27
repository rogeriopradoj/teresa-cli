package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var appCmd = &cobra.Command{
	Use: "app",
}

var appCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create an app",
	Long: `Creates a new application.

The application name is always required, but team name is only required if you
are part of more than one team.`,
	Example: `  teresa app create foo

  teresa app create foo --team bar

  You can optionally specify the number of containers/pods you want your app
  to run with the --scale (defaults to 1) option:

  teresa app create foo --team bar --scale 4`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("=> %+v\n", "Entrou aqui")

		// if len(args) == 0 {
		// 	Usage(cmd)
		// 	return
		// }
		// if appScaleFlag == 0 {
		// 	Fatalf(cmd, "at least one replica is required")
		// }
		//
		// tc := NewTeresa()
		// teamID := tc.GetTeamID(teamNameFlag)
		// app, err := tc.CreateApp(args[0], int64(appScaleFlag), teamID)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// log.Infof("App created. Name: %s Replicas: %d", *app.Name, *app.Scale)
	},
}

var appListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get apps",
	Long:    "Return all apps with address and team.",
	Example: "  $ teresa app list",
	RunE: func(cmd *cobra.Command, args []string) error {
		tc := NewTeresa()
		apps, err := tc.GetApps()
		if err != nil {
			return err
		}
		// rendering app info in a table view
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"TEAM", "APP", "ADDRESS"})
		table.SetRowLine(true)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetRowSeparator("-")
		table.SetAutoWrapText(false)
		for _, a := range apps {
			r := []string{*a.Team, *a.Name, a.AddressList[0]}
			table.Append(r)
		}
		table.Render()
		return nil
	},
}

func init() {
	// add AppCmd
	RootCmd.AddCommand(appCmd)

	// App commands
	appCmd.AddCommand(appCreateCmd)
	appCmd.AddCommand(appListCmd)

	appCreateCmd.Flags().StringVar(&teamNameFlag, "team", "", "team name")
	// appCmdCreate.Flags().IntVar(&appScaleFlag, "scale", 1, "replicas")

}
