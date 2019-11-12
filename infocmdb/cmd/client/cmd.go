package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	cfgFile string
	url string
	username string
	password string
	apikey string
	output string
	limit int
	)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "A brief description of your application #TODO",
	Long: `loooong explanation... #TODO`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	//PersistentPreRun: func(cmd *cobra.Command, args []string) {
	//	// Run parent persistent pre run
	//	if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
	//		cmd.Parent().PersistentPreRun(, args)
	//	}

	//},
}
// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&url, "url", "", "api url")
	if err := viper.BindPFlag("api.url", rootCmd.PersistentFlags().Lookup("url")); err != nil {
		log.Fatalf("error binding flag set, %v", err)
	}
	if err := rootCmd.MarkPersistentFlagRequired("url"); err != nil{
		log.Fatalf("error setting mandatory flag, %v", err)
	}

	rootCmd.PersistentFlags().StringVar(&username, "username", viper.GetString("api.username"), "admin")
	if err := viper.BindPFlag("api.username", rootCmd.PersistentFlags().Lookup("username")); err != nil {
		log.Fatalf("error binding flag set, %v", err)
	}

	rootCmd.PersistentFlags().StringVar(&password, "password", viper.GetString("api.password"), "admin")
	if err := viper.BindPFlag("api.password", rootCmd.PersistentFlags().Lookup("password")); err != nil {
		log.Fatalf("error binding flag set, %v", err)
	}

	rootCmd.PersistentFlags().StringVar(&apikey, "apikey", viper.GetString("api.token"), "UNSET-API-TOKEN")
	if err := viper.BindPFlag("api.token", rootCmd.PersistentFlags().Lookup("apikey")); err != nil {
		log.Fatalf("error binding flag set, %v", err)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".main.yaml", "config file (default is $HOME/.main.yaml)")
	rootCmd.PersistentFlags().StringVar(&output, "output", "json", "output format")
	rootCmd.PersistentFlags().IntVar(&limit, "limit", -1, "limit number of results to fetch, default unlimited(-1)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".main" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".main")
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if _, err = fmt.Fprintf(os.Stderr, "Using config file: %v\n", viper.ConfigFileUsed()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}






