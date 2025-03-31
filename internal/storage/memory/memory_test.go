package memory

import (
	"fmt"
	"matterpoll-bot/internal/entities"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePoll(t *testing.T) {
	store := NewMemoryStore()

	poll := &entities.Poll{
		PollId:  "poll1",
		Options: map[string]int32{"option1": 0, "option2": 0},
		Voters:  map[string]bool{},
		Creator: "user1",
	}

	err := store.CreatePoll(poll)
	require.NoError(t, err)

	storedPoll, exists := store.polls["poll1"]
	require.True(t, exists)
	require.Equal(t, poll, storedPoll)
}
func TestClosePoll(t *testing.T) {
	store := NewMemoryStore()

	poll := &entities.Poll{
		PollId:  "poll1",
		Options: map[string]int32{"option1": 0, "option2": 0},
		Voters:  map[string]bool{},
		Creator: "user1",
		Closed:  false,
	}

	err := store.CreatePoll(poll)
	require.NoError(t, err)

	t.Run("Success closed Poll", func(t *testing.T) {
		pollId := "poll1"
		userId := "user1"
		msg, err := store.ClosePoll(pollId, userId)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has been successfully closed!**", poll.PollId), msg)

		closedPoll, exists := store.polls["poll1"]
		require.True(t, exists)
		require.True(t, closedPoll.Closed)
	})

	t.Run("Already Closed", func(t *testing.T) {
		pollId := "poll1"
		userId := "user1"
		_, err := store.ClosePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has already been closed!**", poll.PollId), err.Error())
	})

	t.Run("Invalid PollId", func(t *testing.T) {
		pollId := "invalid_poll"
		userId := "user1"
		_, err := store.ClosePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
	})

	t.Run("Don't have the permission", func(t *testing.T) {
		pollId := "poll1"
		userId := "user2"
		poll.Closed = false
		msg, err := store.ClosePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, "**You don't have the permission to close a vote!**", err.Error())
		require.Empty(t, msg)
	})
}
func TestVote(t *testing.T) {
	store := NewMemoryStore()

	poll := &entities.Poll{
		PollId:  "poll1",
		Options: map[string]int32{"option1": 0, "option2": 0},
		Voters:  map[string]bool{},
		Creator: "user1",
		Closed:  false,
	}

	err := store.CreatePoll(poll)
	require.NoError(t, err)

	t.Run("Successful Vote", func(t *testing.T) {
		voice := &entities.Voice{
			PollId: "poll1",
			Option: "option1",
			UserId: "user2",
		}
		msg, err := store.Vote(voice)
		require.NoError(t, err)
		require.Equal(t, "**Voice recorded!**", msg)
		require.Equal(t, int32(1), poll.Options["option1"])
		require.True(t, poll.Voters["user2"])
	})

	t.Run("Invalid PollId", func(t *testing.T) {
		voice := &entities.Voice{
			PollId: "invalid_poll",
			Option: "option1",
			UserId: "user2",
		}
		_, err := store.Vote(voice)
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
	})

	t.Run("Invalid Option", func(t *testing.T) {
		voice := &entities.Voice{
			PollId: "poll1",
			Option: "invalid_option",
			UserId: "user2",
		}
		_, err := store.Vote(voice)
		require.Error(t, err)
		require.Equal(t, "**Invalid option!**", err.Error())
	})

	t.Run("Already Voted", func(t *testing.T) {
		voice := &entities.Voice{
			PollId: "poll1",
			Option: "option1",
			UserId: "user2",
		}
		_, err := store.Vote(voice)
		require.Error(t, err)
		require.Equal(t, "**You can't vote again!**", err.Error())
	})

	t.Run("Poll Closed", func(t *testing.T) {
		poll.Closed = true
		voice := &entities.Voice{
			PollId: "poll1",
			Option: "option1",
			UserId: "user3",
		}
		_, err := store.Vote(voice)
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **is already closed!**", voice.PollId), err.Error())
	})
}
func TestGetPollResult(t *testing.T) {
	store := NewMemoryStore()

	poll := &entities.Poll{
		PollId:  "poll1",
		Options: map[string]int32{"option1": 5, "option2": 3},
		Voters:  map[string]bool{"user1": true, "user2": true},
		Creator: "user1",
		Closed:  false,
	}

	err := store.CreatePoll(poll)
	require.NoError(t, err)

	t.Run("Valid PollId", func(t *testing.T) {
		pollId := "poll1"
		result, err := store.GetPollResult(pollId)
		require.NoError(t, err)
		require.NotEmpty(t, result)
	})

	t.Run("Invalid PollId", func(t *testing.T) {
		pollId := "invalid_poll"
		_, err := store.GetPollResult(pollId)
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
	})
}
func TestDeletePoll(t *testing.T) {
	store := NewMemoryStore()

	poll := &entities.Poll{
		PollId:  "poll1",
		Options: map[string]int32{"option1": 0, "option2": 0},
		Voters:  map[string]bool{},
		Creator: "user1",
		Closed:  false,
	}

	err := store.CreatePoll(poll)
	require.NoError(t, err)

	t.Run("Successful Deletion", func(t *testing.T) {
		pollId := "poll1"
		userId := "user1"
		msg, err := store.DeletePoll(pollId, userId)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has been successfully delete!**", pollId), msg)

		_, exists := store.polls[pollId]
		require.False(t, exists)
	})

	t.Run("Invalid PollId", func(t *testing.T) {
		pollId := "invalid_poll"
		userId := "user1"
		_, err := store.DeletePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
	})

	t.Run("Don't have the permission", func(t *testing.T) {
		poll := &entities.Poll{
			PollId:  "poll2",
			Options: map[string]int32{"option1": 0, "option2": 0},
			Voters:  map[string]bool{},
			Creator: "user1",
			Closed:  false,
		}
		err := store.CreatePoll(poll)
		require.NoError(t, err)

		pollId := "poll2"
		userId := "user2"
		_, err = store.DeletePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, "**You don't have the permission to delete a vote!**", err.Error())
	})
}
func TestGetPoll(t *testing.T) {
	store := NewMemoryStore()

	poll := &entities.Poll{
		PollId:  "poll1",
		Options: map[string]int32{"option1": 0, "option2": 0},
		Voters:  map[string]bool{},
		Creator: "user1",
		Closed:  false,
	}

	err := store.CreatePoll(poll)
	require.NoError(t, err)

	t.Run("Valid PollId", func(t *testing.T) {
		retrievedPoll, err := store.getPoll("poll1")
		require.NoError(t, err)
		require.NotNil(t, retrievedPoll)
		require.Equal(t, poll, retrievedPoll)
	})

	t.Run("Invalid PollId", func(t *testing.T) {
		_, err := store.getPoll("invalid_poll")
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
	})
}
