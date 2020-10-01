package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func copyCmd() *cobra.Command {
	var listenAddr string

	cmd := &cobra.Command{
		Use:     "copy",
		Aliases: []string{"cp", "c"},
		Short:   "Copy content and sends it to RClip server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Not implemented yet")
		},
	}

	cmd.Flags().StringVarP(
		&listenAddr, "listen", "l", "localhost:9889",
		"Listen server address",
	)

	return cmd
}
