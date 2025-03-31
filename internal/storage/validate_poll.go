package storage

import (
	"fmt"
	"matterpoll-bot/internal/entities"
)

func ValidateVoice(poll *entities.Poll, voice *entities.Voice) error {
	if _, existsOption := poll.Options[voice.Option]; !existsOption {
		return entities.NewUserError("**Invalid option!**")
	}
	if poll.Voters[voice.UserId] {
		return entities.NewUserError("**You can't vote again!**")
	}
	if poll.Closed {
		return entities.NewUserError(fmt.Sprintf("*Poll*: `%s` **is already closed!**", voice.PollId))
	}
	return nil
}
