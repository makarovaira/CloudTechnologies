package telegram

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io"
)

var fileUrl = "https://api.telegram.org/file/bot" + botToken

type Photo struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`

	Width    int `json:"width"`
	Height   int `json:"height"`
	FileSize int `json:"file_size"`
}

func (photo *Photo) FilePath() (string, error) {
	url := apiUrl + "/getFile?file_id=" + photo.FileID

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("recieved reponse status: %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result struct {
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}
	if err = json.Unmarshal(b, &result); err != nil {
		return "", err
	}
	return result.Result.FilePath, nil
}

func (photo *Photo) Download() ([]byte, error) {
	filePath, err := photo.FilePath()
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(fileUrl + "/" + filePath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recieved reponse status: %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
