package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func daemonCmd() *cobra.Command {
	var listenAddr string

	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "RClip daemon syncs system buffer with RClip server",
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
