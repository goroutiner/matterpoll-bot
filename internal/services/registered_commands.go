package services

import (
	"fmt"
	"matterpoll-bot/config"
	"matterpoll-bot/internal/entities"

	"github.com/mattermost/mattermost-server/v6/model"
)

// RegisterCommands регистрирует команды Mattermost для бота.
// Сначала проверяется наличие команды в списке существующих команд,
// затем создаются новые команды, если они еще не зарегистрированы.
func RegisterCommands() error {
	team, _, err := entities.Bot.GetTeamByName(config.TeamName, "")
	if err != nil {
		return fmt.Errorf("failed to get team: %v", err)
	}

	existingCommands, _, err := entities.Bot.ListCommands(team.Id, false)
	if err != nil {
		return fmt.Errorf("failed to list commands: %v", err)
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
			URL:              fmt.Sprintf("https://matterpoll-bot%s%s", config.BotSocket, cmd.URLPath),
			DisplayName:      cmd.DisplayName,
			Description:      cmd.Description,
			AutoComplete:     true,
			AutoCompleteDesc: cmd.Description,
			AutoCompleteHint: cmd.Hint,
		}

		createdCommand, resp, err := entities.Bot.CreateCommand(newCommand)
		if err != nil || resp.StatusCode != 201 {
			return fmt.Errorf("failed to create command '/%s': %v", cmd.Trigger, err)
		}

		fmt.Printf("Created command '/%s' with token: %s\n", cmd.Trigger, createdCommand.Token)
	}

	return nil
}
