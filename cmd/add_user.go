package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/giteautil"
	"github.com/curiosinauts/platformctl/pkg/postgresutil"
	pwd "github.com/sethvargo/go-password/password"

	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

var addUserCmdUseExistingKeys bool
var addUserCmdUsername string
var addUserCmdUseEmail bool
var addUserCmdRuntimeInstalls string

// addUserCmd represents the user command
var addUserCmd = &cobra.Command{
	Use:     "user",
	Short:   "Adds user",
	Long:    `Adds user to the platform`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		originalEmail := args[0]

		hashedEmail := crypto.Hashed(originalEmail)

		username := haikunator.New().Haikunate()
		if len(addUserCmdUsername) > 0 {
			username = addUserCmdUsername
		}

		password, err := pwd.Generate(32, 10, 0, false, false)
		if err != nil {
			log.Fatal(err)
		}

		email := fmt.Sprintf("%s@curiosityworks.org", username)
		if addUserCmdUseEmail {
			email = originalEmail
		}

		privateKey, publicKey := crypto.GenerateRSASSHKeys()
		if addUserCmdUseExistingKeys {
			privateKey, publicKey = crypto.ReadExistingRSASSHKeys()
		}

		if debug {
			fmt.Printf("hashed email       : %s\n", hashedEmail)
			fmt.Printf("random_username    : %s\n", username)
			fmt.Printf("generated password : %s\n", password)
			fmt.Printf("random_email       : %s\n", email)
			fmt.Printf("private key        : \n%s", privateKey)
			fmt.Printf("public key         : \n%s", publicKey)
		}

		user := database.User{
			Username:        username,
			Password:        password,
			Email:           email,
			HashedEmail:     hashedEmail,
			PrivateKey:      string(privateKey),
			PublicKey:       string(publicKey),
			IsActive:        true,
			GitRepoURI:      fmt.Sprintf("ssh://gitea@git-ssh.curiosityworks.org:2222/%s/project.git", username),
			IDE:             "vscode",
			RuntimeInstalls: addUserCmdRuntimeInstalls,
		}

		eh := ErrorHandler{"adding user"}

		dbs := database.NewUserService(db)
		dberr := dbs.Save(&user)
		eh.HandleError("user insert", dberr)
		eh.HandleError("user id", err)

		gitClient, err := giteautil.NewGitClient()
		eh.HandleError("instantiating git client", err)

		err = gitClient.AddUser(user)
		eh.HandleError("adding user to gitea", err)

		err = gitClient.CreateUserRepo(user.Username)
		eh.HandleError("create user repo", err)

		publicKeyID, err := gitClient.CreateUserPublicKey(user)
		eh.HandleError("create user public key", err)

		user.PublicKeyID = publicKeyID
		_, dberr = dbs.UpdateProfile(user)
		eh.HandleError("updating user profile", dberr)

		postgresUsername := strings.Replace(user.Username, "-", "", -1)
		psql := postgresutil.NewPSQLClientForSharedDB()
		_, err = psql.CreateUser(postgresUsername, password, debug)
		eh.HandleError("creating database user", err)

		_, err = psql.CreateUserSchema(postgresUsername, debug)
		eh.HandleError("creating database user schema", err)

		// jenkins, err := jenkinsutil.NewJenkins()
		// eh.HandleError("accessing Jenkins job", err)

		// params := map[string]string{
		// 	"USERNAME": user.Username,
		// 	"VERSION":  uuid.NewString(),
		// }
		// _, err = jenkins.BuildJob("codeserver", params)
		// eh.HandleError("calling Jenkins job to build codeserver instance", err)

		msg.Success("adding user")
	},
}

func init() {
	addCmd.AddCommand(addUserCmd)
	addUserCmd.Flags().BoolVarP(&addUserCmdUseExistingKeys, "pki", "p", false, "use existing PKI or not")
	addUserCmd.Flags().BoolVarP(&addUserCmdUseEmail, "email", "e", false, "use real email or not")
	addUserCmd.Flags().StringVarP(&addUserCmdUsername, "username", "u", "", "specify username")
	addUserCmd.Flags().StringVarP(&addUserCmdRuntimeInstalls, "runtime-installs", "i", "tmux", "runtime installs")
}
