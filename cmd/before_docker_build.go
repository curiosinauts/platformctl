package cmd

import (
	"os"
	"strings"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/curiosinauts/platformctl/pkg/io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// beforeDockerBuildCmd represents the dockerBuild command
var beforeDockerBuildCmd = &cobra.Command{
	Use:     "docker-build",
	Short:   "Generates files needed by code server docker build",
	Long:    `Generate files for docker build. Generated files are used to customize code server`,
	PreRunE: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		username := args[0]
		dockertag := args[1]

		eh := ErrorHandler{"before docker-build"}

		user := database.User{}
		eh.HandleError("finding user by username", dbs.FindBy(&user, "username=$1", username))

		// code server required authentication file
		// createCoderConfigYamlFile(user.Password)

		// user's private terminal app
		// createGottySSH(user, eh)

		// create k3s deployment service ingress yaml files
		user.DockerTag = dockertag
		createDeploymentServiceIngressYamlFile(user)

		// create k3s secrets file for given user
		createUserSecretsFile(user)

		// disable the following customizations until we are ready
		if false {
			// id_rsa key for git authentication
			createSSHFolder(eh)
			createIDRSAKey(user.PrivateKey)

			// .gitconfig for git operation
			createGITConfig(user)

			// git clone of user projects
			repositories := []string{user.GitRepoURI}
			createRepositoriesTXT(repositories)

			// self service runtime installs config file
			runtimeInstalls := []database.RuntimeInstall{}
			dbs.FindAllRuntimeInstallsForUser(&runtimeInstalls, username)
			createRuntimeInstallSSHFile(runtimeInstalls, eh)

			// environment variables for things like postgresl username and password
			createDotExportsFile(user, eh)
		}

		msg.Success("before docker-build")
	},
}

func init() {
	beforeCmd.AddCommand(beforeDockerBuildCmd)
}

var deployServiceIngressTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: vscode-{{.Username}}
  labels:
    app: vscode-{{.Username}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vscode-{{.Username}}
  template:
    metadata:
      labels:
        app: vscode-{{.Username}}
    spec:
      containers:
        - env:
          - name: VSCODE_USERNAME
            valueFrom:
              secretKeyRef:
                key: vscode-username
                name: vscode-secrets-{{.Username}}	  
          - name: VSCODE_PASSWORD
            valueFrom:
              secretKeyRef:
                key: vscode-password
                name: vscode-secrets-{{.Username}}
          name: vscode-{{.Username}}
          image: curiosinauts/vscode-ext:0.1.0
          ports:
            - containerPort: 9991
      dnsPolicy: Default 

---
apiVersion: v1
kind: Service
metadata:
  name: vscode-{{.Username}}
spec:
  selector:
    app: vscode-{{.Username}}
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9991

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: vscode-{{.Username}}
  annotations:
    kubernetes.io/ingress.class: "traefik"
spec:
  rules:
  - host: vscode-{{.Username}}.curiosityworks.org
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: vscode-{{.Username}}
            port: 
              number: 80
`

var secretsTemplate = `apiVersion: v1
kind: Secret
metadata:
  name: vscode-secrets-{{.Username}}
type: Opaque
stringData:
  vscode-username: "{{.Username}}"
  vscode-password: "{{.Password}}"
`

func createSSHFolder(eh ErrorHandler) {
	if !io.DoesPathExists("./.ssh") {
		err := os.Mkdir("./.ssh", 0755)
		eh.HandleError("creating .ssh folder", err)
	}
}

func createIDRSAKey(privateKey string) {
	io.WriteStringTofile(privateKey, "./.ssh/id_rsa")
}

func createCoderConfigYamlFile(password string) {
	io.WriteTemplate(`bind-addr: 0.0.0.0:9991
auth: password
password: {{.}}
cert: false `, password, "./config.yml")
}

func createGITConfig(user database.User) {
	io.WriteTemplate(`[credential]
    helper = store
[user]
	name = {{.Username}}
	email = {{.Email }}`, user, "./.gitconfig")
}

func createRepositoriesTXT(repositories []string) {
	io.WriteTemplate(`{{range $val := .}}
{{$val}}
{{end}}`, repositories, "./repositories.txt")
}

func createRuntimeInstallSSHFile(runtimeInstalls []database.RuntimeInstall, eh ErrorHandler) {
	io.WriteTemplate(`#!/bin/zsh -e
    
set -x
{{range $v := .}}
{{$v.ScriptBody}}
{{end}}`, runtimeInstalls, "./runtime_install.sh")
	err := os.Chmod("./runtime_install.sh", 0755)
	eh.HandleError("writing runtime install", err)

}

func createGottySSH(user database.User, eh ErrorHandler) {
	io.WriteTemplate(`#!/bin/sh

	export TERM=xterm
	
	/home/coder/go/bin/gotty --ws-origin "vscode-{{.Username}}.curiosityworks.org" -p 2222 -c "{{.Username}}:{{.Password}}" -w /usr/bin/zsh >>/dev/null 2>&1 
	`, user, "./gotty.sh")
	err := os.Chmod("./gotty.sh", 0755)
	eh.HandleError("writing .gotty.sh", err)
}

func createDotExportsFile(user database.User, eh ErrorHandler) {
	user.PostgresUsername = strings.Replace(user.Username, "-", "", -1)
	user.PGHost = viper.Get("shared_database_host").(string)
	user.PGDBName = viper.Get("shared_database_name").(string)
	io.WriteTemplate(`
export PGUSER={{.PostgresUsername}}
export PGPASSWORD={{.Password}}
export PGHOST={{.PGHost}}
export PGDATABASE={{.PGDBName}}
`, user, "./.exports")
	err := os.Chmod("./.exports", 0755)
	eh.HandleError("writing .exports", err)
}

func createDeploymentServiceIngressYamlFile(user database.User) {
	io.WriteTemplate(deployServiceIngressTemplate, user, "./vscode-"+user.Username+".yml")
}

func createUserSecretsFile(user database.User) {
	io.WriteTemplate(secretsTemplate, user, "./vscode-"+user.Username+"-secrets.yml")
}
