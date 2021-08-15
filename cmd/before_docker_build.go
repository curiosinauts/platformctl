package cmd

import (
	"fmt"
	"github.com/curiosinauts/platformctl/internal/database"
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/io"
	"os"

	"github.com/spf13/cobra"
)

// beforeDockerBuildCmd represents the dockerBuild command
var beforeDockerBuildCmd = &cobra.Command{
	Use:   "docker-build",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		cobra.MinimumNArgs(1)(cmd, args)

		email := args[0]
		hashedEmail := crypto.Hashed(email)

		eh := ErrorHandler{"before docker-build"}

		if !io.DoesPathExists("./.ssh") {
			err := os.Mkdir("./.ssh", 0755)
			eh.HandleError("creating .ssh folder", err)
		}

		userService := database.NewUserService(db)

		user, dberr := userService.FindUserByHashedEmail(hashedEmail)
		eh.HandleError("finding user by email", dberr)

		io.WriteStringTofile(user.PrivateKey, "./.ssh/id-rsa")
		io.WriteStringTofile(user.PublicKey, "./.ssh/id-rsa.pub")

		//os.RemoveAll("./.ssh")

		// codeserver .config.yml
		io.WriteTemplate(`bind-addr: 0.0.0.0:9991
auth: password
password: {{.}}
cert: false `, user.Password, "./config.yml")

		// .gitconfig
		io.WriteTemplate(`[credential]
    helper = store
[user]
	name = {{.Username}}
	email = {{.Email }}"`, user, "./.gitconfig")

		repositories, dberr := userService.FindUserIDEReroURIsByUserAndIDE(hashedEmail, "vscode")
		// repositories.txt
		io.WriteTemplate(`{{range $val := .}}
{{$val}}
{{end}}`, repositories, "./repositories.txt")

		msg.Success("before docker-build")
	},
}

func init() {
	beforeCmd.AddCommand(beforeDockerBuildCmd)
}
