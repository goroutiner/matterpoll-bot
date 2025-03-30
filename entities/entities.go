package entities

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

type Poll struct {
	PollId   string           `json:"id"`
	Question string           `json:"question"`
	Options  map[string]int32 `json:"options"`
	Voters   map[string]bool
	Creator  string `json:"creator"`
	Closed   bool   `json:"closed"`
}

type Voice struct {
	PollId string
	UserId string
	Option string
}

type CommandInfo struct {
	Trigger     string
	DisplayName string
	Description string
	Hint        string
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
	Bot         *model.Client4
	SpaceName   = "polls"
	CommandList = []CommandInfo{
		{"poll-create", "Create poll", "Create a new poll", "`[\"question\"]` `[\"option1\"]` `[\"option2\"]` ..."},
		{"poll-vote", "Vote", "Сast a vote", "`[\"poll_id\"]` `[\"option\"]`"},
		{"poll-results", "Results", "Get poll results", "`[\"poll_id\"]`"},
		{"poll-close", "Close poll", "Close an active poll", "`[\"poll_id\"]`"},
		{"poll-delete", "Delete poll", "Delete an exists poll", "`[\"poll_id\"]`"},
	}
)
