package v1

import "fmt"

type YandexRequest struct {
	ModelURI          string `json:"modelUri"`
	CompletionOptions struct {
		Stream      bool    `json:"stream"`
		Temperature float32 `json:"temperature"`
		MaxTokens   string  `json:"maxTokens"`
	} `json:"completionOptions"`
	Messages []struct {
		Role string `json:"role"`
		Text string `json:"text"`
	} `json:"messages"`
}

// getYandexRequestBody generates a Yandex request body based on the provided folder ID, max tokens, and temperature.
//
// folderId: The ID of the folder.
// maxTokens: The maximum number of tokens.
// temperature: The temperature for the completion.
// Returns a pointer to a YandexRequest.
func getYandexRequestBody(folderId string, maxTokens string, temperature float32, extraContext []struct {
	Role string
	Text string
}) *YandexRequest {
	if folderId == "" || maxTokens == "" || temperature < 0 {
		return nil
	}

	messages := []struct {
		Role string `json:"role"`
		Text string `json:"text"`
	}{
		{
			Role: "system",
			Text: SYSTEM_PROMPT,
		},
		{
			Role: "user",
			Text: USER_PROMPT,
		},
	}

	for _, ctx := range extraContext {
		messages = append(messages, struct {
			Role string `json:"role"`
			Text string `json:"text"`
		}{
			Role: ctx.Role,
			Text: ctx.Text,
		})
	}

	return &YandexRequest{
		ModelURI: fmt.Sprintf("gpt://%s/yandexgpt-lite", folderId),
		CompletionOptions: struct {
			Stream      bool    `json:"stream"`
			Temperature float32 `json:"temperature"`
			MaxTokens   string  `json:"maxTokens"`
		}{
			Stream:      false,
			Temperature: temperature,
			MaxTokens:   maxTokens,
		},
		Messages: messages,
	}
}

type yandexResult struct {
	Alternatives []yandexAlternative `json:"alternatives"`
}

type yandexAlternative struct {
	Message yandexMessage `json:"message"`
	Status  string        `json:"status"`
}

type yandexMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type yandexResponseData struct {
	Result yandexResult `json:"result"`
}
