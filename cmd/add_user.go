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

// userCmd represents the user command
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
		dberr := userService.Add(user)
		if dberr != nil {
			fmt.Printf("adding user failed: %v", dberr)
		}
	},
}

func init() {
	addCmd.AddCommand(addUserCmd)
}
