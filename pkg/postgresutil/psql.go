package postgresutil

import (
	"fmt"

	"github.com/curiosinauts/platformctl/pkg/sshutil"
	"github.com/spf13/viper"
)

// PSQLClient psql client
type PSQLClient struct {
	DatabaseName string
	DatabaseHost string
}

// NewPSQLClient returns a new instance of PSQLClient
func NewPSQLClient() *PSQLClient {
	sharedDatabaseName := viper.Get("shared_database_name").(string)
	sharedDatabaseHost := viper.Get("shared_database_host").(string)
	return &PSQLClient{
		DatabaseName: sharedDatabaseName,
		DatabaseHost: sharedDatabaseHost,
	}
}

// CreateUser creates postgres user
func (p *PSQLClient) CreateUser(username, password string) (string, error) {
	script := fmt.Sprintf("psql -d %s -c \"CREATE USER %s ENCRYPTED PASSWORD '%s';\"",
		p.DatabaseName, username, password)
	out, err := sshutil.RemoteSSHExec(p.DatabaseHost, "22", "postgres", script)
	if err != nil {
		return out, err
	}

	return out, nil
}

// CreateUserSchema creates a schema for the given user
func (p *PSQLClient) CreateUserSchema(username string) (string, error) {
	script := fmt.Sprintf("psql -d %s -c \"CREATE SCHEMA AUTHORIZATION %s;\"", p.DatabaseName, username)
	out, err := sshutil.RemoteSSHExec(p.DatabaseHost, "22", "postgres", script)
	if err != nil {
		return out, err
	}

	return out, nil
}

// DropUser drops postgres user
func (p *PSQLClient) DropUser(username string) (string, error) {
	script := fmt.Sprintf("psql -d %s -c \"DROP USER %s;\"", p.DatabaseName, username)
	out, err := sshutil.RemoteSSHExec(p.DatabaseHost, "22", "postgres", script)
	if err != nil {
		return out, err
	}

	return out, nil
}

// DropUserSchema drop schema for the given user
func (p *PSQLClient) DropUserSchema(username string) (string, error) {
	script := fmt.Sprintf("psql -d %s -c \"DROP SCHEMA %s;\"", p.DatabaseName, username)
	out, err := sshutil.RemoteSSHExec(p.DatabaseHost, "22", "postgres", script)
	if err != nil {
		return out, err
	}

	return out, nil
}
