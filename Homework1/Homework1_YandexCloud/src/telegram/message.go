package telegram

import (
	"os"
	"encoding/json"
	"net/http"
	"fmt"
	"bytes"
)

var (
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	apiUrl   = "https://api.telegram.org/bot" + botToken
)

type Chat struct {
	Id int `json:"id"`
}

type Message struct {
	Id     int     `json:"message_id"`
	Chat   Chat    `json:"chat"`
	Text   *string `json:"text"`
	Photos []Photo `json:"photo"`
}

type replyParameters struct {
	MessageId int `json:"message_id"`
}

type sendMessageRequest struct {
	ChatId          int             `json:"chat_id"`
	Text            string          `json:"text"`
	ReplyParameters replyParameters `json:"reply_parameters"`
}

func (message *Message) Reply(text string) error {
	request := sendMessageRequest{
		ChatId:          message.Chat.Id,
		Text:            text,
		ReplyParameters: replyParameters{MessageId: message.Id},
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	resp, err := http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned status code %d", resp.StatusCode)
	}
	return nil
}
