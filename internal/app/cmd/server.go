package cmd

import (
	"net"

	"github.com/NightWolf007/rclip/internal/app/servers"
	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// ServerDefaultAddr is a default server listen address.
const ServerDefaultAddr = "localhost:9889"

var (
	serverBindAddr    string
	serverHistorySize uint
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start RClip server",
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", serverBindAddr)
		if err != nil {
			log.Fatal().
				Err(err).
				Str("addr", serverBindAddr).
				Msg("Failed to listen addr")
		}

		server := grpc.NewServer()

		clipboardServer := servers.NewClipboardServer(serverHistorySize)
		api.RegisterClipboardAPIServer(server, clipboardServer)

		err = server.Serve(lis)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to serve GRPC server")
		}

		grc := newGrace()
		grc.Shutdown = server.GracefulStop

		grc.Run()
	},
}

func init() {
	serverCmd.Flags().StringVarP(
		&serverBindAddr, "bind", "b", ServerDefaultAddr,
		"Bind address",
	)
	serverCmd.Flags().UintVarP(
		&serverHistorySize, "hist-size", "s", 100,
		"Maximum size of clipboard history",
	)
}
