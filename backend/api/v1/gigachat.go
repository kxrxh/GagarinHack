package v1

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type completionRequest struct {
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

func gigachatCompletion(c *fiber.Ctx) error {
	accessToken := c.Locals("sber_access_token").(string)

	baseUrl := viper.GetString("gigachat.baseUrl")

	var requestBody RequestBody

	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if requestBody.RequestMessage == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "request_message query param is required",
		})
	}

	reqBody := completionRequest{
		Model:             "GigaChat:latest",
		Temperature:       0.6,
		TopP:              0.47,
		N:                 1,
		MaxTokens:         1024,
		RepetitionPenalty: 1.07,
		Stream:            false,
		UpdateInterval:    0,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "system",
				Content: SYSTEM_PROMT,
			},
			{
				Role:    "user",
				Content: requestBody.RequestMessage,
			},
		},
	}

	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", baseUrl), bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("X-Request-ID", "79e41a5f-f180-4c7a-b2d9-393086ae20a1")
	req.Header.Set("X-Session-ID", "b6874da0-bf06-410b-a150-fd5f9164a0b2")
	req.Header.Set("X-Client-ID", "b6874da0-bf06-410b-a150-fd5f9164a0b2")

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

	body, _ := io.ReadAll(resp.Body)
	return c.SendString(string(body))
}
