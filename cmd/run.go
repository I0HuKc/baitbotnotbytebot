/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"

	bb "github.com/I0HuKc/baitbotnotbytebot/internal/bot"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "r",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		bot, err := tgbotapi.NewBotAPI(os.Getenv("APP_BOT_TOKEN"))
		if err != nil {
			log.Panic(err)
		}

		pg, err := db.SetPgConn(os.Getenv("APP_DB_URL"))
		if err != nil {
			log.Panic(err)
		}

		fmt.Println(pg)

		if err := bb.CreateBaitbot(bot).Serve(); err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
