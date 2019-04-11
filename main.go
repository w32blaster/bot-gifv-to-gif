package main

import (
	"log"
	"flag"
	"net/http"
	"strconv"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {

	// get the command line arguments and parse them
	pBotToken := flag.String("t", "", "Your bot's token. Mandatory")
	pIsDebug := flag.Bool("d", true, "Is debug or not. Default is true")
	pPort := flag.Int("p", 8444, "Port that the bot will run on. Default value is 8444")
	flag.Parse()

	if len(*pBotToken) == 0 {
		panic("The bot token is missing, this is the mandatory parapeter. Please specify it via -t flag. Exit.")
	}

	bot, err := tgbotapi.NewBotAPI(*pBotToken)
	if err != nil {
		panic("Bot doesn't work. Reason: " + err.Error())
	}

	bot.Debug = *pIsDebug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// recommended to make the bot endpoint ending with its token to make it less guessable
	updates := bot.ListenForWebhook("/" + bot.Token)

	// ok, run the bot and listen given port
	go http.ListenAndServe(":"+strconv.Itoa(*pPort), nil)

	for update := range updates {

		if update.Message != nil {
			log.Printf("Message was received: %+v", update.Message)

			if update.Message.IsCommand() {

				// This is a command starting with slash
				ProcessCommands(bot, update.Message)

			} else {

				// simple message
				ProcessSimpleMessage(bot, update.Message)
			}

		} else if update.CallbackQuery != nil {

			// this is the callback after a button click
			// TODO

		} else if update.InlineQuery != nil {

			// this is inline query
			// TODO
		}

	}

}
