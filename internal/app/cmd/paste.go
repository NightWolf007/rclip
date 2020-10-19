package cmd

import (
	"context"
	"fmt"

	"github.com/NightWolf007/rclip/internal/app/client"
	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/spf13/cobra"
)

var pasteOpts = struct {
	useClipboard bool
}{}

var pasteCmd = &cobra.Command{
	Use:     "paste",
	Aliases: []string{"pt", "p"},
	Short:   "Prints content from RClip server",
	PreRun: func(cmd *cobra.Command, args []string) {
		if pasteOpts.useClipboard {
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
		cli, err := client.Dial()
		errF(err)

		defer cli.Close()

		resp, err := cli.Get(context.Background(), &api.GetRequest{})
		errF(err)

		if pasteOpts.useClipboard {
			errF(writeToClipboard(resp.Value))
		}

		fmt.Printf("%s", resp.Value)
	},
}

func init() {
	pasteCmd.Flags().StringP(
		"target", "t", ServerDefaultAddr,
		"Target server address",
	)
	pasteCmd.Flags().BoolVarP(
		&pasteOpts.useClipboard, "write", "w", false,
		"Also write value to the system clipboard",
	)
}
