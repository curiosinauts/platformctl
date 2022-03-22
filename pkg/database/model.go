package database

import "fmt"

// User user
type User struct {
	ID               int64  `db:"id"`
	GoogleID         string `db:"google_id"`
	Username         string `db:"username"`
	Password         string `db:"password"`
	Email            string `db:"email"`
	HashedEmail      string `db:"hashed_email"`
	IsActive         bool   `db:"is_active"`
	PrivateKey       string `db:"private_key"`
	PublicKey        string `db:"public_key"`
	PublicKeyID      int64  `db:"public_key_id"`
	GitRepoURI       string `db:"git_repo_uri"`
	IDE              string `db:"ide"`
	RuntimeInstalls  string `db:"runtime_installs"`
	DockerTag        string
	Installed        []RuntimeInstall
	PostgresUsername string
	PGHost           string
	PGDBName         string
}

// String string representation of user
func (u User) String() string {
	return fmt.Sprintf("%d: %s\n", u.ID, u.Username)
}

// Meta provides mapping config specific to user
func (u *User) Meta() MappingConfig {
	return MappingConfig{TableName: "users"}
}

// PrimaryKey returns primary key
func (u *User) PrimaryKey() int64 {
	return u.ID
}

// SetPrimaryKey updates the primary key value after insert
func (u *User) SetPrimaryKey(id int64) {
	u.ID = id
}

// RuntimeInstall runtime install
type RuntimeInstall struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	ScriptBody string `db:"script_body"`
}

// Meta provides mapping config specific to runtime install
func (ri RuntimeInstall) Meta() MappingConfig {
	return MappingConfig{TableName: "runtime_install"}
}

// PrimaryKey returns primary key
func (ri RuntimeInstall) PrimaryKey() int64 {
	return ri.ID
}

// SetPrimaryKey updates the primary key value after insert
func (ri RuntimeInstall) SetPrimaryKey(id int64) {
	ri.ID = id
}

func (ri RuntimeInstall) String() string {
	return fmt.Sprintf("%d: %s\n", ri.ID, ri.ScriptBody)
}
