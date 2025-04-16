package services

import (
	"fmt"
	"log"
	"matterpoll-bot/config"
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/storage"

	"github.com/mattermost/mattermost-server/v6/model"
)

type PollService struct {
	Bot   BotInterface
	store storage.StoreInterface
}

// NewPollService возвращает структуру сервиса голосований.
func NewPollService(bot BotInterface, s storage.StoreInterface) *PollService {
	return &PollService{Bot: bot, store: s}
}

// CreatePoll создает новый опрос и сохраняет его в хранилище.
// Возвращает строку с результатом и ошибку, если операция завершилась неудачно.
func (ps *PollService) CreatePoll(poll *entities.Poll) error {
	err := ps.store.CreatePoll(poll)

	return err
}

// Vote регистрирует голос пользователя в опросе,
// в соответствии с выбранным вариантом.
// Возвращает строку с результатом и ошибку, если операция завершилась неудачно.
func (ps *PollService) Vote(voice *entities.Voice) (string, error) {
	res, err := ps.store.Vote(voice)

	return res, err
}

// GetPollResult получает результат опроса по его идентификатору.
// Возвращает строку с результатом и ошибку, если операция завершилась неудачно.
func (ps *PollService) GetPollResult(pollId string) (string, error) {
	res, err := ps.store.GetPollResult(pollId)

	return res, err
}

// ClosePoll завершает опрос с указанным pollId от имени пользователя userId.
// Возвращает строку с результатом и ошибку, если операция завершилась неудачно.
func (ps *PollService) ClosePoll(pollId, userId string) (string, error) {
	res, err := ps.store.ClosePoll(pollId, userId)

	return res, err
}

// DeletePoll удаляет опрос с указанным pollId, если userId имеет необходимые права.
// Возвращает строку с результатом и ошибку, если операция завершилась неудачно.
func (ps *PollService) DeletePoll(pollId, userId string) (string, error) {
	res, err := ps.store.DeletePoll(pollId, userId)

	return res, err
}

// RegisterCommands регистрирует команды Mattermost для бота.
// Сначала проверяется наличие команды в списке существующих команд,
// затем создаются новые команды, если они еще не зарегистрированы.
func (ps *PollService) RegisterCommands() error {
	team, resp, err := ps.Bot.GetTeamByName(config.TeamName, "")
	if err != nil {
		return fmt.Errorf("failed to get team: %w", err)
	}

	if resp == nil || resp.StatusCode != 200 {
		return fmt.Errorf("failed to get team: unexpected status code %d", resp.StatusCode)
	}

	existingCommands, resp, err := ps.Bot.ListCommands(team.Id, false)
	if err != nil {
		return fmt.Errorf("failed to get commands list: %w", err)
	}

	if resp == nil || resp.StatusCode != 200 {
		return fmt.Errorf("failed to get commands list: unexpected status code %d", resp.StatusCode)
	}

	registeredCommands := make(map[string]bool)
	for _, cmd := range existingCommands {
		registeredCommands[cmd.Trigger] = true
	}

	for _, cmd := range entities.CommandList {
		if registeredCommands[cmd.Trigger] {
			continue
		}

		newCommand := &model.Command{
			TeamId:           team.Id,
			Trigger:          cmd.Trigger,
			Method:           "P",
			URL:              fmt.Sprintf("http://%s%s%s", config.BotHostname, config.BotSocket, cmd.URLPath),
			DisplayName:      cmd.DisplayName,
			Description:      cmd.Description,
			AutoComplete:     true,
			AutoCompleteDesc: cmd.Description,
			AutoCompleteHint: cmd.Hint,
		}

		createdCommand, resp, err := ps.Bot.CreateCommand(newCommand)
		if err != nil {
			return fmt.Errorf("failed to create command '%s': %w", cmd.URLPath, err)
		}

		if resp == nil || resp.StatusCode != 201 {
			return fmt.Errorf("failed to create command: unexpected status code %d", resp.StatusCode)
		}

		if err := ps.store.AddCmdToken(cmd.URLPath, createdCommand.Token); err != nil {
			return fmt.Errorf("failed to add cmd token : %w", err)
		}

		log.Printf("Created command '%s' with token: %s\n", cmd.URLPath, createdCommand.Token)
	}

	return nil
}
