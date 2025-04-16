package services

import "github.com/mattermost/mattermost-server/v6/model"

// BotInterface определяет интерфейс для взаимодействия с ботом Mattermost.
type BotInterface interface {
	SetToken(token string)
	GetTeamByName(teamName, etag string) (*model.Team, *model.Response, error)
	ListCommands(teamId string, customOnly bool) ([]*model.Command, *model.Response, error)
	CreateCommand(cmd *model.Command) (*model.Command, *model.Response, error)
	CreatePost(post *model.Post) (*model.Post, *model.Response, error)
}
