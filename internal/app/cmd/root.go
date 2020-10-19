package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to execute root cmd")
	}
}

var rootOpts = struct {
	logVerbose bool
	logPretty  bool
	configPath string
}{}

var rootCmd = &cobra.Command{
	Use:   "rclip",
	Short: "RClip remote clipboard",
	Long:  `RClip is a virtual remote clipboard based on clients-server architecture.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLocale()
		initLogger()
		initConfig()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&rootOpts.logVerbose, "verbose", "v", false,
		"Enable debug log level",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&rootOpts.logPretty, "pretty", "p", false,
		"Print logs in a colorized, human-friendly format",
	)
	rootCmd.PersistentFlags().StringVarP(
		&rootOpts.configPath, "config", "c", "",
		"Path to .toml config file",
	)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(pasteCmd)
	rootCmd.AddCommand(histCmd)
	rootCmd.AddCommand(daemonCmd)
}

func initLocale() {
	time.Local = time.UTC
}

func initLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if rootOpts.logVerbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if rootOpts.logPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func initConfig() {
	if rootOpts.configPath != "" {
		viper.SetConfigFile(rootOpts.configPath)
	} else {
		viper.SetConfigName("rclip")
		viper.SetConfigType("toml")

		viper.AddConfigPath("/etc/rclip/")
		viper.AddConfigPath("$XDG_CONFIG_HOME/rclip/")
		viper.AddConfigPath("$HOME/.rclip/")
		viper.AddConfigPath("./")
	}

	log := log.With().Str("config_path", viper.ConfigFileUsed()).Logger()

	log.Debug().Msg("Loading config file")

	if err := viper.ReadInConfig(); err != nil {
		log.Debug().
			Err(err).
			Msg("Failed to load config file")
	} else {
		log.Debug().
			Interface("config", viper.AllSettings()).
			Msg("Config file loaded")
	}
}
