package cmd

import (
	"fmt"
	"log"

	"github.com/curiosinauts/platformctl/pkg/giteautil"
	"github.com/curiosinauts/platformctl/pkg/jenkinsutil"
	"github.com/google/uuid"

	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
)

var addUserCmdUseExistingKeys bool
var addUserCmdRepos []string
var addUserCmdUsername string
var addUserCmdUseEmail bool

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

		randomPassword, err := password.Generate(32, 10, 0, false, false)
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
			fmt.Printf("generated password : %s\n", randomPassword)
			fmt.Printf("random_email       : %s\n", email)
			fmt.Printf("private key        : \n%s", privateKey)
			fmt.Printf("public key         : \n%s", publicKey)
		}

		user := database.User{
			Username:    username,
			Password:    randomPassword,
			Email:       email,
			HashedEmail: hashedEmail,
			PrivateKey:  string(privateKey),
			PublicKey:   string(publicKey),
			IsActive:    true,
		}

		eh := ErrorHandler{"adding user"}

		result, dberr := userService.Add(user)
		eh.HandleError("user insert", dberr)
		repoURI := fmt.Sprintf("ssh://gitea@git-ssh.curiosityworks.org:2222/%s/project.git", username)
		userID, err := result.LastInsertId()
		eh.HandleError("user id", err)
		user.ID = userID

		_, dberr = userService.AddUserRepo(database.UserRepo{
			URI:    repoURI,
			UserID: userID,
		})

		if len(addUserCmdRepos) > 0 {
			AddUserRepos(user.ID, addUserCmdRepos)
		}

		ide, dberr := userService.FindIDEByName("vscode")
		eh.HandleError("finding ide", dberr)

		result, dberr = userService.AddUserIDE(database.UserIDE{
			UserID: userID,
			IDEID:  ide.ID,
		})
		eh.HandleError("user_ide insert", dberr)

		userIDEID, err := result.LastInsertId()
		eh.HandleError("user_ide new id", err)

		_, dberr = userService.AddIDERepo(database.IDERepo{
			UserIDEID: userIDEID,
			URI:       repoURI,
		})

		if len(addUserCmdRepos) > 0 {
			AddIDERepos(userIDEID, addUserCmdRepos)
		}

		eh.HandleError("ide_repo insert", dberr)

		runtimeInstall, dberr := userService.FindRuntimeInstallByName("tmux")
		eh.HandleError("finding runtime install", dberr)

		_, dberr = userService.AddIDERuntimeInstall(database.IDERuntimeInstall{
			UserIDEID:        userIDEID,
			RuntimeInstallID: runtimeInstall.ID,
		})
		eh.HandleError("ide_runtime_install insert", dberr)

		gitClient, err := giteautil.NewGitClient()
		eh.HandleError("instantiating git client", err)

		err = gitClient.AddUser(user)
		eh.HandleError("adding user to gitea", err)

		err = gitClient.CreateUserRepo(user.Username)
		eh.HandleError("create user repo", err)

		publicKeyID, err := gitClient.CreateUserPublicKey(user)
		eh.HandleError("create user public key", err)

		user.PublicKeyID = publicKeyID
		result, dberr = userService.UpdateProfile(user)
		eh.HandleError("updating user profile", dberr)

		jenkins, err := jenkinsutil.NewJenkins()
		eh.HandleError("accessing Jenkins job", err)

		option := map[string]string{
			"USERNAME": user.Username,
			"VERSION":  uuid.NewString(),
		}
		_, err = jenkins.BuildJob("codeserver", option)
		eh.HandleError("calling Jenkins job to build codeserver instance", err)

		msg.Success("adding user")
	},
}

func AddUserRepos(userID int64, repos []string) *database.DBError {
	for _, repo := range repos {
		_, dberr := userService.AddUserRepo(database.UserRepo{
			URI:    repo,
			UserID: userID,
		})
		if dberr != nil {
			return dberr
		}
	}
	return nil
}

func AddIDERepos(userIDEID int64, repos []string) *database.DBError {
	for _, repo := range repos {
		_, dberr := userService.AddIDERepo(database.IDERepo{
			UserIDEID: userIDEID,
			URI:       repo,
		})
		if dberr != nil {
			return dberr
		}
	}
	return nil
}

func init() {
	addCmd.AddCommand(addUserCmd)
	addUserCmd.Flags().BoolVar(&addUserCmdUseExistingKeys, "pki", false, "use existing PKI")
	addUserCmd.Flags().StringArrayVarP(&addUserCmdRepos, "repo", "r", []string{}, "-r https://example-repo.com/foo")
	addUserCmd.Flags().StringVar(&addUserCmdUsername, "username", "", "specify username instead of auto generated")
	addUserCmd.Flags().BoolVarP(&addUserCmdUseEmail, "real-email", "e", false, "use real email instead of auto generated")
}
