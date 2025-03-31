package services

import (
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/storage"
)

type PollService struct {
	store storage.StoreInterface
}

func NewPollService(s storage.StoreInterface) *PollService {
	return &PollService{store: s}
}

func (ps *PollService) CreatePoll(poll *entities.Poll) error {
	err := ps.store.CreatePoll(poll)

	return err
}
func (ps *PollService) Vote(voice *entities.Voice) (string, error) {
	res, err := ps.store.Vote(voice)

	return res, err
}

func (ps *PollService) GetPollResult(pollId string) (string, error) {
	res, err := ps.store.GetPollResult(pollId)

	return res, err
}

func (ps *PollService) ClosePoll(pollId, userId string) (string, error) {
	res, err := ps.store.ClosePoll(pollId, userId)

	return res, err
}

func (ps *PollService) DeletePoll(pollId, userId string) (string, error) {
	res, err := ps.store.DeletePoll(pollId, userId)

	return res, err
}
