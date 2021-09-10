package postgresutil

import (
	"fmt"

	"github.com/curiosinauts/platformctl/pkg/sshutil"
)

// CreateUser creates postgres user
func CreateUser(username, password string) (string, error) {
	script := fmt.Sprintf("psql -d curiositydb -c \"CREATE USER %s ENCRYPTED PASSWORD '%s';\"", username, password)
	out, err := sshutil.RemoteSSHExec("192.168.0.116", "22", "postgres", script)
	if err != nil {
		return out, err
	}

	return out, nil
}

// CreateUserSchema creates a schema for the given user
func CreateUserSchema(username string) (string, error) {
	script := fmt.Sprintf("psql -d curiositydb -c \"CREATE SCHEMA AUTHORIZATION %s;\"", username)
	out, err := sshutil.RemoteSSHExec("192.168.0.116", "22", "postgres", script)
	if err != nil {
		return out, err
	}

	return out, nil
}

// DropUser drops postgres user
func DropUser(username string) (string, error) {
	script := fmt.Sprintf("psql -d curiositydb -c \"DROP USER %s;\"", username)
	out, err := sshutil.RemoteSSHExec("192.168.0.116", "22", "postgres", script)
	if err != nil {
		return out, err
	}

	return out, nil
}

// DropUserSchema drop schema for the given user
func DropUserSchema(username string) (string, error) {
	script := fmt.Sprintf("psql -d curiositydb -c \"DROP SCHEMA %s;\"", username)
	out, err := sshutil.RemoteSSHExec("192.168.0.116", "22", "postgres", script)
	if err != nil {
		return out, err
	}

	return out, nil
}
