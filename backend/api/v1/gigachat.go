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
)

type gigachatRequest struct {
	Model             string  `json:"model"`
	Temperature       float32 `json:"temperature"`
	TopP              float32 `json:"top_p"`
	N                 int     `json:"n"`
	MaxTokens         int     `json:"max_tokens"`
	RepetitionPenalty float32 `json:"repetition_penalty"`
	Stream            bool    `json:"stream"`
	UpdateInterval    int     `json:"update_interval"`
	Messages          []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type gigachatResponse struct {
	Choices []gigachatChoice `json:"choices"`
	Created int32            `json:"created"`
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

// getGigachatRequestBody generates a request body for GigaChat based on the provided maxTokens and temperature.
//
// Parameters:
// - maxTokens: an integer representing the maximum tokens allowed.
// - temperature: a float32 value for controlling the randomness of the generation.
// Returns a pointer to a gigachatRequest struct.
func getGigachatRequestBody(maxTokens int, temperature float32, extraContext []struct {
	Role    string
	Content string
}) *gigachatRequest {
	if maxTokens == 0 || temperature < 0 {
		return nil
	}

	messages := []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		{
			Role:    "system",
			Content: SYSTEM_PROMPT,
		},
		{
			Role:    "user",
			Content: USER_PROMPT,
		},
	}

	for _, ctx := range extraContext {
		messages = append(messages, struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			Role:    ctx.Role,
			Content: ctx.Content,
		})
	}

	return &gigachatRequest{
		Model:             "GigaChat:latest",
		Temperature:       temperature,
		TopP:              0.47,
		N:                 1,
		MaxTokens:         maxTokens,
		RepetitionPenalty: 1.07,
		Stream:            false,
		UpdateInterval:    0,
		Messages:          messages,
	}
}

func gigachatReq(accessToken, baseUrl string, reqBody []byte) (*gigachatResponse, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", baseUrl), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
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
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response gigachatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
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
