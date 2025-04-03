package memory

import (
	"fmt"
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/storage"
	"sync"
)

type Memory struct {
	polls     map[string]*entities.Poll
	cmdTokens map[string]string
	mu        sync.RWMutex
}

// NewMemoryStore возвращает структуру для хранения ссылок во внутренней памяти.
func NewMemoryStore() *Memory {
	return &Memory{
		polls:     map[string]*entities.Poll{},
		cmdTokens: map[string]string{},
	}
}

// CreatePoll сохраняет новый опрос во внутренней памяти.
func (m *Memory) CreatePoll(poll *entities.Poll) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.polls[poll.PollId] = poll

	return nil
}

// Vote регистрирует голос пользователя в опросе,
// в соответствии с выбранным вариантом и обновляет данные во внутренней памяти.
func (m *Memory) Vote(voice *entities.Voice) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	poll := m.polls[voice.PollId]
	if poll == nil {
		return "", entities.NewUserError("**Invalid Poll_ID or not exists!**")
	}

	if err := storage.ValidateVoice(poll, voice); err != nil {
		return "", err
	}

	poll.Options[voice.Option]++
	poll.Voters[voice.UserId] = true

	return "**Voice recorded!**", nil
}

// GetPollResult получает результаты опроса из внутренней памяти.
func (m *Memory) GetPollResult(pollId string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	poll, err := m.getPoll(pollId)
	if err != nil {
		return "", err
	}

	tbl := storage.PrintTable(poll)

	return tbl, nil
}

// ClosePoll закрывает опрос и обновляет данные во внутренней памяти.
func (m *Memory) ClosePoll(pollId, userId string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	poll, err := m.getPoll(pollId)
	if err != nil {
		return "", err
	}

	if poll.Closed {
		return "", entities.NewUserError(fmt.Sprintf("*Poll*: `%s` **has already been closed!**", pollId))
	}
	if poll.Creator != userId {
		return "", entities.NewUserError("**You don't have the permission to close a vote!**")
	}
	poll.Closed = true

	return fmt.Sprintf("*Poll*: `%s` **has been successfully closed!**", pollId), nil
}

// DeletePoll удаляет опрос из внутренней памяти.
func (m *Memory) DeletePoll(pollId, userId string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	poll, err := m.getPoll(pollId)
	if err != nil {
		return "", err
	}

	if poll.Creator != userId {
		return "", entities.NewUserError("**You don't have the permission to delete a vote!**")
	}
	delete(m.polls, pollId)

	return fmt.Sprintf("*Poll*: `%s` **has been successfully delete!**", pollId), nil
}

func (m *Memory) AddCmdToken(cmdPath, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cmdTokens[cmdPath] = token

	return nil
}

func (m *Memory) ValidateCmdToken(cmdPath, token string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, has := m.cmdTokens[cmdPath]

	return has
}

// getPoll получает структуру опроса по Id и проверяет ее существование.
func (m *Memory) getPoll(pollId string) (*entities.Poll, error) {
	poll := m.polls[pollId]
	if poll == nil {
		return nil, entities.NewUserError("**Invalid Poll_ID or not exists!**")
	}
	
	return poll, nil
}
