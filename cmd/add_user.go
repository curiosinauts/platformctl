package cmd

import (
	"fmt"
	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/curiosinauts/platformctl/internal/database"
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"log"
)

// addUserCmd represents the user command
var addUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Adds user",
	Long: `Adding user means adding user to the user database, user's own dev database, gitea, pgadmin, 
	runtime installs, user repos.`,
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]

		fmt.Println()

		hashedEmail := crypto.Hashed(email)

		randomUsername := haikunator.New().Haikunate()

		randomPassword, err := password.Generate(32, 10, 0, false, false)
		if err != nil {
			log.Fatal(err)
		}

		randomEmail := fmt.Sprintf("%s@example.com", randomUsername)

		privateKey, publicKey := crypto.GenerateRSASSHKeys()

		if debug {
			fmt.Printf("hashed email       : %s\n", hashedEmail)
			fmt.Printf("random_username    : %s\n", randomUsername)
			fmt.Printf("generated password : %s\n", randomPassword)
			fmt.Printf("random_email       : %s\n", randomEmail)
			fmt.Printf("private key        : \n%s", privateKey)
			fmt.Printf("public key         : \n%s", publicKey)
		}

		userService := database.NewUserService(db)

		user := database.User{
			Username:    randomUsername,
			Password:    randomPassword,
			Email:       randomEmail,
			HashedEmail: hashedEmail,
			PrivateKey:  string(privateKey),
			PublicKey:   string(publicKey),
			IsActive:    true,
		}

		eh := ErrorHandler{"adding user"}

		result, dberr := userService.Add(user)
		if dberr != nil {
			fmt.Println(dberr.Err)
		}
		eh.HandleError("user insert", dberr)
		repoURI := fmt.Sprintf("ssh://gitea@git-ssh.curiosityworks.org:2222/%s/project.git", randomUsername)
		userID, err := result.LastInsertId()
		eh.HandleError("user id", err)

		_, dberr = userService.AddUserRepo(database.UserRepo{
			URI:    repoURI,
			UserID: userID,
		})

		ide, dberr := userService.FindIDEByName("vscode")
		eh.HandleError("finding ide", dberr)

		result, dberr = userService.AddUserIDE(database.UserIDE{
			UserID: userID,
			IDEID:  ide.ID,
		})
		eh.HandleError("user_ide insert", dberr)

		userIDEID, err := result.LastInsertId()
		eh.HandleError("user_ide new id", err)

		runtimeInstall, dberr := userService.FindRuntimeInstallName("tmux")
		eh.HandleError("finding runtime install", dberr)

		_, dberr = userService.AddIDERuntimeInstall(database.IDERuntimeInstall{
			UserIDEID:        userIDEID,
			RuntimeInstallID: runtimeInstall.ID,
		})
		eh.HandleError("ide_runtime_install insert", dberr)

		msg.Success("adding user")
	},
}

func init() {
	addCmd.AddCommand(addUserCmd)
}
