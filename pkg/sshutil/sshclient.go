package sshutil

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// CREDIT: https://github.com/Scalingo/go-ssh-examples/blob/master/client.go
// CREDIT: https://gist.github.com/Mebus/c3a437e339481de03a98569090c53b08

// PublicKeyFile derieves public off of private key
func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

// RemoteSSHExec remote ssh and execute script
func RemoteSSHExec(server, port, user, script string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	publicKey, err := PublicKeyFile(fmt.Sprintf("%s/.ssh/id_rsa", home))
	if err != nil {
		return "", err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			publicKey,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", server+":"+port, config)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	out, err := session.CombinedOutput(script)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
