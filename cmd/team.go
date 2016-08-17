package cmd

import (
	"fmt"

	_ "github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

// create a team
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Create a team",
	Long: `Create a team that can have many applications.

eg.:

	$ teresa create team --email sitedev@mydomain.com --name site --url sitedev.mydomain.com
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if teamNameFlag == "" {
			Usage(cmd)
			return
		}
		tc := NewTeresa()
		team, err := tc.CreateTeam(teamNameFlag, teamEmailFlag, teamURLFlag)
		if err != nil {
			log.Fatalf("Failed to create team: %s", err)
		}
		log.Infof("Team created. Name: %s Email: %s URL: %s\n", *team.Name, team.Email, team.URL)
	},
}

// delete team
var deleteTeamCmd = &cobra.Command{
	Use:   "team",
	Short: "Delete a team",
	Long:  `Delete a team`,
	Run: func(cmd *cobra.Command, args []string) {
		if teamIDFlag == 0 {
			Fatalf(cmd, "team ID is required")
		}
		if err := NewTeresa().DeleteTeam(teamIDFlag); err != nil {
			log.Fatalf("Failed to delete team: %s", err)
		}
		log.Infof("Team deleted.")
	},
}

var getTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Get teams",
	// Long:  `Delete a team`,
	Run: func(cmd *cobra.Command, args []string) {
		teams, err := NewTeresa().GetTeams()
		if err != nil {
			log.Fatalf("Failed to retrieve teams: %s", err)
		}

		fmt.Println("\nTeams:")
		for _, t := range teams {
			if t.IAmMember {
				fmt.Printf("- %s (member)\n", *t.Name)
			} else {
				fmt.Printf("- %s\n", *t.Name)
			}
			if t.Email != "" {
				fmt.Printf("  contact: %s\n", t.Email)
			}
			if t.URL != "" {
				fmt.Printf("  url: %s\n", t.URL)
			}
		}
		fmt.Println("")
	},
}

var addUserToTeamCmd = &cobra.Command{
	Use:   "team-user",
	Short: "Add user to team",
	Long: `Add a user to team.

You can add a new user as a member of a team with:

	$ teresa add team-user --email john.doe@teresa.com --team my-team

You need to create a user before use this command.
If the user already is member of the team, you will get an error.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if teamNameFlag == "" {
			Fatalf(cmd, "team name is required")
		}
		if userEmailFlag == "" {
			Fatalf(cmd, "user e-mail is required")
		}
		tc := NewTeresa()
		err := tc.AddUserToTeam(teamNameFlag, userEmailFlag)
		if err == nil {
			log.Infof("user [%s] is now member of the team [%s]", userEmailFlag, teamNameFlag)
			return
		}
		if err.Code() == 500 {
			log.Fatalf("Error with the command")
		}
		if err.Code() == 422 {
			log.Fatal(err.Payload.Message)
		}
	},
}

func init() {
	createCmd.AddCommand(teamCmd)
	teamCmd.Flags().StringVarP(&teamNameFlag, "name", "n", "", "team name [required]")
	teamCmd.Flags().StringVarP(&teamEmailFlag, "email", "e", "", "team email, if any")
	teamCmd.Flags().StringVarP(&teamURLFlag, "url", "u", "", "team site's URL, if any")

	deleteCmd.AddCommand(deleteTeamCmd)
	deleteTeamCmd.Flags().Int64Var(&teamIDFlag, "id", 0, "team ID [required]")

	getCmd.AddCommand(getTeamsCmd)

	// add user to team
	addCmd.AddCommand(addUserToTeamCmd)
	addUserToTeamCmd.Flags().StringVar(&userEmailFlag, "email", "", "user email [required]")
	addUserToTeamCmd.Flags().StringVar(&teamNameFlag, "team", "", "team name [required]")
}
