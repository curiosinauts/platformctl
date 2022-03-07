package postgresutil

import (
	"fmt"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/executil"
	"github.com/curiosinauts/platformctl/pkg/sshutil"
	"github.com/spf13/viper"
)

// PSQLClient psql client
type PSQLClient struct {
	DatabaseName string
	DatabaseHost string
}

// NewPSQLClientForSharedDB returns a new instance of PSQLClient
func NewPSQLClientForSharedDB() *PSQLClient {
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
func (p *PSQLClient) CreateUser(username, password string, debug bool) (string, error) {
	script := fmt.Sprintf("psql -d %s -c \"CREATE USER %s ENCRYPTED PASSWORD '%s';\"",
		p.DatabaseName, username, password)
	return p.ExecutePSQLScriptOverSSH(script, debug)
}

// CreateUserSchema creates a schema for the given user
func (p *PSQLClient) CreateUserSchema(username string, debug bool) (string, error) {
	script := fmt.Sprintf("psql -d %s -c \"CREATE SCHEMA AUTHORIZATION %s;\"", p.DatabaseName, username)
	return p.ExecutePSQLScriptOverSSH(script, debug)
}

// DropUser drops postgres user
func (p *PSQLClient) DropUser(username string, debug bool) (string, error) {
	script := fmt.Sprintf("psql -d %s -c \"DROP USER %s;\"", p.DatabaseName, username)
	return p.ExecutePSQLScriptOverSSH(script, debug)
}

// DropUserSchema drop schema for the given user
func (p *PSQLClient) DropUserSchema(username string, debug bool) (string, error) {
	script := fmt.Sprintf("psql -d %s -c \"DROP SCHEMA %s;\"", p.DatabaseName, username)
	return p.ExecutePSQLScriptOverSSH(script, debug)
}

// ListDBUsers list db users
func (p *PSQLClient) ListDBUsers(debug bool) (string, error) {
	script := fmt.Sprintf(`psql -c "\du"`)
	return p.ExecutePSQLScriptOverSSH(script, debug)
}

// ExecutePSQLScriptOverSSH executes psql script over ssh
func (p *PSQLClient) ExecutePSQLScriptOverSSH(script string, debug bool) (string, error) {
	if debug {
		fmt.Println("  + " + script)
	}
	out, err := sshutil.RemoteSSHExec(p.DatabaseHost, "22", "postgres", script)
	out = strings.TrimSpace(out)
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

// BackUpSchemaWithData backs up all tables from a schema using pg_dump
func (p *PSQLClient) BackUpSchemaWithData(password string, username string, host string, dbname string, schemaName string, debug bool) (string, string, error) {
	sqlFilePath := fmt.Sprintf("/tmp/%s.sql", schemaName)
	script := fmt.Sprintf("export PGPASSWORD=\"%s\"; pg_dump -U %s -h %s -d %s -n %s > "+sqlFilePath,
		password, username, p.DatabaseHost, p.DatabaseName, schemaName)
	out, err := p.ExecuteLocalShell(script, debug)
	return out, sqlFilePath, err
}

// RestoreSchemaWithData backs up all tables from a schema using pg_dump
func (p *PSQLClient) RestoreSchemaWithData(password string, username string, host string, dbname string, schemaName string, debug bool) (string, string, error) {
	sqlFilePath := fmt.Sprintf("/tmp/%s.sql", schemaName)
	script := fmt.Sprintf("export PGPASSWORD=\"%s\"; psql -h %s -U %s \"dbname=%s options=--search_path=%s\" < "+sqlFilePath,
		password, p.DatabaseHost, username, p.DatabaseName, schemaName)
	out, err := p.ExecuteLocalShell(script, debug)
	return out, sqlFilePath, err
}

// ExecuteLocalShell executes local shell script
func (p *PSQLClient) ExecuteLocalShell(script string, debug bool) (string, error) {
	if debug {
		fmt.Println("  + " + script)
	}
	out, err := executil.ExecuteShell(script, debug)
	out = strings.TrimSpace(out)
	if err != nil {
		return out, err
	}

	return out, nil
}
