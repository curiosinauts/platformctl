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
}

type UserRepo struct {
	ID     int64  `db:"id"`
	URI    string `db:"uri"`
	UserID int64  `db:"user_id"`
}
