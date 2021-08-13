package cmd

import (
	"fmt"
	"github.com/curiosinauts/platformctl/internal/database"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"log"

	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
)

// addUserCmd represents the user command
var addUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Adds user",
	Long: `Adding user means adding user to the user database, user's own dev database, gitea, pgadmin, 
	runtime installs, user repos.`,
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]

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
		result, dberr := userService.Add(user)
		if dberr != nil {
			fmt.Printf("adding user failed: %v", dberr)
		}

		repoURI := fmt.Sprintf("ssh://gitea@git-ssh.curiosityworks.org:2222/%s/project.git", randomUsername)
		userID, err := result.LastInsertId()
		if err != nil {
			fmt.Printf("adding user failed: %v", err)
		}
		userRepo := database.UserRepo{
			URI:    repoURI,
			UserID: userID,
		}
		userService.AddRepo(userRepo)

		ide, dberr := userService.FindIDEByName("vscode")
		if dberr != nil {
			fmt.Printf("adding user failed: %v", dberr)
		}

		userIDE := database.UserIDE{
			UserID: userID,
			IDEID:  ide.ID,
		}
		result, dberr = userService.AddUserIDE(userIDE)
		if dberr != nil {
			fmt.Printf("adding user failed: %v", dberr)
		}

		userIDEID, err := result.LastInsertId()
		if err != nil {
			fmt.Printf("adding user failed: %v", err)
		}

		runtimeInstall, dberr := userService.FindRuntimeInstallName("tmux")
		if dberr != nil {
			fmt.Printf("adding user failed: %v", dberr)
		}

		ideRuntimeInstall := database.IDERuntimeInstall{
			UserIDEID:        userIDEID,
			RuntimeInstallID: runtimeInstall.ID,
		}
		userService.AddIDERuntimeInstall(ideRuntimeInstall)
	},
}

func init() {
	addCmd.AddCommand(addUserCmd)
}
