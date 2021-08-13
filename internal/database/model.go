package database

type User struct {
	ID           string `json:"id"              db:"user_id"               api:"users"`
	PlatformName string `json:"-"               db:"platform_name"`
	Email        string `json:"email"           db:"email"                 api:"attr"`
	Password     string `json:"password"        db:"passhash"              api:"attr"`
	FirstName    string `json:"firstname"       db:"firstname"             api:"attr"`
	LastName     string `json:"lastname"        db:"lastname"              api:"attr"`
	Created      int64  `json:"created"         db:"created_date"          api:"attr"`
}
