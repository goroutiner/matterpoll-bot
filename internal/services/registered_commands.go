package services

import (
	"fmt"
	"matterpoll-bot/config"
	"matterpoll-bot/entities"

	"github.com/mattermost/mattermost-server/v6/model"
)

func RegisterCommands() error {
	existingCommands, resp, err := entities.Bot.ListCommands(config.TeamId, false)
	if err != nil || resp.StatusCode != 200 {
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
			TeamId:           config.TeamId,
			Trigger:          cmd.Trigger,
			Method:           "P",
			URL:              config.ServerURL,
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
