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
var userService database.UserService
var dbs database.UserService

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
	// rely on connectiont string to set the search_path instead of hard coding it here
	// db.Exec(`set search_path='platformctl'`)

	userService = database.NewUserService(db)

	options := []database.DBOption{
		database.DBOptionDebug(debug),
	}

	dbs = database.NewUserServiceWithOptions(db, options...)
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
		msg.Info("using config file:" + viper.ConfigFileUsed())
		fmt.Println()
	}
	initDB()
}

// ErrorHandler handles errors
type ErrorHandler struct {
	message string
}

// HandleError handles the given error
func (eh ErrorHandler) HandleError(step string, err error) {
	if debug {
		fmt.Println("step:", step+"\n")
	}
	var e *database.DBError
	if errors.As(err, &e) {
		if e != nil {
			msg.Info("step  = " + step)
			msg.Info("error = " + e.Err.Error())
			msg.Info("sql   = " + e.Query)
			msg.Failure(eh.message)
			os.Exit(1)
		} else {
			return
		}
	}

	if err != nil {
		msg.Info("step  = " + step)
		msg.Info("error = " + err.Error())
		msg.Failure(eh.message)
		os.Exit(1)
	}
}

// HandleErrorWithOutput handles the given error and prints output
func (eh ErrorHandler) HandleErrorWithOutput(step string, err error, output string) {
	var e *database.DBError
	if errors.As(err, &e) {
		if e != nil {
			msg.Info("step = [" + step + "]" + " error = [" + e.Err.Error() + "] output = [" + output + "]")
			msg.Failure(eh.message)
			os.Exit(1)
		} else {
			return
		}
	}

	if err != nil {
		msg.Info("step = [" + step + "]" + " error = [" + err.Error() + "] output = [" + output + "]")
		msg.Failure(eh.message)
		os.Exit(1)
	}
}

// PrintErrorWithOutput prints the given error and output
func (eh ErrorHandler) PrintErrorWithOutput(step string, err error, output string) {
	var e *database.DBError
	if errors.As(err, &e) {
		if e != nil {
			msg.Info("step = [" + step + "]" + " error = [" + e.Err.Error() + "] output = [" + output + "]")
		} else {
			return
		}
	}

	if err != nil {
		msg.Info("step = [" + step + "]" + " error = [" + err.Error() + "] output = [" + output + "]")
	}
}

// PrintError prints the error, no exit
func (eh ErrorHandler) PrintError(step string, err error) {
	var e *database.DBError
	if errors.As(err, &e) {
		if e != nil {
			msg.Info("step = [" + step + "]" + " error = [" + e.Err.Error() + "]")
		} else {
			return
		}
	}

	if err != nil {
		msg.Info("step = [" + step + "]" + " error = [" + err.Error() + "]")
	}
}
