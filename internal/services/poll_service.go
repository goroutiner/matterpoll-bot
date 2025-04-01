package services

import (
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/storage"
)

type PollService struct {
	store storage.StoreInterface
}

// NewPollService возвращает структуру сервиса голосований.
func NewPollService(s storage.StoreInterface) *PollService {
	return &PollService{store: s}
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
