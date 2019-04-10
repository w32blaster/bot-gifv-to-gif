package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var urlRegex = regexp.MustCompile(`^http(s)?://i\.imgur\.com\/[a-zA-Z0-9]+\.(mp4|gifv)$`)

// ProcessCommands process commands
func ProcessCommands(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	command := extractCommand(message.Command())

	switch command {

	case "help":
		sendMsg(bot, chatID, "Just send me a URL link to a GIFV file and that's it")

	case "start":
		sendMsg(bot, chatID, "Hi there! I can convert animated GIFV file (usually from Imgur website) to usual GIF file that you can add to your telegram and use in your chatting. Enjoy!")

	case "about":
		sendMsg(bot, chatID, "Github page: https://github.com/w32blaster/bot-gifv-to-gif")

	case "list":
		sendMsg(bot, chatID, "Here are list of all the saved items to watch")
	}
}

// ProcessSimpleMessage is called when simeone sends a plain text, expected to be URL, but we need to check
func ProcessSimpleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if !ifGifvURL(message.Text) {
		sendMsg(bot, message.Chat.ID, "Send me an URL please, for example, `https://i.imgur.com/pniMdmr.gifv`")
		return
	}

	ConvertFile(message.Text)

	data, _ := ioutil.ReadFile("/tmp/temp_file_to_convert2.gif")
	b := tgbotapi.FileBytes{Name: "image.gif", Bytes: data}

	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, b)
	msg.Caption = "Test"
	bot.Send(msg)
}

// test that the given URL is valid GIFV file
func ifGifvURL(url string) bool {
	return urlRegex.MatchString(url)
}

// properly extracts command from the input string, removing all unnecessary parts
// please refer to unit tests for details
func extractCommand(rawCommand string) string {

	command := rawCommand

	// remove slash if necessary
	if rawCommand[0] == '/' {
		command = command[1:]
	}

	// if command contains the name of our bot, remote it
	command = strings.Split(command, "@")[0]
	command = strings.Split(command, " ")[0]

	return command
}

// simply send a message to bot in Markdown format
func sendMsg(bot *tgbotapi.BotAPI, chatID int64, textMarkdown string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, textMarkdown)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true

	// send the message
	resp, err := bot.Send(msg)
	if err != nil {
		log.Println("bot.Send:", err, resp)
		return resp, err
	}

	return resp, err
}
