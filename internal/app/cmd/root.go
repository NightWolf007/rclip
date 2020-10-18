package cmd

import (
	"os"
	"time"

	"github.com/NightWolf007/rclip/internal/pkg/grace"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

var rootCmd = &cobra.Command{
	Use:   "rclip",
	Short: "RClip remote clipboard",
	Long:  `RClip is a virtual remote clipboard based on clients-server architecture.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		registerViperKey(
			"log.verbose",
			"LOG_VERBOSE",
			cmd.PersistentFlags().Lookup("verbose"),
			false,
		)
		registerViperKey(
			"log.pretty",
			"LOG_pretty",
			cmd.PersistentFlags().Lookup("pretty"),
			false,
		)

		initLocale()
		initLogger(
			viper.GetBool("log.verbose"),
			viper.GetBool("log.pretty"),
		)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP(
		"verbose", "v", false,
		"Enable debug log level",
	)
	rootCmd.PersistentFlags().BoolP(
		"pretty", "p", false,
		"Print logs in a colorized, human-friendly format",
	)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(pasteCmd)
	rootCmd.AddCommand(daemonCmd)
}

func initLocale() {
	time.Local = time.UTC
}

func initLogger(verbose, pretty bool) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func registerViperKey(key string, env string, flag *pflag.Flag, defValue interface{}) {
	viper.SetDefault(key, defValue)

	err := viper.BindEnv(key, env)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("key", key).
			Str("env", env).
			Msg("Failed to bind env variable")
	}

	err = viper.BindPFlag(key, flag)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("key", key).
			Interface("flag", flag).
			Msg("Failed to bind flag")
	}
}

func newGrace() *grace.Grace {
	return &grace.Grace{
		Timeout: 5 * time.Second, // nolint:gomnd // This is the default value
		OnShutdown: func() {
			log.Info().Msg("Gracefully stopping... (press Ctrl+C again to force)")
		},
		OnDone: func() {
			log.Info().Msg("Bye!")
		},
		OnForceQuit: func() {
			log.Info().Msg("Force stopping...")
			log.Info().Msg("Bye!")
		},
		OnTimeout: func() {
			log.Info().Msg("Gracefully stop is taking too long. Force stopping...")
			log.Info().Msg("Bye!")
		},
	}
}
