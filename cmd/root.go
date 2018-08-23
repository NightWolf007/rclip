package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var debugFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rclip",
	Short: "Teminal clipboard",
	Long:  `Terminal clipboard. Includes server and client.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if debugFlag {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	},
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
	rootCmd.PersistentFlags().BoolVar(
		&debugFlag, "debug", false,
		"Enable debug output",
	)
}
