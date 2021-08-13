package cmd

import (
	"fmt"
	"log"

	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/curiosinauts/platformctl/pkg/enc"
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
		hashedEmail := enc.Hashed(email)

		fmt.Printf("hashed email %s\n", hashedEmail)

		haikunator := haikunator.New()
		randomUsername := haikunator.Haikunate()

		fmt.Printf("random_username = %s\n", randomUsername)

		// Generate a password that is 64 characters long with 10 digits, 10 symbols,
		// allowing upper and lower case letters, disallowing repeat characters.
		res, err := password.Generate(32, 10, 0, false, false)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("generated password = %s\n", res)

		randomEmail := fmt.Sprintf("%s@example.com", randomUsername)

		fmt.Printf("random_email = %s\n", randomEmail)

		privateKey, publicKey := enc.GenerateRSAKeys()

		fmt.Printf("private key = %s", privateKey)
		fmt.Printf("public key = %s", publicKey)
	},
}

func init() {
	addCmd.AddCommand(addUserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
