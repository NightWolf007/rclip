package cmd

import (
	"context"

	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	daemonListenAddr string
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "RClip daemon syncs system clipboard buffer with RClip server",
	Run: func(cmd *cobra.Command, args []string) {
		if !clipboard.IsSupported() {
			log.Fatal().Msg("System clipboard unsupported")
		}

		ctx := context.Background()

		go daemonListenClipboard(ctx)

		go daemonListenRemote(ctx)
	},
}

func init() {
	daemonCmd.Flags().StringVarP(
		&daemonListenAddr, "listen", "l", ServerDefaultAddr,
		"Listen server address",
	)
}
