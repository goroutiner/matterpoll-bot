package services_test

import (
	"errors"
	"fmt"
	"testing"

	"matterpoll-bot/entities"
	"matterpoll-bot/internal/services"
	"matterpoll-bot/internal/storage/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatePoll(t *testing.T) {
	mockStore := mocks.NewStoreInterface(t)
	pollService := services.NewPollService(mockStore)

	poll := &entities.Poll{
		PollId:   "poll1",
		Question: "What is your favorite color?",
		Options:  map[string]int32{"Red": 0, "Blue": 0},
		Creator:  "user1",
		Closed:   false,
	}

	t.Run("Success created Poll", func(t *testing.T){
		mockStore.On("CreatePoll", mock.Anything).Return(nil)

		err := pollService.CreatePoll(poll)
		require.NoError(t, err)
		mockStore.AssertCalled(t, "CreatePoll", poll)	
	})

	t.Run("Failed created Poll", func(t *testing.T){
		mockStore.ExpectedCalls = nil
		mockStore.On("CreatePoll", mock.Anything).Return(errors.New("failed to create poll"))
	
		err := pollService.CreatePoll(poll)
		require.Error(t, err)
		require.Equal(t, "failed to create poll", err.Error())
		mockStore.AssertCalled(t, "CreatePoll", poll)	
	})
}

func TestVote(t *testing.T) {
	mockStore := mocks.NewStoreInterface(t)
	pollService := services.NewPollService(mockStore)

	voice := &entities.Voice{
		PollId: "poll1",
		UserId: "user1",
		Option: "Red",
	}

	t.Run("Success Vote", func(t *testing.T) {
		mockStore.On("Vote", mock.Anything).Return("**Voice recorded!**", nil)

		msg, err := pollService.Vote(voice)
		require.NoError(t, err)
		require.Equal(t, "**Voice recorded!**", msg)
	})

	t.Run("Failed Vote", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		mockStore.On("Vote", mock.Anything).Return("", errors.New("**Invalid Poll_ID or not exists!**"))

		msg, err := pollService.Vote(voice)
		require.Empty(t, msg)
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
		mockStore.AssertCalled(t, "Vote", voice)
	})
}

func TestClosePoll(t *testing.T) {
	mockStore := mocks.NewStoreInterface(t)
	pollService := services.NewPollService(mockStore)

	pollId := "poll1"
	userId := "user1"

	t.Run("Success closed Poll", func(t *testing.T) {
		mockStore.On("ClosePoll", mock.Anything, mock.Anything).Return(fmt.Sprintf("*Poll*: `%s` **has been successfully closed!**", pollId), nil)

		msg, err := pollService.ClosePoll(pollId, userId)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has been successfully closed!**", pollId), msg)
		mockStore.AssertCalled(t, "ClosePoll", pollId, userId)

	})

	t.Run("Failed closed Poll", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		mockStore.On("ClosePoll", mock.Anything, mock.Anything).Return("", fmt.Errorf("*Poll*: `%s` **has already been closed!**", pollId))

		msg, err := pollService.ClosePoll(pollId, userId)
		require.Error(t, err)
		require.Empty(t, msg)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has already been closed!**", pollId), err.Error())
		mockStore.AssertCalled(t, "ClosePoll", pollId, userId)
	})
}

func TestDeletePoll(t *testing.T) {
	mockStore := mocks.NewStoreInterface(t)
	pollService := services.NewPollService(mockStore)

	pollId := "poll1"
	userId := "user1"

	// Успешное удаление опроса
	t.Run("Success deleted Poll", func(t *testing.T) {
		mockStore.On("DeletePoll", mock.Anything, mock.Anything).Return(fmt.Sprintf("*Poll*: `%s` **has been successfully deleted!**", pollId), nil)

		msg, err := pollService.DeletePoll(pollId, userId)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has been successfully deleted!**", pollId), msg)
		mockStore.AssertCalled(t, "DeletePoll", pollId, userId)

	})

	t.Run("Failed closed Poll", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		mockStore.On("DeletePoll", mock.Anything, mock.Anything).Return("", fmt.Errorf("**Invalid Poll_ID or not exists!**"))

		msg, err := pollService.DeletePoll(pollId, userId)
		require.Error(t, err)
		require.Empty(t, msg)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
		mockStore.AssertCalled(t, "DeletePoll", pollId, userId)
	})
}

func TestGetPollResult(t *testing.T) {
	mockStore := mocks.NewStoreInterface(t)
	pollService := services.NewPollService(mockStore)

	pollId := "poll1"

	t.Run("Success got Poll results", func(t *testing.T) {
		mockStore.On("GetPollResult", mock.Anything).Return("**Poll Results:** Red: 5, Blue: 3", nil)

		result, err := pollService.GetPollResult(pollId)
		require.NoError(t, err)
		require.Equal(t, "**Poll Results:** Red: 5, Blue: 3", result)
		mockStore.AssertCalled(t, "GetPollResult", pollId)
	})

	t.Run("Failed got Poll results", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		mockStore.On("GetPollResult", mock.Anything).Return("", fmt.Errorf("**Invalid Poll_ID or not exists!**"))

		result, err := pollService.GetPollResult(pollId)
		require.Error(t, err)
		require.Empty(t, result)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
		mockStore.AssertCalled(t, "GetPollResult", pollId)

	})
}
