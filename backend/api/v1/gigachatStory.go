package v1

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func gigachatStory(c *fiber.Ctx) error {
	accessToken := c.Locals("sber_access_token").(string)

	baseUrl := viper.GetString("gigachat.baseUrl")

	var requestBody RequestBody

	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if requestBody.TypeOfStory == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type_of_story query param is required",
		})
	}

	if requestBody.TypeOfStory == "epitaph" {
		USER_PROMPT = USER_PROMPT_EPITAPH
		USER_PROMPT += "\nС именем: " + requestBody.HumanInfo.Name + "\n" + "Дата рождения: " + requestBody.HumanInfo.DateOfBirth + "\n" + "Дата смерти: " + requestBody.HumanInfo.DateOfDeath + "\n" + "Пол человека: " + requestBody.HumanInfo.Sex
		for key, value := range requestBody.HumanInfo.Questions {
			USER_PROMPT += "\nВопрос: " + key + " Ответ: " + "" + value
		}
	} else if requestBody.TypeOfStory == "biography" {
		USER_PROMPT = USER_PROMPT_BIOGRAPHY
		USER_PROMPT += "\nИмя: " + requestBody.HumanInfo.Name + "\n" + "Пол: " + requestBody.HumanInfo.Sex + "\n" + "Дата рождения: " + requestBody.HumanInfo.DateOfBirth + "\n" + "Пол человека: " + requestBody.HumanInfo.Sex
		for key, value := range requestBody.HumanInfo.Questions {
			USER_PROMPT += "\nВопрос: " + key + " Ответ: " + "" + value
		}
	}

	reqBody := gigachatRequest{
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
				Content: SYSTEM_PROMPT,
			},
			{
				Role:    "user",
				Content: USER_PROMPT,
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

	var data gigachatResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"response": data.Choices[0].Message.Content,
	})
}
