package cmd

import (
	"fmt"
	"net"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/NightWolf007/rclip/pb"
	"github.com/NightWolf007/rclip/servers"
)

var cfgFile string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start rclip server",
	Long:  `Start rclip server`,
	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "",
		"config file (default is $HOME/.rclipd.yaml)",
	)
}

// initConfig reads in config file
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
