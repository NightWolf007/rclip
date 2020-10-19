package cmd

import (
	"context"
	"sync"

	"github.com/NightWolf007/rclip/internal/app/syncer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "RClip daemon syncs system clipboard buffer with RClip server",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkClipboardSupport()

		registerViperKey(
			"client.target",
			"CLIENT_TARGET",
			cmd.Flags().Lookup("target"),
			ServerDefaultAddr,
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelFn := context.WithCancel(context.Background())

		rtlLogger := log.With().
			Str("module", "rtl-sync").
			Logger()
		rtlSyncer := syncer.New(rtlLogger)

		ltrLogger := log.With().
			Str("module", "ltr-sync").
			Logger()
		ltrSyncer := syncer.New(ltrLogger)

		wg := sync.WaitGroup{}
		wg.Add(2) // nolint:gomnd // Starting only two goroutines

		go runDaemon(&wg, rtlLogger, func() error {
			return rtlSyncer.RemoteToLocal(ctx)
		})

		go runDaemon(&wg, ltrLogger, func() error {
			return ltrSyncer.RemoteToLocal(ctx)
		})

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

func runDaemon(wg *sync.WaitGroup, logger zerolog.Logger, fn func() error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error().
				Interface("error", r).
				Msg("Recovered from panic")
		}
	}()

	for {
		logger.Debug().
			Msg("Starting deamon")

		err := fn()
		if err == nil {
			logger.Debug().Msg("Shutting down daemon")
			break
		}

		logger.Error().
			Err(err).
			Msg("Received error from daemon")

		logger.Debug().Msg("Restarting daemon")
	}

	logger.Debug().Msg("Deamon is shutted down")
	wg.Done()
}
