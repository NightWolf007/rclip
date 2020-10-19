package cmd

import (
	"net"
	"sync"

	"github.com/NightWolf007/rclip/internal/app/server/interceptors"
	"github.com/NightWolf007/rclip/internal/app/server/servers"
	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// ServerDefaultAddr is a default server listen address.
const ServerDefaultAddr = "localhost:9889"

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start RClip server",
	PreRun: func(cmd *cobra.Command, args []string) {
		registerViperKey(
			"server.listen",
			"SERVER_LISTEN",
			cmd.Flags().Lookup("listen"),
			ServerDefaultAddr,
		)

		registerViperKey(
			"server.hist_size",
			"SERVER_HIST_SIZE",
			cmd.Flags().Lookup("hist-size"),
			100,
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		listenAddr := viper.GetString("server.listen")
		histSize := viper.GetUint("server.hist_size")

		log := log.With().
			Str("listen", listenAddr).
			Uint("hist_size", histSize).
			Logger()

		log.Debug().
			Msg("Starting listenning address")

		lis, err := net.Listen("tcp", listenAddr)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to listen address")
		}

		log.Debug().
			Msg("Initializing server")

		server := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptors.UnaryLogInterceptor(log),
			),
			grpc.ChainStreamInterceptor(
				interceptors.StreamLogInterceptor(log),
			),
		)

		clipboardServer := servers.NewClipboardServer(histSize)
		api.RegisterClipboardAPIServer(server, clipboardServer)

		log.Debug().
			Msg("Starting server")

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			err := server.Serve(lis)
			if err != nil {
				log.Fatal().
					Err(err).
					Str("listen", listenAddr).
					Uint("hist_size", histSize).
					Msg("Failed to serve GRPC server")
			}

			log.Debug().
				Msg("Server is stopped")

			wg.Done()
		}()

		log.Info().
			Msg("Server is running")

		grc := newGrace()
		grc.Shutdown = server.GracefulStop
		grc.Wait = wg.Wait

		grc.Run()
	},
}

func init() {
	serverCmd.Flags().StringP(
		"listen", "l", ServerDefaultAddr,
		"Listen address",
	)
	serverCmd.Flags().UintP(
		"hist-size", "s", 100,
		"Maximum size of clipboard history",
	)
}
