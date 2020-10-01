package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func serverCmd() *cobra.Command {
	var bindAddr string

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start RClip server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Not implemented yet")
		},
	}

	cmd.Flags().StringVarP(
		&bindAddr, "bind", "b", "localhost:9889",
		"Bind address",
	)

	return cmd
}
