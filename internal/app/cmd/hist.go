package cmd

import (
	"context"
	"fmt"

	"github.com/NightWolf007/rclip/internal/app/client"
	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var histOpts = struct {
	size         int
	useClipboard bool
}{}

var histCmd = &cobra.Command{
	Use:     "hist",
	Aliases: []string{"ht", "h"},
	Short:   "Shows content of RClip server history",
	PreRun: func(cmd *cobra.Command, args []string) {
		if histOpts.useClipboard {
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

		resp, err := cli.Hist(context.Background(), &api.HistRequest{})
		errF(err)

		selItems := make([]string, len(resp.Values))
		for i, val := range resp.Values {
			maxLen := 100
			if maxLen > len(val) {
				maxLen = len(val)
			}

			selItems[i] = string(val[0:maxLen])
		}

		sel := promptui.Select{
			Label:        "Select value",
			Size:         histOpts.size,
			Items:        selItems,
			HideSelected: true,
		}

		_, selResult, err := sel.Run()
		errF(err)

		if histOpts.useClipboard {
			errF(writeToClipboard([]byte(selResult)))
		}

		fmt.Printf("%s", selResult)
	},
}

func init() {
	histCmd.Flags().StringP(
		"target", "t", ServerDefaultAddr,
		"Target server address",
	)
	histCmd.Flags().IntVarP(
		&histOpts.size, "size", "s", 10,
		"The number of items that should appear on the select",
	)
	histCmd.Flags().BoolVarP(
		&histOpts.useClipboard, "write", "w", false,
		"Also write value to the system clipboard",
	)
}
