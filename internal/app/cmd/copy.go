package cmd

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/NightWolf007/rclip/internal/app/client"
	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var copyOpts = struct {
	data         string
	useClipboard bool
}{}

var copyCmd = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"cp", "c"},
	Short:   "Copy content and sends it to RClip server",
	PreRun: func(cmd *cobra.Command, args []string) {
		if copyOpts.useClipboard {
			checkClipboardSupport()
		}

		registerViperKey(
			"client.target",
			"CLIENT_TARGET",
			cmd.Flags().Lookup("target"),
			ServerDefaultAddr,
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		data := []byte(copyOpts.data)

		if len(copyOpts.data) == 0 {
			var err error

			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to read stdin")
			}
		}

		cli, err := client.Dial()
		errF(err)

		defer cli.Close()

		_, err = cli.Push(context.Background(), &api.PushRequest{Value: data})
		errF(err)

		if copyOpts.useClipboard {
			errF(writeToClipboard(data))
		}
	},
}

func init() {
	copyCmd.Flags().StringP(
		"target", "t", ServerDefaultAddr,
		"Target server address",
	)
	copyCmd.Flags().StringVarP(
		&copyOpts.data, "data", "d", "",
		"Use the given data instead of stdin",
	)
	copyCmd.Flags().BoolVarP(
		&copyOpts.useClipboard, "write", "w", false,
		"Also write value to the system clipboard",
	)
}
