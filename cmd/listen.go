package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/NightWolf007/rclip/pb"
)

var listenAddr string
var binaryOutput bool
var connTimeout time.Duration

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen clipboard buffer updates and puts it to stdout",
	Long:  `Listen clipboard buffer updates and puts it to stdout`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(
			listenAddr,
			grpc.WithInsecure(),
			grpc.WithTimeout(connTimeout),
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer conn.Close()

		client := pb.NewClipboardClient(conn)
		stream, err := client.Subscribe(context.Background(), &pb.SubscribeRequest{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for {
			clip, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if binaryOutput {
				fmt.Printf("%X\n", clip.Data)
			} else {
				fmt.Printf("%s\n", clip.Data)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	listenCmd.PersistentFlags().StringVarP(
		&listenAddr, "addr", "a", "localhost:9889",
		"RClip server address",
	)

	listenCmd.PersistentFlags().DurationVarP(
		&connTimeout, "timeout", "t", 5*time.Second,
		"RClip connection timeout",
	)

	listenCmd.PersistentFlags().BoolVarP(
		&binaryOutput, "binary", "b", false,
		"Print hex when set",
	)
}
