package entities

// Poll представляет сущность опроса.
// Поля структуры:
// - PollId: уникальный идентификатор опроса.
// - Question: текст вопроса опроса.
// - Options: варианты ответа с количеством голосов за каждый вариант.
// Poll представляет сущность опроса.
type Poll struct {
	PollId   string           // PollId - уникальный идентификатор опроса
	Question string           // Question - текст вопроса опроса.
	Options  map[string]int32 // Options - варианты ответа с количеством голосов за каждый вариант.
	Voters   map[string]bool  // Voters: список пользователей, проголосовавших в опросе (по идентификатору).
	Creator  string           // Creator - идентификатор создателя опроса.
	Closed   bool             // Closed - флаг, указывающий, закрыт ли опрос.

}

// Voice представляет сущность голоса пользователя в опросе.
type Voice struct {
	PollId string // PollId - уникальный идентификатор опроса
	UserId string // UserId - идентификатор пользователя, который проголосовал.
	Option string // Option - выбранный пользователем вариант ответа.

}

// CommandInfo представляет информацию о команде бота.
type CommandInfo struct {
	Trigger     string // Trigger - триггер команды, который пользователь вводит для её вызова.
	URLPath     string // URLPath - путь URL, связанный с командой, для обработки запросов.
	DisplayName string // DisplayName - отображаемое имя команды, показываемое пользователю.
	Description string // Description - описание команды, объясняющее её функциональность.
	Hint        string // Hint - подсказка или дополнительная информация о том, как использовать команду.
}

type TarantoolConfig struct {
	Address  string
	User     string
	Password string
}

// UserError представляет ошибки, которые можно показывать пользователям.
type UserError struct {
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}

// NewUserError создает новую пользовательскую ошибку.
func NewUserError(message string) error {
	return &UserError{Message: message}
}

var (
	PollsSpaceName  = "polls"      // PollsSpaceName - имя пространства для хранения опросов в Tarantool.
	TokensSpaceName = "cmd_tokens" // PollsSpaceName - имя пространства для хранения токенов команд в Tarantool.
	CommandList     = []CommandInfo{
		{"poll-create", "/poll-create", "Create poll", "Create a new poll", "[\"question\"] [\"option1\"] [\"option2\"] ..."},
		{"poll-vote", "/poll-vote", "Vote", "Сast a vote", "[\"poll_id\"] [\"option\"]"},
		{"poll-results", "/poll-results", "Results", "Get poll results", "[\"poll_id\"]"},
		{"poll-close", "/poll-close", "Close poll", "Close an active poll", "[\"poll_id\"]"},
		{"poll-delete", "/poll-delete", "Delete poll", "Delete an exists poll", "[\"poll_id\"]"},
	}
)
