package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/jmoiron/sqlx"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool
var db *sqlx.DB

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "platformctl",
	Short: "CuriosityWorks learning platform tool",
	Long:  `CuriosityWorks learning platform tool. The tool helps in managing users and in automating the code server build, etc`,
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
	connStr := viper.Get("database_conn").(string)
	newdb, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	db = newdb
	db.Exec(`set search_path='curiosity'`)
	if debug {
		fmt.Printf("database_conn : %s\n", viper.Get("database_conn"))
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

	viper.SetEnvPrefix("platformctl")
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
			msg.Info("step = [" + step + "]" + " error = [" + e.Err.Error() + "]")
			msg.Failure(eh.message)
			os.Exit(1)
		} else {
			return
		}
	}

	if err != nil {
		msg.Info("step = [" + step + "]" + " error = [" + err.Error() + "]")
		msg.Failure(eh.message)
		os.Exit(1)
	}
}
