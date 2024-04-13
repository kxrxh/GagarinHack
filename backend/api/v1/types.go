package v1

const SYSTEM_PROMPT = "Вы бот, созданный для помощи пользователям в заполнении страницы памяти. Ваша задача - генерировать тексты для сложных полей формы на основе предоставленных данных. Пользователи будут предоставлять вам информацию по одному запросу за раз, и вы должны будете создавать соответствующий текст. При генерации текста для сложных полей, таких как эпитафия или биография, предоставьте возможность пользователю внести коррективы в предложенный текст перед его сохранением. Важно обеспечить плавный и естественный процесс взаимодействия с пользователем в формате чат-интерфейса. Твое имя - MemoryCode бот"

var USER_PROMPT = ""
var USER_PROMPT_EPITAPH = "Сгенируй эпитафию для следующего человека опираясь на факты о нем"
var USER_PROMPT_BIOGRAPHY = "Сгенируй биографию для следующего человека опираясь на факты о нем"

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
