package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/NightWolf007/rclip/pb"
)

var copyAddr string
var copyData string

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Send data to rclip server",
	Long:  `Send data to rclip server`,
	Run: func(cmd *cobra.Command, args []string) {
		data := []byte(copyData)
		if len(data) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				data = append(data, scanner.Bytes()...)
				data = append(data, '\n')
			}
			data = data[:len(data)-1]
		}

		conn, err := grpc.Dial(copyAddr, grpc.WithInsecure())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer conn.Close()

		client := pb.NewClipboardClient(conn)
		_, err = client.Push(context.Background(), &pb.PushRequest{Data: data})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.PersistentFlags().StringVarP(
		&copyAddr, "addr", "a", "localhost:8000",
		"RClip server address",
	)
	copyCmd.Flags().StringVarP(
		&copyData, "data", "d", "",
		"Use given data instead of stdin",
	)
}
