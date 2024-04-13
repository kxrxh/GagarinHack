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

	USER_PROMPT = USER_PROMPT_QUESTIONS + " человека " + requestBody.HumanInfo.Sex + " пола " + "  по имени " + requestBody.HumanInfo.Name + "," + " дата рождения " + requestBody.HumanInfo.DateOfBirth + SYSTEM_PROMPT_QUESTIONS_FORMAT + "\n" + "Примеры вопросов:Что [имя] ценил больше всего в друзьях?\nКакие увлечения были у [имя]?\nКакие черты характера делали [имя] особенным?\nО чём [имя] любил рассказывать?\nКакой совет [имя] дал бы молодым?\nКакие слова или фразы чаще всего говорил [имя]?\nКакие моменты из детства [имя] Вам известны?\nЧем [имя] любил заниматься в свободное время?\nКакой смешной случай из жизни [имя] Вы можете вспомнить?\nКакие качества [имя] восхищали его друзей и близких?\nЧто [имя] считал своим главным достижением?\nКакие истории из молодости [имя] наиболее запомнились?\nКак [имя] относился к изменениям в жизни?\nКакие увлечения [имя] могли удивить его знакомых?\nЧто для [имя] значила семья?\nКакими мудрыми словами [имя] любил делиться?\nЧто в поведении [имя] вдохновляло окружающих?\nЧто из детства [имя] оказало на него наибольшее влияние?\nКакие личные качества помогли [имя] в жизни?\nКакие моменты в карьере [имя] были особенно значимыми?\nКакие увлечения у [имя] были вне работы?\nКак [имя] относился к принятию решений?\nЧто [имя] считал своим главным достижением?\nКакие черты [имя] делали его незабываемым для друзей?\nКакие моменты из жизни [имя] вызывают улыбку?\nКак [имя] любил проводить время с семьёй?\nЧему [имя] могла научить любого человека?\nЧто [имя] считала своей гордостью?\nКакие события в жизни [имя] были особенно важны?\nКакой совет [имя] давал чаще всего?\nКакие важные события в жизни [имя] оставили отпечаток на его характере?\nЧто [имя] ценил в людях больше всего?\nКакие увлечения были у [имя]?\nКак [имя] относился к своим ошибкам и неудачам?\nКакие традиции [имя] считал важными в своей семье?\nКакие хобби у [имя] были наиболее необычными?\nЧто [имя] могла рассказать о своих путешествиях?\nКак [имя] выражала свою любовь к близким?\nКакие события [имя] считала поворотными в своей жизни?\nЧем [имя] гордилась в своей жизни?\nКакой была мечта [имя], которую она так и не смогла осуществить?\nЧему [имя] посвящал большую часть своего времени?\nКакие слова или выражения были фирменными для [имя]?\nКакой историей из своей жизни [имя] любил делиться?\nКакие черты [имя] вызывали восхищение у его друзей?\nЧто для [имя] значила семья?\nКакой совет [имя] считал самым важным для передачи молодому поколению?\nКак [имя] относился к своим достижениям?"

	extraContext := []struct {
		Role string
		Text string
	}{}

	reqBodyBytes, _ := json.Marshal(getYandexRequestBody(folderId, "1024", 0.2, extraContext))

	req, _ := http.NewRequest("POST", "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewBuffer(reqBodyBytes))

	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", apiKey))
	req.Header.Set("Text-Type", "application/json")
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
