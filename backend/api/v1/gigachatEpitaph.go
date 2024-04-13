package v1

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// gigachatStory handles the story generation based on the type of story requested.
//
// It takes a Fiber context object as a parameter.
// Returns an error.
func gigachatEpitaph(c *fiber.Ctx) error {
	accessToken := c.Locals("sber_access_token").(string)

	baseUrl := viper.GetString("gigachat.baseUrl")

	var requestBody RequestBody

	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	USER_PROMPT = USER_PROMPT_EPITAPH
	USER_PROMPT += "\nС именем: " + requestBody.HumanInfo.Name + "\n" + "Дата рождения: " + requestBody.HumanInfo.DateOfBirth + "\n" + "Дата смерти: " + requestBody.HumanInfo.DateOfDeath + "\n" + "Пол человека: " + requestBody.HumanInfo.Sex
	for key, value := range requestBody.HumanInfo.Questions {
		USER_PROMPT += "\nВопрос: " + key + " Ответ: " + "" + value
	}

	extraContext := []struct {
		Role    string
		Content string
	}{}
	reqBodyBytes, _ := json.Marshal(getGigachatRequestBody(1024, 0.6, extraContext))

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

	text := data.Choices[0].Message.Content
	if len(text) > 300 {
		trimmedText := text[:300]

		lastIndex := strings.LastIndexAny(trimmedText, ".!?…")

		if lastIndex != -1 {
			trimmedText = trimmedText[:lastIndex+1]
		}
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"response": strings.ReplaceAll(trimmedText, "\"", ""),
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"response": strings.ReplaceAll(text, "\"", ""),
		})
	}
}
