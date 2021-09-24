package database

import (
	"testing"

	"github.com/curiosinauts/platformctl/pkg/testutil"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func init() {
	testutil.InitConfig()
}

func Test(t *testing.T) {

	connStr := viper.Get("database_conn").(string)
	newdb, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatal("test db connection failed:", err.Error())
	}
	newdb.Exec(`set search_path='curiosity'`)

	if err != nil {
		t.Fatal("test db connection failed")
	}
	dbs := NewUserService(newdb)

	user := User{
		Username:    "username",
		Password:    "password",
		Email:       "email",
		HashedEmail: "hashed_email",
		PrivateKey:  "private_key",
		PublicKey:   "public_key",
		IsActive:    true,
	}

	dbs.Save(&user)
	dbs.Del(&user)
}
