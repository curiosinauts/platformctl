package database

type User struct {
	ID          int    `db:"id"`
	GoogleID    string `db:"google_id"`
	Username    string `db:"username"`
	Password    string `db:"password"`
	Email       string `db:"email"`
	HashedEmail string `db:"hashed_email"`
	IsActive    bool   `db:"is_active"`
	PrivateKey  string `db:"private_key"`
	PublicKey   string `db:"public_key"`
}
