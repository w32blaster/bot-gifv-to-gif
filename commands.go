package main

import (
	"fmt"
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
		sendMsg(bot, message.Chat.ID, "Send me an URL please, for example, `https://i.imgur.com/pniMdmr.gifv` or `https://i.imgur.com/GuJXSQ9.mp4`")
		return
	}

	msgProcessing, _ := sendMsg(bot, message.Chat.ID, "Ok, processing...")

	fileName, err := ConvertFile(message.Text)
	defer CleanUp(fileName)

	if err != nil {
		log.Println(err.Error())
		sendMsg(bot, message.Chat.ID, "Error ocurred, sorry :(")
		return
	}

	// now send the file to chat
	data, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s.gif", StorageDirPath, fileName))
	b := tgbotapi.FileBytes{Name: "image.gif", Bytes: data}

	msg := tgbotapi.NewAnimationUpload(message.Chat.ID, b)
	bot.Send(msg)

	msgToDelete := tgbotapi.NewDeleteMessage(message.Chat.ID, msgProcessing.MessageID)
	bot.Send(msgToDelete)

	fileSize := ByteSize(int64(len(data))).String()
	sendMsg(bot, message.Chat.ID, "Done! The result GIF file is "+fileSize)

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

// simple function that prints size of bytes in human readable form
func byteCount(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
