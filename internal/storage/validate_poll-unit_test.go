package storage_test

import (
	"fmt"
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/storage"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestValidateVoice тестирует функционал проверки корректности голосования в опросе.
func TestValidateVoice(t *testing.T) {
    t.Run("Invalid option", func(t *testing.T) {
        poll := &entities.Poll{
            Options: map[string]int32{"1": 1},
            Voters:  map[string]bool{},
            Closed:  false,
        }
        voice := &entities.Voice{
            Option: "2",
            UserId: "user1",
            PollId: "poll1",
        }

        err := storage.ValidateVoice(poll, voice)
        require.Error(t, err)
        require.Equal(t, "**Invalid option!**", err.Error())
    })

    t.Run("Repeat voice", func(t *testing.T) {
        poll := &entities.Poll{
            Options: map[string]int32{"1": 1},
            Voters:  map[string]bool{"user1": true},
            Closed:  false,
        }
        voice := &entities.Voice{
            Option: "1",
            UserId: "user1",
            PollId: "poll1",
        }

        err := storage.ValidateVoice(poll, voice)
        require.Error(t, err)
        require.Equal(t, "**You can't vote again!**", err.Error())
    })

    t.Run("Closed poll", func(t *testing.T) {
        poll := &entities.Poll{
            Options: map[string]int32{"1": 1},
            Voters:  map[string]bool{},
            Closed:  true,
        }
        voice := &entities.Voice{
            Option: "1",
            UserId: "user1",
            PollId: "poll1",
        }

        err := storage.ValidateVoice(poll, voice)
        require.Error(t, err)
        require.Equal(t, fmt.Sprintf("*Poll*: `%s` **is already closed!**", voice.PollId), err.Error())
    })

    t.Run("Valid vote", func(t *testing.T) {
        poll := &entities.Poll{
            Options: map[string]int32{"1": 1},
            Voters:  map[string]bool{},
            Closed:  false,
        }
        voice := &entities.Voice{
            Option: "1",
            UserId: "user1",
            PollId: "poll1",
        }

        err := storage.ValidateVoice(poll, voice)
        require.NoError(t, err)
    })
}
