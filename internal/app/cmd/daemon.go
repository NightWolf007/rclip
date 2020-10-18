package cmd

import (
	"context"
	"sync"

	"github.com/NightWolf007/rclip/internal/app/syncer"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "RClip daemon syncs system clipboard buffer with RClip server",
	PreRun: func(cmd *cobra.Command, args []string) {
		registerViperKey(
			"client.target",
			"CLIENT_TARGET",
			cmd.Flags().Lookup("target"),
			ServerDefaultAddr,
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !clipboard.IsSupported() {
			log.Fatal().Msg("System clipboard is unsupported")
		}

		ctx, cancelFn := context.WithCancel(context.Background())

		targetAddr := viper.GetString("client.target")

		rtlLogger := log.With().
			Str("module", "rtl-sync").
			Str("taget", targetAddr).
			Logger()
		rtlSyncer := syncer.New(targetAddr, rtlLogger)

		ltrLogger := log.With().
			Str("module", "ltr-sync").
			Str("target", targetAddr).
			Logger()
		ltrSyncer := syncer.New(targetAddr, ltrLogger)

		wg := sync.WaitGroup{}
		wg.Add(2) // nolint:gomnd // Starting only two goroutines

		go runDaemon(func() error {
			return rtlSyncer.RemoteToLocal(ctx)
		}, &wg, rtlLogger)

		go runDaemon(func() error {
			return ltrSyncer.RemoteToLocal(ctx)
		}, &wg, ltrLogger)

		grc := newGrace()
		grc.Shutdown = cancelFn
		grc.Wait = wg.Wait

		grc.Run()
	},
}

func init() {
	daemonCmd.Flags().StringP(
		"target", "t", ServerDefaultAddr,
		"Target server address",
	)
}

func runDaemon(fn func() error, wg *sync.WaitGroup, logger zerolog.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error().
				Interface("error", r).
				Msg("Daemon recovered from panic")
		}
	}()

	for {
		err := fn()
		if err == nil {
			break
		}

		logger.Error().
			Err(err).
			Msg("Received error from daemon")

		logger.Debug().Msg("Restarting daemon")
	}

	wg.Done()
}
