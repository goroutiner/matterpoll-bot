package storage

import "matterpoll-bot/internal/entities"

// StoreInterface определяет интерфейс для работы с хранилищем данных,
// используемым в приложении для управления опросами.
type StoreInterface interface {
	CreatePoll(poll *entities.Poll) error
	Vote(voice *entities.Voice) (string, error)
	GetPollResult(pollId string) (string, error)
	ClosePoll(pollId, userId string) (string, error)
	DeletePoll(pollId, userId string) (string, error)
	AddCmdToken(cmdPath, token string) error
	ValidateCmdToken(cmdPath, token string) bool
}
