/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"log"
	"os"

	bb "github.com/I0HuKc/baitbotnotbytebot/internal/bot"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db/sqlstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "r",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		bot, err := tgbotapi.NewBotAPI(os.Getenv("APP_BOT_TOKEN"))
		if err != nil {
			log.Panic(err)
		}

		pg, err := db.SetPgConn(ctx, os.Getenv("APP_DB_URL"))
		if err != nil {
			log.Panic(err)
		}
		defer pg.Close()

		baitbot := bb.CreateBaitbot(bot, sqlstore.CreateSqlStore(pg))
		if err := baitbot.Serve(ctx); err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
