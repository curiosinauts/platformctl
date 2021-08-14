package cmd

import (
	"errors"
	"fmt"
	"github.com/curiosinauts/platformctl/internal/database"
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/jmoiron/sqlx"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool
var db *sqlx.DB

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "platformctl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.platformctl.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debugging")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initDB() {
	connStr := viper.Get("database.conn").(string)
	newdb, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	db = newdb
	db.Exec(`set search_path='curiosity'`)
	if debug {
		fmt.Printf("database.conn : %s\n", viper.Get("database.conn"))
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".platformctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".platformctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && debug {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	initDB()
}

type ErrorHandler struct {
	message string
}

func (eh ErrorHandler) HandleError(step string, err error) {
	var e *database.DBError
	if errors.As(err, &e) {
		if e != nil {
			msg.Info(step)
			msg.Failure(eh.message)
			msg.Formaterr(e.Err)
			os.Exit(1)
		} else {
			return
		}
	}

	if err != nil {
		msg.Info(step)
		msg.Failure(eh.message)
		msg.Formaterr(err)
		os.Exit(1)
	}
}
