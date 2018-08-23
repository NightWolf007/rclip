package cmd

import (
	"fmt"
	"net"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/NightWolf007/rclipd/pb"
	"github.com/NightWolf007/rclipd/servers"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rclipd",
	Short: "Teminal clipboard",
	Long:  `Terminal clipboard. Includes server and client.`,
	Run:   run,
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

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "",
		"config file (default is $HOME/.rclipd.yaml)",
	)
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

		// Search config in home directory with name ".rclipd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rclipd")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func run(cmd *cobra.Command, args []string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	listenAddr := viper.GetString("listen")
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot start TCP server on %s", listenAddr)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterClipboardServer(grpcServer, servers.NewClipboardServer())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal().Err(err).Msg("Server unexpected shutdown")
	}
}
