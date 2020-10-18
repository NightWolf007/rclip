package cmd

import (
	"context"
	"fmt"

	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var pasteCmd = &cobra.Command{
	Use:     "paste",
	Aliases: []string{"pt", "p"},
	Short:   "Prints content from RClip server",
	PreRun: func(cmd *cobra.Command, args []string) {
		registerViperKey(
			"client.target",
			"CLIENT_TARGET",
			cmd.Flags().Lookup("target"),
			ServerDefaultAddr,
		)
		registerViperKey(
			"client.paste.clipboard",
			"CLIENT_PASTE_CLIPBOARD",
			cmd.Flags().Lookup("clipboard"),
			false,
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		targetAddr := viper.GetString("client.target")

		conn, err := grpc.Dial(targetAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatal().
				Err(err).
				Str("target", targetAddr).
				Msg("Failed to connect to the server")
		}

		defer conn.Close()

		client := api.NewClipboardAPIClient(conn)
		resp, err := client.Get(context.Background(), &api.GetRequest{})
		if err != nil {
			log.Fatal().
				Err(err).
				Str("target", targetAddr).
				Str("method", "Get").
				Msg("Failed to execute RPC method")
		}

		if resp.Value != nil {
			if viper.GetBool("client.paste.clipboard") {
				err := clipboard.Write(resp.Value)
				if err != nil {
					log.Fatal().
						Err(err).
						Bytes("value", resp.Value).
						Msg("Failed to write value to clipboard")
				}
			}

			fmt.Print(resp.Value)
		}
	},
}

func init() {
	pasteCmd.Flags().StringP(
		"target", "t", ServerDefaultAddr,
		"Target server address",
	)
	pasteCmd.Flags().BoolP(
		"clipboard", "c", false,
		"Also paste value to the system clipboard",
	)
}
