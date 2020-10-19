package cmd

import (
	"time"

	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/NightWolf007/rclip/internal/pkg/grace"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func checkClipboardSupport() {
	if !clipboard.IsSupported() {
		log.Fatal().
			Msg("System clipboard is unsupported")
	}
}

func errF(err error) {
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Received fatal error")
	}
}

func newGrace() *grace.Grace {
	return &grace.Grace{
		Timeout: 5 * time.Second, // nolint:gomnd // This is the default value
		OnShutdown: func() {
			log.Info().Msg("Gracefully stopping... (press Ctrl+C again to force)")
		},
		OnDone: func() {
			log.Info().Msg("Bye!")
		},
		OnForceQuit: func() {
			log.Info().Msg("Force stopping...")
			log.Info().Msg("Bye!")
		},
		OnTimeout: func() {
			log.Info().Msg("Gracefully stop is taking too long. Force stopping...")
			log.Info().Msg("Bye!")
		},
	}
}

func registerViperKey(key string, env string, flag *pflag.Flag, defValue interface{}) {
	viper.SetDefault(key, defValue)

	err := viper.BindEnv(key, env)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("key", key).
			Str("env", env).
			Msg("Failed to bind env variable")
	}

	err = viper.BindPFlag(key, flag)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("key", key).
			Interface("flag", flag).
			Msg("Failed to bind flag")
	}
}

func writeToClipboard(val []byte) error {
	log := log.With().Bytes("value", val).Logger()

	log.Debug().
		Msg("Writing value to the system clipboard")

	err := clipboard.Write(val)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to write value to clipboard")

		return err
	}

	log.Debug().
		Msg("System clipboard successfully updated")

	return nil
}
