package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// yandexQuestions handles the processing of a Yandex questions request.
//
// c: Context provided by Fiber framework.
// error: An error interface.
func yandexQuestions(c *fiber.Ctx) error {
	apiKey := viper.GetString("yandex.api_key")
	folderId := viper.GetString("yandex.folder_id")

	var requestBody RequestBody

	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if requestBody.HumanInfo.Sex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sex query param is required",
		})
	}

	if requestBody.HumanInfo.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name query param is required",
		})
	}

	if requestBody.HumanInfo.DateOfBirth == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "date_of_birth query param is required",
		})
	}

	USER_PROMPT = USER_PROMPT_QUESTIONS + " человека " + requestBody.HumanInfo.Sex + " пола " + "  по имени " + requestBody.HumanInfo.Name + "," + " дата рождения " + requestBody.HumanInfo.DateOfBirth + SYSTEM_PROMPT_QUESTIONS_FORMAT

	reqBodyBytes, _ := json.Marshal(getYandexRequestBody(folderId, "1024", 0.1))

	req, _ := http.NewRequest("POST", "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewBuffer(reqBodyBytes))

	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "PostmanRuntime/7.37.0")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	var data yandexResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	questionRegex := regexp.MustCompile(`^\d*\.?\s*([А-Я].*\?)$`)
	numberDotRegex := regexp.MustCompile(`^\d+\.\s*`)

	var questions []string
	for _, line := range strings.Split(data.Result.Alternatives[0].Message.Text, "\n") {
		if questionRegex.MatchString(line) {
			cleanLine := numberDotRegex.ReplaceAllString(line, "") // Remove the number and dot at the beginning
			questions = append(questions, cleanLine)
		}
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"response": questions,
	})

}
