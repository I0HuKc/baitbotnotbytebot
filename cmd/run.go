package cmd

import (
	"context"
	"log"

	"github.com/I0HuKc/baitbotnotbytebot/internal/bot"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "r",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		rclient, err := db.SetRedisConn(ctx)
		if err != nil {
			log.Panic(err)
		}
		defer rclient.Close()

		bot, err := bot.NewBaitbot(rclient)
		if err != nil {
			log.Fatal(err)
		}

		bot.Configure(ctx).Serve()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
