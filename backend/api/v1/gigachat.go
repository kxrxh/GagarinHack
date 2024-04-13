package v1

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type gigachatRequest struct {
	Model             string  `json:"model"`
	Temperature       float64 `json:"temperature"`
	TopP              float64 `json:"top_p"`
	N                 int     `json:"n"`
	MaxTokens         int     `json:"max_tokens"`
	RepetitionPenalty float64 `json:"repetition_penalty"`
	Stream            bool    `json:"stream"`
	UpdateInterval    int     `json:"update_interval"`
	Messages          []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type gigachatResponse struct {
	Choices []gigachatChoice `json:"choices"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Object  string           `json:"object"`
	Usage   gigachatUsage    `json:"usage"`
}

type gigachatChoice struct {
	Message      gigachatMessage `json:"message"`
	Index        int             `json:"index"`
	FinishReason string          `json:"finish_reason"`
}

type gigachatMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type gigachatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
	SystemTokens     int `json:"system_tokens"`
}

func getAccessToken(c *fiber.Ctx) error {
	scope := viper.GetString("gigachat.scope")
	authUrl := viper.GetString("gigachat.authUrl")
	credentials := viper.GetString("gigachat.credentials")

	form := url.Values{}
	form.Add("scope", scope)

	req, _ := http.NewRequest("POST", authUrl, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("RqUID", "6f0b1291-c7f3-43c6-bb2e-9f3efb2dc98e")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var accessTokenResponse AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&accessTokenResponse); err != nil {
		return err
	}

	c.Locals("sber_access_token", accessTokenResponse.AccessToken)
	return c.Next()
}
