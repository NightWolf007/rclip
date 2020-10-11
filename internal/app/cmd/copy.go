package cmd

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	copyListenAddr  string
	copyData        string
	copyToClipboard bool
)

var copyCmd = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"cp", "c"},
	Short:   "Copy content and sends it to RClip server",
	Run: func(cmd *cobra.Command, args []string) {
		data := []byte(copyData)

		if len(copyData) == 0 {
			var err error

			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to read stdin")
			}
		}

		conn, err := grpc.Dial(copyListenAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatal().
				Err(err).
				Str("addr", copyListenAddr).
				Msg("Failed to connect to the server")
		}

		defer conn.Close()

		client := api.NewClipboardAPIClient(conn)
		_, err = client.Push(context.Background(), &api.PushRequest{Value: data})
		if err != nil {
			log.Fatal().
				Err(err).
				Str("method", "Push").
				Msg("Failed to execute RPC method")
		}

		if copyToClipboard {
			clipboard.Write(data)
		}
	},
}

func init() {
	copyCmd.Flags().StringVarP(
		&copyListenAddr, "listen", "l", ServerDefaultAddr,
		"Listen server address",
	)
	copyCmd.Flags().StringVarP(
		&copyData, "data", "d", "",
		"Use the given data instead of stdin",
	)
	pasteCmd.Flags().BoolVarP(
		&copyToClipboard, "clipboard", "c", false,
		"Also copy value to the system clipboard",
	)
}
