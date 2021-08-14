package cmd

import (
	"fmt"
	"github.com/curiosinauts/platformctl/internal/database"
	"github.com/curiosinauts/platformctl/internal/msg"
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

		eh := ErrorHandler{"before docker-build"}

		if !io.DoesPathExists("./.ssh") {
			err := os.Mkdir("./.ssh", 0755)
			eh.HandleError("creating .ssh folder", err)
		}

		userService := database.NewUserService(db)

		user, dberr := userService.FindByEmail(email)
		eh.HandleError("finding user by email", dberr)

		username := user.Username
		password := user.Password
		randonemail := user.Email

		if debug {
			fmt.Printf("username = %s\n", username)
			fmt.Printf("password = %s\n", password)
			fmt.Printf("email    = %s\n", randonemail)
		}

		io.WriteStringTofile(user.PrivateKey, "./.ssh/id-rsa")
		io.WriteStringTofile(user.PublicKey, "./.ssh/id-rsa.pub")

		//os.RemoveAll("./.ssh")

		// codeserver .config.yml
		io.WriteTemplate(`bind-addr: 0.0.0.0:9991
auth: password
password: {{.}}
cert: false `, password, "./config.yml")

		// .gitconfig
		io.WriteTemplate(`[credential]
    helper = store
[user]
	name = {{.Username}}
	email = {{.Email }}"`, user, "./.gitconfig")

		// repositories.txt
		io.WriteTemplate(`{{range $val := .}}
{{$val}}
{{end}}`, nil, "./repositories.txt")

		msg.Success("before docker-build")
	},
}

func init() {
	beforeCmd.AddCommand(beforeDockerBuildCmd)
}
