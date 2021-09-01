package database

import "fmt"

// User user
type User struct {
	ID              int64  `db:"id"`
	GoogleID        string `db:"google_id"`
	Username        string `db:"username"`
	Password        string `db:"password"`
	Email           string `db:"email"`
	HashedEmail     string `db:"hashed_email"`
	IsActive        bool   `db:"is_active"`
	PrivateKey      string `db:"private_key"`
	PublicKey       string `db:"public_key"`
	PublicKeyID     int64  `db:"public_key_id"`
	DockerTag       string
	IDEs            []IDE
	RuntimeInstalls []string
	Repos           []string
}

func (u User) String() string {
	return fmt.Sprintf("%d: %s", u.ID, u.Username)
}

// Meta provides mapping config specific to user
func (u *User) Meta() MappingConfig {
	return MappingConfig{TableName: "curiosity.user"}
}

// PrimaryKey returns primary key
func (u *User) PrimaryKey() int64 {
	return u.ID
}

// SetPrimaryKey updates the primary key value after insert
func (u *User) SetPrimaryKey(id int64) {
	u.ID = id
}

// UserRepo user repo
type UserRepo struct {
	ID     int64  `db:"id"`
	URI    string `db:"uri"`
	UserID int64  `db:"user_id"`
}

// NewUserRepo returns new instance of user repo
func NewUserRepo(uri string, userID int64) *UserRepo {
	return &UserRepo{URI: uri, UserID: userID}
}

// Meta provides mapping config specific to user
func (ur *UserRepo) Meta() MappingConfig {
	return MappingConfig{TableName: "user_repo"}
}

// PrimaryKey returns primary key
func (ur *UserRepo) PrimaryKey() int64 {
	return ur.ID
}

// SetPrimaryKey updates the primary key value after insert
func (ur *UserRepo) SetPrimaryKey(id int64) {
	ur.ID = id
}

// UserIDE user ide
type UserIDE struct {
	ID     int64 `db:"id"`
	UserID int64 `db:"user_id"`
	IDEID  int64 `db:"ide_id"`
}

// Meta provides mapping config specific to user ide
func (ur *UserIDE) Meta() MappingConfig {
	return MappingConfig{TableName: "user_ide"}
}

// PrimaryKey returns primary key
func (ur *UserIDE) PrimaryKey() int64 {
	return ur.ID
}

// SetPrimaryKey updates the primary key value after insert
func (ur *UserIDE) SetPrimaryKey(id int64) {
	ur.ID = id
}

// IDE ide
type IDE struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

// Meta provides mapping config specific to user
func (ur *IDE) Meta() MappingConfig {
	return MappingConfig{TableName: "ide"}
}

// PrimaryKey returns primary key
func (ur *IDE) PrimaryKey() int64 {
	return ur.ID
}

// SetPrimaryKey updates the primary key value after insert
func (ur *IDE) SetPrimaryKey(id int64) {
	ur.ID = id
}

// IDERepo ide repo
type IDERepo struct {
	ID        int64  `db:"id"`
	UserIDEID int64  `db:"user_ide_id"`
	URI       string `db:"uri"`
}

// Meta provides mapping config specific to ide repo
func (ur *IDERepo) Meta() MappingConfig {
	return MappingConfig{TableName: "ide_repo"}
}

// PrimaryKey returns primary key
func (ur *IDERepo) PrimaryKey() int64 {
	return ur.ID
}

// SetPrimaryKey updates the primary key value after insert
func (ur *IDERepo) SetPrimaryKey(id int64) {
	ur.ID = id
}

// IDERuntimeInstall ide runtime install
type IDERuntimeInstall struct {
	ID               int64 `db:"id"`
	UserIDEID        int64 `db:"user_ide_id"`
	RuntimeInstallID int64 `db:"runtime_install_id"`
}

// Meta provides mapping config specific to ide runtime install
func (ur *IDERuntimeInstall) Meta() MappingConfig {
	return MappingConfig{TableName: "ide_runtime_install"}
}

// PrimaryKey returns primary key
func (ur *IDERuntimeInstall) PrimaryKey() int64 {
	return ur.ID
}

// SetPrimaryKey updates the primary key value after insert
func (ur *IDERuntimeInstall) SetPrimaryKey(id int64) {
	ur.ID = id
}

// RuntimeInstall runtime install
type RuntimeInstall struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	ScriptBody string `db:"script_body"`
}

// Meta provides mapping config specific to runtime install
func (ur *RuntimeInstall) Meta() MappingConfig {
	return MappingConfig{TableName: "runtime_install"}
}

// PrimaryKey returns primary key
func (ur *RuntimeInstall) PrimaryKey() int64 {
	return ur.ID
}

// SetPrimaryKey updates the primary key value after insert
func (ur *RuntimeInstall) SetPrimaryKey(id int64) {
	ur.ID = id
}
