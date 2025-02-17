package yandex

import (
	"encoding/base64"
	"net/http"
	"encoding/json"
	"bytes"
	"io"
	"fmt"
)

type LanguageCode = string

const (
	LanguageCodeEnglish LanguageCode = "en"
	LanguageCodeRussian LanguageCode = "ru"
	LanguageCodeAuto    LanguageCode = "*"
)

type MimeType = string

const (
	MimeTypeJPEG  = "JPEG"
	MimeTypePNG   = "PNG"
	MimeTypeEmpty = ""
)

type ocrRequest struct {
	MimeType      MimeType       `json:"mimeType"`
	LanguageCodes []LanguageCode `json:"languageCodes"`
	Model         string         `json:"model"`
	Content       string         `json:"content"`
}

type ocrResponse struct {
	Result struct {
		TextAnnotation struct {
			Width  string `json:"width"`
			Height string `json:"height"`
			Blocks []struct {
				BoundingBox struct {
					Vertices []struct {
						X string `json:"x"`
						Y string `json:"y"`
					} `json:"vertices"`
				} `json:"boundingBox"`
				Lines []struct {
					BoundingBox struct {
						Vertices []struct {
							X string `json:"x"`
							Y string `json:"y"`
						} `json:"vertices"`
					} `json:"boundingBox"`
					Text  string `json:"text"`
					Words []struct {
						BoundingBox struct {
							Vertices []struct {
								X string `json:"x"`
								Y string `json:"y"`
							} `json:"vertices"`
						} `json:"boundingBox"`
						Text         string `json:"text"`
						EntityIndex  string `json:"entityIndex"`
						TextSegments []struct {
							StartIndex string `json:"startIndex"`
							Length     string `json:"length"`
						} `json:"textSegments"`
					} `json:"words"`
					TextSegments []struct {
						StartIndex string `json:"startIndex"`
						Length     string `json:"length"`
					} `json:"textSegments"`
					Orientation string `json:"orientation"`
				} `json:"lines"`
				Languages []struct {
					LanguageCode string `json:"languageCode"`
				} `json:"languages"`
				TextSegments []struct {
					StartIndex string `json:"startIndex"`
					Length     string `json:"length"`
				} `json:"textSegments"`
			} `json:"blocks"`
			Entities []interface{} `json:"entities"`
			Tables   []interface{} `json:"tables"`
			FullText string        `json:"fullText"`
			Rotate   string        `json:"rotate"`
		} `json:"textAnnotation"`
		Page string `json:"page"`
	} `json:"result"`
}

func (service *Service) ImageToText(image []byte, mimeType MimeType, langCodes []LanguageCode) (string, error) {
	ocrReq := ocrRequest{
		MimeType:      mimeType,
		LanguageCodes: langCodes,
		Model:         "page",
		Content:       base64.StdEncoding.EncodeToString(image),
	}
	ocrReqBytes, err := json.Marshal(ocrReq)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, "https://ocr.api.cloud.yandex.net/ocr/v1/recognizeText", bytes.NewReader(ocrReqBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+service.IamToken)
	req.Header.Set("x-folder-id", service.FolderId)
	req.Header.Set("x-data-logging-enabled", "true")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed with status code %d", resp.StatusCode)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ocrResp ocrResponse
	if err = json.Unmarshal(respBytes, &ocrResp); err != nil {
		return "", err
	}
	var res string
	for _, block := range ocrResp.Result.TextAnnotation.Blocks {
		for _, line := range block.Lines {
			res += line.Text + " "
		}
	}
	return res, nil
}
