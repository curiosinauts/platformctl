package database

type User struct {
	ID          int64  `db:"id"`
	GoogleID    string `db:"google_id"`
	Username    string `db:"username"`
	Password    string `db:"password"`
	Email       string `db:"email"`
	HashedEmail string `db:"hashed_email"`
	IsActive    bool   `db:"is_active"`
	PrivateKey  string `db:"private_key"`
	PublicKey   string `db:"public_key"`
	DockerTag   string
}

type UserRepo struct {
	ID     int64  `db:"id"`
	URI    string `db:"uri"`
	UserID int64  `db:"user_id"`
}

type UserIDE struct {
	ID     int64 `db:"id"`
	UserID int64 `db:"user_id"`
	IDEID  int64 `db:"ide_id"`
}

type IDE struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type IDERepo struct {
	ID        int64  `db:"id"`
	UserIDEID int64  `db:"user_ide_id"`
	URI       string `db:"uri"`
}

type IDERuntimeInstall struct {
	ID               int64 `db:"id"`
	UserIDEID        int64 `db:"user_ide_id"`
	RuntimeInstallID int64 `db:"runtime_install_id"`
}

type RuntimeInstall struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	ScriptBody string `db:"script_body"`
}
