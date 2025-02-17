package yandex

import (
	"time"
	"os"
	"net/http"
	"fmt"
	"io"
	"context"
	"encoding/json"
)

type IamToken struct {
	Token     string    `json:"iamToken"`
	ExpiresAt time.Time `json:"expiresAt"`
}

const (
	InstanceMetadataOverrideEnvVar = "YC_METADATA_ADDR"
	InstanceMetadataAddr           = "169.254.169.254"
)

func getMetadataServiceAddr() string {
	if nonDefaultAddr := os.Getenv(InstanceMetadataOverrideEnvVar); nonDefaultAddr != "" {
		return nonDefaultAddr
	}
	return InstanceMetadataAddr
}

func GetServiceIamToken(ctx context.Context) (*IamToken, error) {
	url := fmt.Sprintf("http://%s/computeMetadata/v1/instance/service-accounts/default/token", getMetadataServiceAddr())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("%s.\n"+
			"Are you inside compute instance?",
			err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s.\n"+
			"Is this compute instance running using Service Account? That is, Instance.service_account_id should not be empty.",
			resp.Status)
	}
	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if err != nil {
			body = []byte(fmt.Sprintf("Failed response body read failed: %s", err.Error()))
		}
		return nil, fmt.Errorf("%s", resp.Status)
	}
	if err != nil {
		return nil, fmt.Errorf("reponse read failed: %s", err)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, err
	}
	return &IamToken{
		Token:     tokenResponse.AccessToken,
		ExpiresAt: time.UnixMilli(tokenResponse.ExpiresIn),
	}, nil
}
