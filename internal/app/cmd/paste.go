package cmd

import (
	"context"
	"fmt"

	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	pasteListenAddr  string
	pasteToClipboard bool
)

var pasteCmd = &cobra.Command{
	Use:     "paste",
	Aliases: []string{"pt", "p"},
	Short:   "Prints content from RClip server",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(pasteListenAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatal().
				Err(err).
				Str("addr", daemonListenAddr).
				Msg("Failed to connect to the server")
		}

		defer conn.Close()

		client := api.NewClipboardAPIClient(conn)
		resp, err := client.Get(context.Background(), &api.GetRequest{})
		if err != nil {
			log.Fatal().
				Err(err).
				Str("method", "Get").
				Msg("Failed to execute RPC method")
		}

		if resp.Value != nil {
			if pasteToClipboard {
				clipboard.Write(resp.Value)
			}

			fmt.Print(resp.Value)
		}
	},
}

func init() {
	pasteCmd.Flags().StringVarP(
		&pasteListenAddr, "listen", "l", ServerDefaultAddr,
		"Listen server address",
	)
	pasteCmd.Flags().BoolVarP(
		&pasteToClipboard, "clipboard", "c", false,
		"Also paste value to the system clipboard",
	)
}
