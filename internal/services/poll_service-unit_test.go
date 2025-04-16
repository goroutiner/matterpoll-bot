package services_test

import (
	"errors"
	"fmt"
	"testing"

	"matterpoll-bot/config"
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/services"
	"matterpoll-bot/internal/services/service_mocks"
	"matterpoll-bot/internal/storage/store_mocks"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestCreatePoll проверяет функциональность создания опроса.
func TestCreatePoll(t *testing.T) {
	mockStore := store_mocks.NewStoreInterface(t)
	pollService := services.NewPollService(nil, mockStore)

	poll := &entities.Poll{
		PollId:   "poll1",
		Question: "What is your favorite color?",
		Options:  map[string]int32{"Red": 0, "Blue": 0},
		Creator:  "user1",
		Closed:   false,
	}

	t.Run("success created Poll", func(t *testing.T) {
		mockStore.On("CreatePoll", mock.Anything).Return(nil)

		err := pollService.CreatePoll(poll)
		require.NoError(t, err)
		mockStore.AssertCalled(t, "CreatePoll", poll)
	})

	t.Run("failed created Poll", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		mockStore.On("CreatePoll", mock.Anything).Return(errors.New("failed to create poll"))

		err := pollService.CreatePoll(poll)
		require.Error(t, err)
		require.Equal(t, "failed to create poll", err.Error())
		mockStore.AssertCalled(t, "CreatePoll", poll)
	})
}

// TestVote проверяет функциональность голосования в опросе.
func TestVote(t *testing.T) {
	mockStore := store_mocks.NewStoreInterface(t)
	pollService := services.NewPollService(nil, mockStore)

	voice := &entities.Voice{
		PollId: "poll1",
		UserId: "user1",
		Option: "Red",
	}

	t.Run("success Vote", func(t *testing.T) {
		mockStore.On("Vote", mock.Anything).Return("**Voice recorded!**", nil)

		msg, err := pollService.Vote(voice)
		require.NoError(t, err)
		require.Equal(t, "**Voice recorded!**", msg)
	})

	t.Run("failed Vote", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		mockStore.On("Vote", mock.Anything).Return("", errors.New("**Invalid Poll_ID or not exists!**"))

		msg, err := pollService.Vote(voice)
		require.Empty(t, msg)
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
		mockStore.AssertCalled(t, "Vote", voice)
	})
}

// TestClosePoll проверяет функциональность закрытия опроса.
func TestClosePoll(t *testing.T) {
	mockStore := store_mocks.NewStoreInterface(t)
	pollService := services.NewPollService(nil, mockStore)

	pollId := "poll1"
	userId := "user1"

	t.Run("success closed Poll", func(t *testing.T) {
		mockStore.On("ClosePoll", mock.Anything, mock.Anything).Return(fmt.Sprintf("*Poll*: `%s` **has been successfully closed!**", pollId), nil)

		msg, err := pollService.ClosePoll(pollId, userId)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has been successfully closed!**", pollId), msg)
		mockStore.AssertCalled(t, "ClosePoll", pollId, userId)
	})

	t.Run("failed closed Poll", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		mockStore.On("ClosePoll", mock.Anything, mock.Anything).Return("", fmt.Errorf("*Poll*: `%s` **has already been closed!**", pollId))

		msg, err := pollService.ClosePoll(pollId, userId)
		require.Error(t, err)
		require.Empty(t, msg)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has already been closed!**", pollId), err.Error())
		mockStore.AssertCalled(t, "ClosePoll", pollId, userId)
	})
}

// TestDeletePoll проверяет функциональность удаления опроса.
func TestDeletePoll(t *testing.T) {
	mockStore := store_mocks.NewStoreInterface(t)
	pollService := services.NewPollService(nil, mockStore)

	pollId := "poll1"
	userId := "user1"

	t.Run("success deleted Poll", func(t *testing.T) {
		mockStore.On("DeletePoll", mock.Anything, mock.Anything).Return(fmt.Sprintf("*Poll*: `%s` **has been successfully deleted!**", pollId), nil)

		msg, err := pollService.DeletePoll(pollId, userId)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has been successfully deleted!**", pollId), msg)
		mockStore.AssertCalled(t, "DeletePoll", pollId, userId)

	})

	t.Run("failed closed Poll", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		mockStore.On("DeletePoll", mock.Anything, mock.Anything).Return("", fmt.Errorf("**Invalid Poll_ID or not exists!**"))

		msg, err := pollService.DeletePoll(pollId, userId)
		require.Error(t, err)
		require.Empty(t, msg)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
		mockStore.AssertCalled(t, "DeletePoll", pollId, userId)
	})
}

// TestGetPollResult проверяет функциональность получения результатов опроса.
func TestGetPollResult(t *testing.T) {
	mockStore := store_mocks.NewStoreInterface(t)
	pollService := services.NewPollService(nil, mockStore)

	pollId := "poll1"

	t.Run("success got Poll results", func(t *testing.T) {
		mockStore.On("GetPollResult", mock.Anything).Return("**Poll Results:** Red: 5, Blue: 3", nil)

		result, err := pollService.GetPollResult(pollId)
		require.NoError(t, err)
		require.Equal(t, "**Poll Results:** Red: 5, Blue: 3", result)
		mockStore.AssertCalled(t, "GetPollResult", pollId)
	})

	t.Run("failed got Poll results", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		mockStore.On("GetPollResult", mock.Anything).Return("", fmt.Errorf("**Invalid Poll_ID or not exists!**"))

		result, err := pollService.GetPollResult(pollId)
		require.Error(t, err)
		require.Empty(t, result)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
		mockStore.AssertCalled(t, "GetPollResult", pollId)
	})
}

// TestRegisterCommands проверяет функциональность регистрации команд.
func TestRegisterCommands(t *testing.T) {
	mockStore := store_mocks.NewStoreInterface(t)
	mockBot := service_mocks.NewBotInterface(t)
	pollService := services.NewPollService(mockBot, mockStore)

	config.TeamName = "test_team"
	config.BotHostname = "localhost"
	config.BotSocket = ":8080"

	team := &model.Team{Id: "team_id"}
	existingCommands := []*model.Command{
		{Trigger: "existing_command"},
	}
	newCommand := entities.CommandInfo{
		Trigger:     "new_command",
		URLPath:     "/new_command",
		DisplayName: "New Command",
		Description: "Description of new command",
		Hint:        "[hint]",
	}
	entities.CommandList = []entities.CommandInfo{newCommand}

	t.Run("successfully registers new commands", func(t *testing.T) {
		getStatusCode := 200
		createStatusCode := 201

		mockBot.On("GetTeamByName", config.TeamName, "").Return(team, &model.Response{StatusCode: getStatusCode}, nil)
		mockBot.On("ListCommands", team.Id, false).Return(existingCommands, &model.Response{StatusCode: getStatusCode}, nil)
		mockBot.On("CreateCommand", mock.Anything).Return(&model.Command{Token: "new_token"}, &model.Response{StatusCode: createStatusCode}, nil)
		mockStore.On("AddCmdToken", newCommand.URLPath, "new_token").Return(nil)

		err := pollService.RegisterCommands()
		require.NoError(t, err)

		mockBot.AssertCalled(t, "GetTeamByName", config.TeamName, "")
		mockBot.AssertCalled(t, "ListCommands", team.Id, false)
		mockBot.AssertCalled(t, "CreateCommand", mock.MatchedBy(func(cmd *model.Command) bool {
			return cmd.Trigger == newCommand.Trigger
		}))
		mockStore.AssertCalled(t, "AddCmdToken", newCommand.URLPath, "new_token")
	})

	t.Run("failed to get team", func(t *testing.T) {
		// Проверяем обработку ошибки при получении команды
		getStatusCode := 200
		testErr := errors.New("error text")

		mockBot.ExpectedCalls = nil

		mockBot.On("GetTeamByName", config.TeamName, "").Return(nil, &model.Response{StatusCode: getStatusCode}, testErr)

		err := pollService.RegisterCommands()
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("failed to get team: %v", testErr), err.Error())

		// Проверяем обработку 500 статуса ответа при получении команды
		getStatusCode = 500

		mockBot.ExpectedCalls = nil

		mockBot.On("GetTeamByName", config.TeamName, "").Return(nil, &model.Response{StatusCode: getStatusCode}, nil)

		err = pollService.RegisterCommands()
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("failed to get team: unexpected status code %d", getStatusCode), err.Error())

		mockBot.AssertCalled(t, "GetTeamByName", config.TeamName, "")
	})

	t.Run("failed to get commands list", func(t *testing.T) {
		// Проверяем обработку ошибки при получении списка команд
		mockBot.ExpectedCalls = nil

		getStatusCode := 200

		mockBot.On("GetTeamByName", config.TeamName, "").Return(team, &model.Response{StatusCode: getStatusCode}, nil)

		testErr := errors.New("error text")
		mockBot.On("ListCommands", team.Id, false).Return(existingCommands, &model.Response{StatusCode: getStatusCode}, testErr)

		err := pollService.RegisterCommands()
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("failed to get commands list: %v", testErr), err.Error())

		// Проверяем обработку 500 статуса ответа при получении списка команд
		mockBot.ExpectedCalls = nil

		mockBot.On("GetTeamByName", config.TeamName, "").Return(team, &model.Response{StatusCode: getStatusCode}, nil)

		getStatusCode = 500
		mockBot.On("ListCommands", team.Id, false).Return(existingCommands, &model.Response{StatusCode: getStatusCode}, nil)

		err = pollService.RegisterCommands()
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("failed to get commands list: unexpected status code %d", getStatusCode), err.Error())

		mockBot.AssertCalled(t, "GetTeamByName", config.TeamName, "")
		mockBot.AssertCalled(t, "ListCommands", team.Id, false)
	})

	t.Run("failed to create commands", func(t *testing.T) {
		// Проверяем обработку ошибки при создании команд
		mockBot.ExpectedCalls = nil

		getStatusCode := 200
		createStatusCode := 500
		testErr := errors.New("error text")

		mockBot.On("GetTeamByName", config.TeamName, "").Return(team, &model.Response{StatusCode: getStatusCode}, nil)
		mockBot.On("ListCommands", team.Id, false).Return(existingCommands, &model.Response{StatusCode: getStatusCode}, nil)
		mockCreateCommand := mockBot.On("CreateCommand", mock.Anything).Return(&model.Command{Token: "new_token"}, &model.Response{StatusCode: createStatusCode}, testErr)

		err := pollService.RegisterCommands()
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("failed to create command '%s': %v", newCommand.URLPath, testErr), err.Error())

		// Проверяем обработку 500 статуса ответа при создании команд
		mockCreateCommand.Unset()

		mockCreateCommand.On("CreateCommand", mock.Anything).Return(&model.Command{Token: "new_token"}, &model.Response{StatusCode: createStatusCode}, nil)

		err = pollService.RegisterCommands()
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("failed to create command: unexpected status code %d", createStatusCode), err.Error())

		mockBot.AssertCalled(t, "GetTeamByName", config.TeamName, "")
		mockBot.AssertCalled(t, "ListCommands", team.Id, false)
		mockBot.AssertCalled(t, "CreateCommand", mock.MatchedBy(func(cmd *model.Command) bool {
			return cmd.Trigger == newCommand.Trigger
		}))
	})

	t.Run("failed to add cmd token", func(t *testing.T) {
		// Проверяем обработку ошибки при добавлении токенов команд
		mockBot.ExpectedCalls = nil
		mockStore.ExpectedCalls = nil

		getStatusCode := 200
		createStatusCode := 201

		mockBot.On("GetTeamByName", config.TeamName, "").Return(team, &model.Response{StatusCode: getStatusCode}, nil)
		mockBot.On("ListCommands", team.Id, false).Return(existingCommands, &model.Response{StatusCode: getStatusCode}, nil)
		mockBot.On("CreateCommand", mock.Anything).Return(&model.Command{Token: "new_token"}, &model.Response{StatusCode: createStatusCode}, nil)

		testErr := errors.New("error text")
		mockStore.On("AddCmdToken", newCommand.URLPath, "new_token").Return(testErr)

		err := pollService.RegisterCommands()
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("failed to add cmd token : %v", testErr), err.Error())

		mockBot.AssertCalled(t, "GetTeamByName", config.TeamName, "")
		mockBot.AssertCalled(t, "ListCommands", team.Id, false)
		mockBot.AssertCalled(t, "CreateCommand", mock.MatchedBy(func(cmd *model.Command) bool {
			return cmd.Trigger == newCommand.Trigger
		}))
		mockStore.AssertCalled(t, "AddCmdToken", newCommand.URLPath, "new_token")
	})
}
