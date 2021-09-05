package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/giteautil"
	"github.com/curiosinauts/platformctl/pkg/jenkinsutil"
	"github.com/google/uuid"
	pwd "github.com/sethvargo/go-password/password"

	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
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
			Username:    username,
			Password:    password,
			Email:       email,
			HashedEmail: hashedEmail,
			PrivateKey:  string(privateKey),
			PublicKey:   string(publicKey),
			IsActive:    true,
		}

		eh := ErrorHandler{"adding user"}

		dbs := database.NewUserService(db)

		dberr := dbs.Save(&user)
		eh.HandleError("user insert", dberr)
		repoURI := fmt.Sprintf("ssh://gitea@git-ssh.curiosityworks.org:2222/%s/project.git", username)
		eh.HandleError("user id", err)

		dberr = dbs.Save(database.NewUserRepo(repoURI, user.ID))
		eh.HandleError("saving new user repo", dberr)

		if len(addUserCmdRepos) > 0 {
			AddUserRepos(user.ID, addUserCmdRepos)
		}

		ide := new(database.IDE)
		dberr = dbs.FindBy(ide, "name=$1", "vscode")
		eh.HandleError("finding ide", dberr)

		userIDE := database.UserIDE{
			UserID: user.ID,
			IDEID:  ide.ID,
		}
		dberr = dbs.Save(&userIDE)
		eh.HandleError("user_ide insert", dberr)

		// userIDEID, err := result.LastInsertId()
		eh.HandleError("user_ide new id", err)

		dberr = dbs.Save(&database.IDERepo{
			UserIDEID: userIDE.ID,
			URI:       repoURI,
		})

		if len(addUserCmdRepos) > 0 {
			AddIDERepos(userIDE.ID, addUserCmdRepos)
		}

		eh.HandleError("ide_repo insert", dberr)

		runtimeInstall := database.RuntimeInstall{}
		eh.HandleError("finding runtime install", dbs.FindBy(&runtimeInstall, "name=$1", "tmux"))

		dberr = dbs.Save(&database.IDERuntimeInstall{
			UserIDEID:        userIDE.ID,
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
		_, dberr = dbs.UpdateProfile(user)
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

// AddUserRepos adds user repos
func AddUserRepos(userID int64, repos []string) *database.DBError {
	for _, repo := range repos {
		dberr := dbs.Save(&database.UserRepo{
			URI:    repo,
			UserID: userID,
		})
		if dberr != nil {
			return dberr
		}
	}
	return nil
}

// AddIDERepos adds ide repos
func AddIDERepos(userIDEID int64, repos []string) *database.DBError {
	ideRepos := []database.IDERepo{}
	dbs.ListBy("ide_repo", &ideRepos, "user_ide_id=$1", userIDEID)

nextRepo:
	for _, repo := range repos {
		for _, ideRepo := range ideRepos {
			if strings.TrimSpace(ideRepo.URI) == strings.TrimSpace(repo) {
				msg.Info("repo already exists. skipping. " + repo)
				continue nextRepo
			}
		}
		dberr := dbs.Save(&database.IDERepo{
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
	addUserCmd.Flags().BoolVarP(&addUserCmdUseExistingKeys, "pki", "p", false, "use existing PKI or not")
	addUserCmd.Flags().BoolVarP(&addUserCmdUseEmail, "email", "e", false, "use real email or not")
	addUserCmd.Flags().StringVarP(&addUserCmdUsername, "username", "u", "", "specify username")
	addUserCmd.Flags().StringArrayVarP(&addUserCmdRepos, "repo", "r", []string{}, "specify personal git repo")
}
