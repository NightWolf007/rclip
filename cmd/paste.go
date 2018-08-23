package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/NightWolf007/rclip/pb"
)

var pasteAddr string

// pasteCmd represents the paste command
var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "Puts data from RClip server into stdout",
	Long:  `Puts data from RClip server into stdout`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(pasteAddr, grpc.WithInsecure())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer conn.Close()

		client := pb.NewClipboardClient(conn)
		resp, err := client.Get(context.Background(), &pb.GetRequest{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(resp.Data)
	},
}

func init() {
	rootCmd.AddCommand(pasteCmd)

	pasteCmd.PersistentFlags().StringVarP(
		&pasteAddr, "addr", "a", "localhost:8000",
		"RClip server address",
	)
}
