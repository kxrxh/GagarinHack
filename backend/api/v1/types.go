package v1

const SYSTEM_PROMPT = "Вы бот, созданный для помощи пользователям в заполнении страницы памяти. Ваша задача - генерировать тексты для сложных полей формы на основе предоставленных данных. Пользователи будут предоставлять вам информацию по одному запросу за раз, и вы должны будете создавать соответствующий текст. При генерации текста для сложных полей, таких как эпитафия или биография, предоставьте возможность пользователю внести коррективы в предложенный текст перед его сохранением. Важно обеспечить плавный и естественный процесс взаимодействия с пользователем в формате чат-интерфейса. Твое имя - MemoryCode бот"
const SYSTEM_PROMPT_QUESTIONS_FORMAT = "Используй строго следующий формат:\n1. ВОПРОС\n2. ВОПРОС\n3. ВОПРОС\n4. ВОПРОС\n5. ВОПРОС\n..."
const SYSTEM_PROMPT_HEADING_FORMAT = "Используй строго следующий формат:\n1. ЗАГОЛОВОК\n2. ЗАГОЛОВОК\n3. ЗАГОЛОВОК\n4. ЗАГОЛОВОК\n5. ЗАГОЛОВОК\n..."

var USER_PROMPT = ""
var USER_PROMPT_EPITAPH = "Сгенируй обязательно краткую (не более 200 символов) эпитафию для следующего человека опираясь на факты о нем"
var USER_PROMPT_BIOGRAPHY_START = "На данный момент в биографии человека имеется следующее: "
var USER_PROMPT_BIOGRAPHY = "Сгенируй продолжение в виде детальной биографии для описания"
var USER_PROMPT_ENDING = "Сгенируй детальную заключение для биографии "
var USER_PROMPT_QUESTIONS = "Сделай короткие вопросы (от 6 до 10) для создания "

var USER_PROMPT_HEADING = "Придумай короткий заголовок этого материала так, чтобы он отражал"

type RequestBody struct {
	RequestMessage string    `json:"request_message"`
	TypeOfStory    string    `json:"type_of_story"`
	PeriodOfLife   string    `json:"period_of_life"`
	Previous       string    `json:"previous"`
	HumanInfo      HumanInfo `json:"human_info"`
}

type ShortBody struct {
	PartYoung  string    `json:"part_young"`
	PartMiddle string    `json:"part_middle"`
	PartOld    string    `json:"part_old"`
	HumanInfo  HumanInfo `json:"human_info"`
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
