package v1

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
