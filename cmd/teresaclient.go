package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/luizalabs/teresa-api/client"
	"github.com/luizalabs/teresa-api/client/apps"
	"github.com/luizalabs/teresa-api/client/auth"
	"github.com/luizalabs/teresa-api/client/deployments"
	"github.com/luizalabs/teresa-api/client/teams"
	"github.com/luizalabs/teresa-api/client/users"
	"github.com/luizalabs/teresa-api/models"
	_ "github.com/prometheus/common/log" // still needed?
)

// TeresaClient foo bar
type TeresaClient struct {
	teresa         *apiclient.Teresa
	apiKeyAuthFunc runtime.ClientAuthInfoWriter
}

// TeresaServer scheme and host where the api server is running
type TeresaServer struct {
	scheme string
	host   string
}

// AppInfo foo bar
type AppInfo struct {
	AppID  int64
	TeamID int64
}

// ParseServerURL parse the server url and ensure its in the format we expect
func ParseServerURL(s string) (TeresaServer, error) {
	u, err := url.Parse(s)
	if err != nil {
		return TeresaServer{}, fmt.Errorf("Failed to parse server: %+v\n", err)
	}
	if u.Scheme == "" || u.Scheme != "http" && u.Scheme != "https" {
		return TeresaServer{}, errors.New("accepted server url format: http(s)://hostname[:port]")
	}
	ts := TeresaServer{scheme: u.Scheme, host: u.Host}
	return ts, nil
}

// NewTeresa foo bar
func NewTeresa() TeresaClient {
	cfg, err := readConfigFile(cfgFile)
	if err != nil {
		log.Fatalf("Failed to read config file, err: %+v\n", err)
	}
	n, err := getCurrentClusterName()
	if err != nil {
		log.Fatalf("Failed to get current cluster name, err: %+v\n", err)
	}
	cluster := cfg.Clusters[n]
	suffix := "/v1"

	tc := TeresaClient{teresa: apiclient.Default}
	log.Debugf(`Setting new teresa client. server: %s, api suffix: %s`, cluster.Server, suffix)

	ts, err := ParseServerURL(cluster.Server)
	if err != nil {
		log.Fatal(err)
	}

	client.DefaultTimeout = 5 * time.Minute // 5 minutes to wait deploy proccess
	c := client.New(ts.host, suffix, []string{ts.scheme})

	tc.teresa.SetTransport(c)

	if cluster.Token != "" {
		tc.apiKeyAuthFunc = httptransport.APIKeyAuth("Authorization", "header", cluster.Token)
	}
	return tc
}

// Login login the user
func (tc TeresaClient) Login(email strfmt.Email, password strfmt.Password) (token string, err error) {
	params := auth.NewUserLoginParams()
	params.WithBody(&models.Login{Email: &email, Password: &password})

	r, err := tc.teresa.Auth.UserLogin(params)
	if err != nil {
		return "", err
	}
	return r.Payload.Token, nil
}

// CreateTeam Creates a team
func (tc TeresaClient) CreateTeam(name, email, URL string) (*models.Team, error) {
	params := teams.NewCreateTeamParams()
	e := strfmt.Email(email)
	params.WithBody(&models.Team{Name: &name, Email: e, URL: URL})
	r, err := tc.teresa.Teams.CreateTeam(params, tc.apiKeyAuthFunc)
	if err != nil {
		return nil, err
	}
	return r.Payload, nil
}

// DeleteTeam Deletes a team
func (tc TeresaClient) DeleteTeam(ID int64) error {
	params := teams.NewDeleteTeamParams()
	params.TeamID = ID
	_, err := tc.teresa.Teams.DeleteTeam(params, tc.apiKeyAuthFunc)
	return err
}

// CreateApp creates an user
func (tc TeresaClient) CreateApp(name string, scale int64, teamID int64) (app *models.App, err error) {
	params := apps.NewCreateAppParams()
	params.TeamID = teamID
	params.WithBody(&models.App{Name: &name, Scale: &scale})
	r, err := tc.teresa.Apps.CreateApp(params, tc.apiKeyAuthFunc)
	if err != nil {
		return nil, err
	}
	return r.Payload, nil
}

// GetApps return apps for a specific team
func (tc TeresaClient) GetApps(teamID int64) (app []*models.App, err error) {
	params := apps.NewGetAppsParams().WithTeamID(teamID)
	r, err := tc.teresa.Apps.GetApps(params, tc.apiKeyAuthFunc)
	if err != nil {
		return nil, err
	}
	return r.Payload.Items, nil
}

// GetAppDetail Create app attributes
func (tc TeresaClient) GetAppDetail(teamID, appID int64) (app *models.App, err error) {
	params := apps.NewGetAppDetailsParams().WithTeamID(teamID).WithAppID(appID)
	r, err := tc.teresa.Apps.GetAppDetails(params, tc.apiKeyAuthFunc)
	if err != nil {
		return nil, err
	}
	return r.Payload, nil
}

// CreateUser Create an user
func (tc TeresaClient) CreateUser(name, email, password string, isAdmin bool) (user *models.User, err error) {
	params := users.NewCreateUserParams()
	params.WithBody(&models.User{Email: &email, Name: &name, Password: &password, IsAdmin: &isAdmin})

	r, err := tc.teresa.Users.CreateUser(params, tc.apiKeyAuthFunc)
	if err != nil {
		return nil, err
	}
	return r.Payload, nil
}

// DeleteUser Delete an user
func (tc TeresaClient) DeleteUser(ID int64) error {
	params := users.NewDeleteUserParams()
	params.UserID = ID
	_, err := tc.teresa.Users.DeleteUser(params, tc.apiKeyAuthFunc)
	return err
}

// Me get's the user infos + teams + apps
func (tc TeresaClient) Me() (user *models.User, err error) {
	r, err := tc.teresa.Users.GetCurrentUser(nil, tc.apiKeyAuthFunc)
	if err != nil {
		return nil, err
	}
	return r.Payload, nil
}

// GetAppInfo return teamID and appID
func (tc TeresaClient) GetAppInfo(teamName, appName string) (appInfo AppInfo) {
	me, err := tc.Me()
	if err != nil {
		log.Fatalf("unable to get user information: %s", err)
	}
	if len(me.Teams) > 1 && teamName == "" {
		log.Fatalln("User is in more than one team and provided none")
	}
	for _, t := range me.Teams {
		if teamName == "" || *t.Name == teamName {
			appInfo.TeamID = t.ID
			for _, a := range t.Apps {
				if *a.Name == appName {
					appInfo.AppID = a.ID
					break
				}
			}
			break
		}
	}
	if appInfo.TeamID == 0 || appInfo.AppID == 0 {
		log.Fatalf("Invalid Team [%s] or App [%s]\n", teamName, appName)
	}
	return
}

// GetTeamID returns teamID from team_name
func (tc TeresaClient) GetTeamID(teamName string) (teamID int64) {
	me, err := tc.Me()
	if err != nil {
		log.Fatalf("unable to get user information: %s", err)
	}
	if len(me.Teams) > 1 && teamName == "" {
		log.Fatalln("User is in more than one team and provided none")
	}
	for _, t := range me.Teams {
		if teamName == "" || *t.Name == teamName {
			return t.ID
		}
	}

	log.Fatalf("Invalid Team [%s]\n", teamName)
	return
}

// GetTeams returns a list with my teams
func (tc TeresaClient) GetTeams() (teamsList []*models.Team, err error) {
	params := teams.NewGetTeamsParams()
	r, err := tc.teresa.Teams.GetTeams(params, tc.apiKeyAuthFunc)
	if err != nil {
		return nil, err
	}
	return r.Payload.Items, nil
}

// CreateDeploy creates a new deploy
func (tc TeresaClient) CreateDeploy(teamID, appID int64, description string, tarBall *os.File, writer io.Writer) (io.Writer, error) {
	p := deployments.NewCreateDeploymentParams()
	p.TeamID = teamID
	p.AppID = appID
	p.Description = &description
	p.AppTarball = *tarBall

	r, err := tc.teresa.Deployments.CreateDeployment(p, tc.apiKeyAuthFunc, writer)
	if err != nil {
		return nil, err
	}
	return r.Payload, err
}

// PartialUpdateApp partial updates app... for now, updates only envvars
func (tc TeresaClient) PartialUpdateApp(teamID, appID int64, operations []*models.PatchAppRequest) error {
	p := apps.NewPartialUpdateAppParams()
	p.TeamID = teamID
	p.AppID = appID
	p.Body = operations

	_, err := tc.teresa.Apps.PartialUpdateApp(p, tc.apiKeyAuthFunc)
	return err
}

// AddUserToTeam adds a user (by email) to a team.
// if the user is already part of the team, returns error
func (tc TeresaClient) AddUserToTeam(team, userEmail string) *teams.AddUserToTeamDefault {
	p := teams.NewAddUserToTeamParams()
	p.TeamName = team
	email := strfmt.Email(userEmail)
	p.User.Email = &email
	_, err := tc.teresa.Teams.AddUserToTeam(p, tc.apiKeyAuthFunc)
	if err != nil {
		return err.(*teams.AddUserToTeamDefault)
	}
	return nil
}
