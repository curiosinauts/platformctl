package postgresutil

import (
	"fmt"

	"github.com/curiosinauts/platformctl/pkg/executil"
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

// NewPSQLClientByHostAndDBName returns a new instance of PSQLClient
func NewPSQLClientByHostAndDBName(host string, dbname string) *PSQLClient {
	return &PSQLClient{
		DatabaseName: dbname,
		DatabaseHost: host,
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

// BackUpDBWithData backs up all tables from a schema using pg_dump
func (p *PSQLClient) BackUpDBWithData(password string, username string, host string, dbname string, debug bool) (string, error) {
	script := fmt.Sprintf("export PGPASSWORD=\"%s\"; pg_dump -U %s -h %s %s > /tmp/%s.sql",
		username, password, p.DatabaseHost, p.DatabaseName, p.DatabaseName)
	out, err := executil.Execute(script, debug)
	if err != nil {
		return out, err
	}

	return out, nil
}

// BackUpSchemaOnlyWithData backs up all tables from a schema using pg_dump
func (p *PSQLClient) BackUpSchemaOnlyWithData(password string, username string, host string, dbname string, schemaName string, debug bool) (string, error) {
	script := fmt.Sprintf("export PGPASSWORD=\"%s\"; pg_dump -U %s -h %s -d %s -n %s > /tmp/%s.sql",
		password, username, p.DatabaseHost, p.DatabaseName, schemaName, schemaName)
	if debug {
		fmt.Println("+ " + script)
	}
	out, err := executil.ExecuteShell(script, debug)
	if err != nil {
		return out, err
	}

	return out, nil
}

// RestoreSchemaWithData backs up all tables from a schema using pg_dump
func (p *PSQLClient) RestoreSchemaWithData(password string, username string, host string, dbname string, schemaName string, debug bool) (string, error) {
	script := fmt.Sprintf("export PGPASSWORD=\"%s\"; psql -h %s -U %s \"dbname=%s options=--search_path=%s\" < /tmp/%s.sql",
		password, p.DatabaseHost, username, p.DatabaseName, schemaName, schemaName)
	if debug {
		fmt.Println("+ " + script)
	}
	out, err := executil.ExecuteShell(script, debug)
	if err != nil {
		return out, err
	}

	return out, nil
}
