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
func gigachatGenerateQuestions(c *fiber.Ctx) error {
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

	USER_PROMPT = USER_PROMPT_QUESTIONS + " человека " + requestBody.HumanInfo.Sex + " пола" + " по имени " + requestBody.HumanInfo.Name + "," + " дата рождения " + requestBody.HumanInfo.DateOfBirth + " " + SYSTEM_PROMPT_QUESTIONS_FORMAT + "\n " + "Примеры вопросов: Что Валерий ценил больше всего в друзьях?\nКакие увлечения были у Валерия?\nКакие черты характера делали Валерия особенным?\nО чём Валерий любил рассказывать?\nКакой совет Валерий дал бы молодым?\nКакие слова или фразы чаще всего говорил Илья?\nКакие моменты из детства Ильи Вам известны?\nЧем Илья любил заниматься в свободное время?\nКакой смешной случай из жизни Ильи Вы можете вспомнить?\nКакие качества Ильи восхищали его друзей и близких?\nЧто Илья считал своим главным достижением?\nКакие истории из молодости Митрофана наиболее запомнились?\nКак Митрофан относился к изменениям в жизни?\nКакие увлечения Митрофана могли удивить его знакомых?\nЧто для Митрофана значила семья?\nКакими мудрыми словами Митрофан любил делиться?\nЧто в поведении Митрофана вдохновляло окружающих?\nЧто из детства Бориса Николаевича оказало на него наибольшее влияние?\nКакие личные качества помогли Борису Николаевичу в жизни?\nКакие моменты в карьере Бориса Николаевича были особенно значимыми?\nКакие увлечения у Бориса Николаевича были вне работы?\nКак Борис Николаевич относился к принятию решений?\nЧто Борис Николаевич считал своим главным достижением?\nКакие черты Надежды делали её незабываемой для друзей?\nКакие моменты из жизни Надежды вызывают улыбку?\nКак Надежда любила проводить время с семьёй?\nЧему Надежда могла научить любого человека?\nЧто Надежда считала своей гордостью?\nКакие события в жизни Надежды были особенно важны?\nКакой совет Юрий Александрович давал чаще всего?\nКакие важные события в жизни Юрия оставили отпечаток на его характере?\nЧто Юрий ценил в людях больше всего?\nКакие увлечения были у Юрия?\nКак Юрий относился к своим ошибкам и неудачам?\nКакие традиции Юрий считал важными в своей семье?\nКакие хобби у Елены были наиболее необычными?\nЧто Елена могла рассказать о своих путешествиях?\nКак Елена выражала свою любовь к близким?\nКакие события Елена считала поворотными в своей жизни?\nЧем Елена гордилась в своей жизни?\nКакой была мечта Елены, которую она так и не смогла осуществить?\nЧему Алексей посвящал большую часть своего времени?\nКакие слова или выражения были фирменными для Алексея?\nКакой историей из своей жизни Алексей любил делиться?\nКакие черты Алексея вызывали восхищение у его друзей?\nЧто для Алексея значила семья?\nКакой совет Алексей считал самым важным для передачи молодому поколению?\nКак Амет-Хан относился к своим достижениям?\nКакие личные истории Амет-Хана наиболее вдохновляющие?\nЧему Амет-Хан учил своих детей и внуков?\nКакие увлечения Амет-Хана были наиболее неожиданными для его окружения?\nЧто Амет-Хан считал своим главным жизненным уроком?\nКак Амет-Хан предпочитал проводить своё свободное время?\nКакие воспоминания о Елене Михайловне наиболее дороги её друзьям?\nЧто Елена Михайловна считала своими главными достоинствами?\nКакие книги Елена Михайловна любила больше всего?\nЧему Елена Михайловна стремилась научить своих детей?\nКакие события в жизни Елены Михайловны были особенно значимыми?\nКакие советы Елена Михайловна часто повторяла?\nКакие увлечения Игоря Владимировича были наиболее увлекательными?\nКак Игорь Владимирович относился к новым вызовам и возможностям?\nКакие личные качества Игоря Владимировича наиболее восхищали его окружение?\nКакие события в жизни Игоря оставили наибольший след в его характере?\nЧто Игорь Владимирович считал своим главным достижением?\nКак Игорь Владимирович предпочитал отдыхать?\nКакие стихотворения Сергея Васильевича были особенно значимы для него самого?\nКакие черты характера Сергея Васильевича помогли ему в его карьере?\nКакие важные моменты из жизни Сергея Васильевича он сам любил вспоминать?\nКак Сергей Васильевич относился к своим неудачам и успехам?\nЧем Сергей Васильевич любил заниматься в свободное время?\nКакие советы Сергей Васильевич давал молодым писателям?\nКакие хобби у Ксении были наиболее необычными?\nЧто Ксения могла рассказать о своих путешествиях?\nКак Ксения выражала свою любовь к близким?\nКакие события Ксения считала поворотными в своей жизни?\nЧем Ксения гордилась в своей жизни?\nКакой была мечта Ксении, которую она так и не смогла осуществить?\nЧему Лидия Михайловна посвящала большую часть своего времени?\nКакие слова или выражения были фирменными для Лидии Михайловны?\nКакой историей из своей жизни Лидия Михайловна любила делиться?\nКакие черты Лидии вызывали восхищение у её друзей?\nЧто для Лидии значила семья?\nКакой совет Лидия Михайловна считала самым важным для передачи молодому поколению?"

	extraConContent := []struct {
		Role    string
		Content string
	}{}

	reqBodyBytes, _ := json.Marshal(getGigachatRequestBody(1024, 0.6, extraConContent))

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
			cleanLine := numberDotRegex.ReplaceAllString(line, "") // Remove the number and dot at the beginning
			questions = append(questions, cleanLine)
		}
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"response": questions,
	})
}
