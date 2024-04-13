package v1

const SYSTEM_PROMPT = "Вы бот, созданный для помощи пользователям в заполнении страницы памяти. Ваша задача - генерировать тексты для сложных полей формы на основе предоставленных данных. Пользователи будут предоставлять вам информацию по одному запросу за раз, и вы должны будете создавать соответствующий текст. При генерации текста для сложных полей, таких как эпитафия или биография, предоставьте возможность пользователю внести коррективы в предложенный текст перед его сохранением. Важно обеспечить плавный и естественный процесс взаимодействия с пользователем в формате чат-интерфейса. Твое имя - MemoryCode бот"
const SYSTEM_PROMPT_QUESTIONS_FORMAT = "Используй строго следующий формат:\n1. ВОПРОС\n2. ВОПРОС\n3. ВОПРОС\n4. ВОПРОС\n5. ВОПРОС\n..."

var USER_PROMPT = ""
var USER_PROMPT_EPITAPH = "Сгенируй обязательно краткую (не более 200 символов) питафию для следующего человека опираясь на факты о нем"
var USER_PROMPT_BIOGRAPHY = "Сгенируй обязательно краткую биографию (не более 200 символов) для следующего человека опираясь на факты о нем"
var USER_PROMPT_QUESTIONS = "Сделай короткие вопросы (от 6 до 10) для создания страницы памяти для"

type RequestBody struct {
	RequestMessage string    `json:"request_message"`
	TypeOfStory    string    `json:"type_of_story"`
	HumanInfo      HumanInfo `json:"human_info"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type HumanInfo struct {
	Name        string            `json:"name"`
	Sex         string            `json:"sex"`
	DateOfBirth string            `json:"birth_date"`
	DateOfDeath string            `json:"death_date"`
	Questions   map[string]string `json:"questions"`
}
