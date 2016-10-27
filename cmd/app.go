package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Everything about apps",
}

var appCreateCmd = &cobra.Command{
	Use:   "create <app_name>",
	Short: "Creates an app",
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
	Short:   "List all apps",
	Long:    "Return all apps with address and team.",
	Example: "  $ teresa app list",
	RunE: func(cmd *cobra.Command, args []string) error {
		tc := NewTeresa()
		apps, err := tc.GetApps()
		if err != nil {
			if isNotFound(err) {
				return newCmdError("You have no apps")
			}
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

var appInfoCmd = &cobra.Command{
	Use:     "info <app_name>",
	Short:   "All infos about the app",
	Long:    "Return all infos about an specific app, like addresses, scale, auto scale, etc...",
	Example: "  $ teresa app info foo",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return newUsageError("You should provide the name of the app in order to continue")
		}
		appName := args[0]
		tc := NewTeresa()
		app, err := tc.GetAppInfo(appName)
		if err != nil {
			if isNotFound(err) {
				return newCmdError("App not found")
			}
			return err
		}

		color.New(color.FgCyan, color.Bold).Printf("[%s]\n", *app.Name)
		bold := color.New(color.Bold).SprintFunc()

		fmt.Println(bold("team:"), *app.Team)
		fmt.Println(bold("addresses:"))
		for _, a := range app.AddressList {
			fmt.Printf("  %s\n", a)
		}
		fmt.Println(bold("env vars:"))
		for _, e := range app.EnvVars {
			fmt.Printf("  %s=%s\n", *e.Key, *e.Value)
		}
		fmt.Println(bold("scale:"), app.Scale)
		fmt.Println(bold("autoscale:"))
		fmt.Printf("  %s %d%%\n", bold("cpu:"), *app.AutoScale.CPUTargetUtilization)
		fmt.Printf("  %s %d\n", bold("max:"), app.AutoScale.Max)
		fmt.Printf("  %s %d\n", bold("min:"), app.AutoScale.Min)
		fmt.Println(bold("limits:"))
		if len(app.Limits.Default) > 0 {
			fmt.Println(bold("  defaults"))
			for _, l := range app.Limits.Default {
				fmt.Printf("    %s %s\n", bold(*l.Resource), *l.Quantity)
			}
		}
		if len(app.Limits.DefaultRequest) > 0 {
			fmt.Println(bold("  request"))
			for _, l := range app.Limits.DefaultRequest {
				fmt.Printf("    %s %s\n", bold(*l.Resource), *l.Quantity)
			}
		}
		if app.HealthCheck != nil && (app.HealthCheck.Liveness != nil || app.HealthCheck.Readiness != nil) {
			fmt.Println(bold("healthcheck:"))
			if app.HealthCheck.Liveness != nil {
				fmt.Println(bold("  liveness:"))
				fmt.Printf("    %s %s\n", bold("path:"), app.HealthCheck.Liveness.Path)
				fmt.Printf("    %s %ds\n", bold("period:"), app.HealthCheck.Liveness.PeriodSeconds)
				fmt.Printf("    %s %ds\n", bold("timeout:"), app.HealthCheck.Liveness.TimeoutSeconds)
				fmt.Printf("    %s %ds\n", bold("initial delay:"), app.HealthCheck.Liveness.InitialDelaySeconds)
				fmt.Printf("    %s %d\n", bold("success threshold:"), app.HealthCheck.Liveness.SuccessThreshold)
				fmt.Printf("    %s %d\n", bold("failure threshold:"), app.HealthCheck.Liveness.FailureThreshold)
			}
			if app.HealthCheck.Readiness != nil {
				fmt.Println(bold("  readiness:"))
				fmt.Printf("    %s %s\n", bold("path:"), app.HealthCheck.Readiness.Path)
				fmt.Printf("    %s %ds\n", bold("period:"), app.HealthCheck.Readiness.PeriodSeconds)
				fmt.Printf("    %s %ds\n", bold("timeout:"), app.HealthCheck.Readiness.TimeoutSeconds)
				fmt.Printf("    %s %ds\n", bold("initial delay:"), app.HealthCheck.Readiness.InitialDelaySeconds)
				fmt.Printf("    %s %d\n", bold("success threshold:"), app.HealthCheck.Readiness.SuccessThreshold)
				fmt.Printf("    %s %d\n", bold("failure threshold:"), app.HealthCheck.Readiness.FailureThreshold)
			}
		}
		return nil
	},
}

func init() {
	// add AppCmd
	RootCmd.AddCommand(appCmd)
	// App commands
	appCmd.AddCommand(appCreateCmd)
	appCmd.AddCommand(appListCmd)
	appCmd.AddCommand(appInfoCmd)

	// appCreateCmd.Flags().StringVar(&teamNameFlag, "team", "", "team name")
	// appCmdCreate.Flags().IntVar(&appScaleFlag, "scale", 1, "replicas")

}
