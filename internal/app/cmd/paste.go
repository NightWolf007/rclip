package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func pasteCmd() *cobra.Command {
	var listenAddr string

	cmd := &cobra.Command{
		Use:     "paste",
		Aliases: []string{"pt", "p"},
		Short:   "Prints content from RClip server",
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
