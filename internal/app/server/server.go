// Package server implements RClip server.
package server

import (
	"net"

	"github.com/rs/zerolog/log"
)

func Run(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	messages := make([]string)

	defer l.Close()

	for {
		c, err := l.Accept()

		log.Debug().Msg("Received incomming connection")

		if err != nil {
			log.Error().Err(err).Msg("Cannot accept connection")

			continue
		}

		log.Debug().Str("address", c.LocalAddr().String()).Msg("Connection established")

		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {

}
