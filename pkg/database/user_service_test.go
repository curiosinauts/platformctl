package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// init reads in config file and ENV variables if set.
func Init() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".platformctl" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".platformctl")

	viper.SetEnvPrefix("platformctl")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func Test(t *testing.T) {

	Init()

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
