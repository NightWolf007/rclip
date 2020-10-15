package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/NightWolf007/rclip/internal/pkg/grace"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Println("FATAL! Cannot execute cmd")
	}
}

func rootCmd() *cobra.Command {
	var (
		logVerbose bool
		logPretty  bool
	)

	cmd := &cobra.Command{
		Use:   "rclip",
		Short: "RClip remote clipboard",
		Long:  `RClip is a virtual remote clipboard based on clients-server architecture.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initLocale()
			initLogger(logVerbose, logPretty)
		},
	}

	cmd.PersistentFlags().BoolVarP(
		&logVerbose, "verbose", "v", false,
		"Enable debug log level",
	)
	cmd.PersistentFlags().BoolVarP(
		&logPretty, "pretty", "p", false,
		"Print logs in a colorized, human-friendly format",
	)

	cmd.AddCommand(serverCmd)
	cmd.AddCommand(copyCmd)
	cmd.AddCommand(pasteCmd)
	cmd.AddCommand(daemonCmd)

	return cmd
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

func newGrace() *grace.Grace {
	return &grace.Grace{
		Timeout: 5 * time.Second,
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
