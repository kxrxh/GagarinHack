package v1

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// gigachatGenerateQuestions generates questions for a chat based on the provided request body parameters.
//
// Parameters:
//   - c: ConContent object for handling HTTP requests and responses.
//
// Return type: error
// gigachatGenerateQuestions generates questions using the Gigachat API.
func gigachatGenerateBioQuestions(c *fiber.Ctx) error {
	accessToken := c.Locals("sber_access_token").(string)

	baseUrl := viper.GetString("gigachat.baseUrl")

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

	if requestBody.TypeOfStory == "youth" {
		USER_PROMPT = USER_PROMPT_QUESTIONS + " биографии, которая детально раскрывает детство и юношество для его биографии" + requestBody.HumanInfo.Sex + " пола по имени " + requestBody.HumanInfo.Name + "," + " дата рождения " + requestBody.HumanInfo.DateOfBirth + " " + SYSTEM_PROMPT_QUESTIONS_FORMAT + "\n " + "Вот пример вопросов: Какие ранние воспоминания из детства [имя] оказали на него наибольшее влияние?\nВ каких играх и занятиях [имя] предпочитал участвовать в детстве?\nКакие особенные события или традиции в семье [имя] он особенно запомнил?\nКакой предмет в школе был любимым у [имя], и почему?\nКакие книги или истории вдохновляли [имя] в юном возрасте?\nБыли ли у [имя] какие-то необычные увлечения или таланты в детстве?\nКакие вызовы или трудности пришлось преодолеть [имя] в ранние годы?\nКто был самым значимым человеком для [имя] в его юности и почему?\nКакой значимый совет или урок [имя] получил от своих родителей или наставников?\nКакие мечты и амбиции были у [имя] в юношеском возрасте?"
	} else if requestBody.TypeOfStory == "middle_age" {
		USER_PROMPT = USER_PROMPT_QUESTIONS + " биографии, которая детально раскрывает средний возраст для его биографии " + requestBody.HumanInfo.Sex + " пола по имени " + requestBody.HumanInfo.Name + "," + " дата рождения " + requestBody.HumanInfo.DateOfBirth + " " + SYSTEM_PROMPT_QUESTIONS_FORMAT + "\n " + "Вот пример вопросов: Какие ключевые карьерные достижения [имя] произошли в среднем возрасте?\nКакие жизненные уроки [имя] считает наиболее значимыми за это время?\nКакие события в личной жизни [имя] оказали наибольшее влияние на его развитие и взгляды?\nКак [имя] справлялся с жизненными испытаниями и переменами в этот период?\nКакие хобби или увлечения приобрел [имя] в среднем возрасте и почему?\nКак изменились отношения [имя] с близкими и друзьями в этот период его жизни?\nЧем [имя] гордится больше всего при взгляде назад на эти годы?\nКакие места или путешествия оказали особое влияние на [имя] в среднем возрасте?\nКакие профессиональные или личные вызовы [имя] успешно преодолел в этот период?\nКак [имя] видит свое влияние на общество и окружающих в эти годы?"
	} else if requestBody.TypeOfStory == "old_age" {
		USER_PROMPT = USER_PROMPT_QUESTIONS + " биографии, которая детально раскрывает последние года жизни для его биографии " + requestBody.HumanInfo.Sex + " пола по имени " + requestBody.HumanInfo.Name + "," + " дата рождения " + requestBody.HumanInfo.DateOfBirth + " " + SYSTEM_PROMPT_QUESTIONS_FORMAT + "\n " + "Как [имя] описал бы свои ощущения и переживания в последние годы жизни?\nКакие моменты или события оставили наиболее значимый след в последние годы его жизни?\nКакие увлечения или хобби [имя] нашёл или продолжал заниматься в последние годы?\nКак [имя] относился к старению и изменениям в своей жизни?\nЕсть ли какие-то особенные моменты или истории, которыми [имя] любил делиться о своем прошлом в последние годы?\nКакие отношения и связи были особенно важны для [имя] в его последние годы?\nКакие надежды и мечты [имя] выразил бы, глядя назад на свою жизнь?\nКакие уроки [имя] хотел бы передать младшим поколениям на основе своего опыта?\nКакие изменения в мире оказали на [имя] наибольшее влияние в последние годы?\nКак [имя] видел своё наследие и что он хотел бы оставить после себя?"
	}

	extraConContent := []struct {
		Role    string
		Content string
	}{}

	reqBodyBytes, _ := json.Marshal(getGigachatRequestBody(1024, 0.8, extraConContent))

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
	questionRegex := regexp.MustCompile(`^\d*\.?\s*([А-Я].*\?)$`)
	numberDotRegex := regexp.MustCompile(`^\d+\.\s*`)

	var questions []string
	for _, line := range strings.Split(data.Choices[0].Message.Content, "\n") {
		if questionRegex.MatchString(line) {
			cleanLine := strings.ReplaceAll(numberDotRegex.ReplaceAllString(line, ""), "\"", "") // Remove the number and dot at the beginning
			questions = append(questions, cleanLine)
		}
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"response": questions,
	})
}
