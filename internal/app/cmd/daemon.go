package cmd

import (
	"context"
	"sync"

	"github.com/NightWolf007/rclip/internal/app/syncer"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var daemonListenAddr string

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "RClip daemon syncs system clipboard buffer with RClip server",
	Run: func(cmd *cobra.Command, args []string) {
		if !clipboard.IsSupported() {
			log.Fatal().Msg("System clipboard unsupported")
		}

		ctx, cancelFn := context.WithCancel(context.Background())

		rtlLogger := log.With().
			Str("module", "rtl-sync").
			Str("addr", daemonListenAddr).
			Logger()
		rtlSyncer := syncer.New(daemonListenAddr, rtlLogger)

		ltrLogger := log.With().
			Str("module", "ltr-sync").
			Str("addr", daemonListenAddr).
			Logger()
		ltrSyncer := syncer.New(daemonListenAddr, ltrLogger)

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
	daemonCmd.Flags().StringVarP(
		&daemonListenAddr, "listen", "l", ServerDefaultAddr,
		"Listen server address",
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
