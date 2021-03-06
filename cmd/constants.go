package cmd

// variables used to capture the cli flags
var (
	cfgFile            string
	serverFlag         string
	currentFlag        bool
	teamIDFlag         int64
	teamNameFlag       string
	teamEmailFlag      string
	teamURLFlag        string
	userIDFlag         int64
	userNameFlag       string
	userEmailFlag      string
	userPasswordFlag   string
	appNameFlag        string
	appScaleFlag       int
	descriptionFlag    string
	autocompleteTarget string
	isAdminFlag        bool
)

const (
	version               = "0.1.2"
	archiveTempFolder     = "/tmp"
	deploymentSuccessMark = "----------deployment-success----------"
	deploymentErrorMark   = "----------deployment-error----------"
)
