package yandex

import (
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
	"io"
)

type ReasoningOptions struct {
	Mode string `json:"mode"`
}

type CompletionOptions struct {
	Stream           bool             `json:"stream"`
	Temperature      float64          `json:"temperature"`
	MaxTokens        string           `json:"maxTokens"`
	ReasoningOptions ReasoningOptions `json:"reasoningOptions"`
}

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type GptRequest struct {
	ModelUri          string            `json:"modelUri"`
	CompletionOptions CompletionOptions `json:"completionOptions"`
	Messages          []Message         `json:"messages"`
}

type GptResponse struct {
	Result struct {
		Alternatives []struct {
			Message struct {
				Role string `json:"role"`
				Text string `json:"text"`
			} `json:"message"`
			Status string `json:"status"`
		} `json:"alternatives"`
		Usage struct {
			InputTextTokens  string `json:"inputTextTokens"`
			CompletionTokens string `json:"completionTokens"`
			TotalTokens      string `json:"totalTokens"`
		} `json:"usage"`
		ModelVersion string `json:"modelVersion"`
	} `json:"result"`
}

func (service *Service) AskGPT(prompt, query string) (string, error) {
	gptReq := GptRequest{
		ModelUri: fmt.Sprintf("gpt://%s/yandexgpt", service.FolderId),
		CompletionOptions: CompletionOptions{
			Stream:      false,
			Temperature: 0.6,
			MaxTokens:   "2000",
			ReasoningOptions: ReasoningOptions{
				Mode: "DISABLED",
			},
		},
		Messages: []Message{
			{
				Role: "system",
				Text: prompt,
			},
			{
				Role: "user",
				Text: query,
			},
		},
	}
	gptReqBytes, err := json.Marshal(gptReq)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewReader(gptReqBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", service.IamToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected 200 status code, got %d", resp.StatusCode)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var gptResponse GptResponse
	if err = json.Unmarshal(respBytes, &gptResponse); err != nil {
		return "", err
	}
	var res string
	for _, alternative := range gptResponse.Result.Alternatives {
		res += alternative.Message.Text + " "
	}
	return res, nil
}
