package main

import (
	"os"
	"cheatsheet_bot/yandex"
	"net/http"
	"cheatsheet_bot/telegram"
	"log"
	"context"
	"encoding/json"
	"fmt"
)

const (
	helpMessage         = "Я помогу подготовить ответ на экзаменационный вопрос по дисциплине \"Операционные системы\".\nПришлите мне фотографию с вопросом или наберите его текстом."
	gptErrorMessage     = "Я не смог подготовить ответ на экзаменационный вопрос."
	invalidPhotoMessage = "Я не могу обработать эту фотографию."
	notSupportedMessage = "Я могу обработать только текстовое сообщение или фотографию."
)

var (
	folderId      = os.Getenv("FOLDER_ID")
	mountPoint    = os.Getenv("MOUNT_POINT")
	bucketObjKey  = os.Getenv("BUCKET_OBJECT_KEY")
	defaultPrompt = readPrompt()
)

type Request struct {
	Method string `json:"httpMethod"`
	Body   string `json:"body"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func Handler(ctx context.Context, requestBytes []byte) (*Response, error) {
	var request Request
	if err := json.Unmarshal(requestBytes, &request); err != nil {
		log.Printf("failed ot unmarshal request: %v\n", err)
		return nil, err
	}
	var messageRequest struct {
		Message telegram.Message `json:"message"`
	}
	if err := json.Unmarshal([]byte(request.Body), &messageRequest); err != nil {
		log.Printf("failed to unmarshal message request: %v\n", err)
		return nil, err
	}
	iam, err := yandex.GetServiceIamToken(ctx)
	if err != nil {
		log.Printf("failed to get IAM token: %v\n", err)
		return nil, err
	}
	if err = handleMessage(iam.Token, &messageRequest.Message); err != nil {
		log.Printf("failed to handle message: %v\n", err)
		return nil, err
	}
	return &Response{StatusCode: http.StatusOK}, nil
}

func handleMessage(iamToken string, message *telegram.Message) error {
	if messageText := message.Text; messageText != nil {
		return handleText(*messageText, iamToken, message)
	}
	if len(message.Photos) != 0 {
		return handlePhoto(iamToken, message)
	}
	return message.Reply(notSupportedMessage)
}

func handleText(text, iamToken string, message *telegram.Message) error {
	switch text {
	case "/start", "/help":
		return message.Reply(helpMessage)
	default:
		result, err := yandex.NewService(folderId, iamToken).AskGPT(defaultPrompt, text)
		if err != nil {
			log.Printf("failed to ask gpt for %s: %v\n", text, err)
			return message.Reply(gptErrorMessage)
		}
		return message.Reply(result)
	}
}

func handlePhoto(iamToken string, message *telegram.Message) error {
	photo := message.Photos[len(message.Photos)-1]
	img, err := photo.Download()
	if err != nil {
		log.Printf("failed to download photo: %v\n", err)
		return message.Reply(invalidPhotoMessage)
	}
	mime, err := detectImageType(img)
	if err != nil {
		log.Printf("failed to detect image type: %v\n", err)
		return message.Reply(invalidPhotoMessage)
	}
	service := yandex.NewService(folderId, iamToken)
	text, err := service.ImageToText(img, mime, []yandex.LanguageCode{yandex.LanguageCodeAuto})
	if err != nil {
		log.Printf("failed to convert image to text: %v\n", err)
		return message.Reply(invalidPhotoMessage)
	}
	result, err := service.AskGPT(defaultPrompt, text)
	if err != nil {
		log.Printf("failed to ask gpt: %v\n", err)
		return message.Reply(gptErrorMessage)
	}
	return message.Reply(result)
}

func detectImageType(buffer []byte) (yandex.MimeType, error) {
	if buffer[0] == 0xFF && buffer[1] == 0xD8 && buffer[2] == 0xFF {
		return yandex.MimeTypeJPEG, nil
	} else if buffer[0] == 0x89 && buffer[1] == 0x50 && buffer[2] == 0x4E && buffer[3] == 0x47 {
		return yandex.MimeTypePNG, nil
	}
	return yandex.MimeTypeEmpty, fmt.Errorf("unsupported file type")
}

func readPrompt() string {
	data, _ := os.ReadFile("/function/storage/mnt/prompt.txt")
	return string(data)
}
