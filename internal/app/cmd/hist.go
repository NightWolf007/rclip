package cmd

import (
	"context"
	"fmt"

	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var histCmd = &cobra.Command{
	Use:     "hist",
	Aliases: []string{"ht", "h"},
	Short:   "Shows content of RClip server history",
	PreRun: func(cmd *cobra.Command, args []string) {
		registerViperKey(
			"client.target",
			"CLIENT_TARGET",
			cmd.Flags().Lookup("target"),
			ServerDefaultAddr,
		)
		registerViperKey(
			"client.hist.clipboard",
			"CLIENT_HIST_CLIPBOARD",
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
		resp, err := client.Hist(context.Background(), &api.HistRequest{})
		if err != nil {
			log.Fatal().
				Err(err).
				Str("target", targetAddr).
				Str("method", "Get").
				Msg("Failed to execute RPC method")
		}

		selItems := make([]string, len(resp.Values))
		for i, val := range resp.Values {
			selItems[i] = string(val[0:100])
		}

		sel := promptui.Select{
			Size:  10,
			Items: selItems,
		}

		_, selResult, err := sel.Run()
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to execute select")
		}

		if viper.GetBool("client.hist.clipboard") {
			err := clipboard.Write([]byte(selResult))
			if err != nil {
				log.Fatal().
					Err(err).
					Str("value", selResult).
					Msg("Failed to write value to clipboard")
			}
		}

		fmt.Print(selResult)
	},
}

func init() {
	histCmd.Flags().StringP(
		"target", "t", ServerDefaultAddr,
		"Target server address",
	)
	histCmd.Flags().BoolP(
		"clipboard", "c", false,
		"Also paste value to the system clipboard",
	)
}
