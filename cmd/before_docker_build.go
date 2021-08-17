package cmd

import (
	"fmt"
	"os"

	"github.com/curiosinauts/platformctl/internal/database"
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/io"

	"github.com/spf13/cobra"
)

// beforeDockerBuildCmd represents the dockerBuild command
var beforeDockerBuildCmd = &cobra.Command{
	Use:   "docker-build",
	Short: "Generates files needed by code server docker build",
	Long:  `Generate files for docker build. Generated files are used to customize code server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		username := args[0]
		dockertag := args[1]

		eh := ErrorHandler{"before docker-build"}

		if !io.DoesPathExists("./.ssh") {
			err := os.Mkdir("./.ssh", 0755)
			eh.HandleError("creating .ssh folder", err)
		}

		userService := database.NewUserService(db)

		user, dberr := userService.FindUserByUsername(username)
		eh.HandleError("finding user by email", dberr)

		io.WriteStringTofile(user.PrivateKey, "./.ssh/id_rsa")

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
	email = {{.Email }}`, user, "./.gitconfig")

		repositories, dberr := userService.FindUserIDEReroURIsByUserAndIDE(username, "vscode")
		// repositories.txt
		io.WriteTemplate(`{{range $val := .}}
{{$val}}
{{end}}`, repositories, "./repositories.txt")

		runtimeInstalls, dberr := userService.FindUserIDERuntimeInstallsByUserAndIDE(username, "vscode")
		io.WriteTemplate(`#!/bin/bash -e
    
set -x
{{range $val := .}}
{{$val}}
{{end}}`, runtimeInstalls, "./runtime_install.sh")

		err := os.Chmod("./runtime_install.sh", 0755)
		eh.HandleError("writing runtime install", err)

		io.WriteTemplate(`#!/bin/sh

export TERM=xterm

/home/coder/go/bin/gotty --ws-origin "vscode-{{.Username}}.curiosityworks.org" -p 2222 -c "{{.Username}}:{{.Password}}" -w /usr/bin/zsh >>/dev/null 2>&1 
`, user, "./gotty.sh")
		err = os.Chmod("./gotty.sh", 0755)

		user.DockerTag = dockertag
		io.WriteTemplate(deployServiceIngressTemplate, user, "./vscode-"+user.Username+".yml")

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
      - name: vscode-{{.Username}}
        image: docker-registry.int.curiosityworks.org/7onetella/vscode-{{.Username}}:{{.DockerTag}}
        ports:
        - containerPort: 9991

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
