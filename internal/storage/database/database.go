package database

import (
	"context"
	"fmt"
	"log"
	"matterpoll-bot/config"
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/storage"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	_ "github.com/tarantool/go-tarantool/v2/datetime"
	_ "github.com/tarantool/go-tarantool/v2/decimal"
	_ "github.com/tarantool/go-tarantool/v2/uuid"
)

type Database struct {
	conn *tarantool.Connection
}

// NewDatabaseConection возвращает структуру соединения с БД.
func NewDatabaseStore(conn *tarantool.Connection) *Database {
	return &Database{conn: conn}
}

// NewDatabaseConection создает соединение с БД.
func NewDatabaseConection() (*tarantool.Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  config.DbSocket,
		User:     "user",
		Password: "secret",
	}
	opts := tarantool.Opts{
		Timeout:       1 * time.Minute,
		Reconnect:     500 * time.Millisecond,
		MaxReconnects: 10,
		RateLimit:     100,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		return nil, fmt.Errorf("connection refused: %v", err)
	}

	return conn, err
}

// CreatePoll добавляет новый опрос в базу данных.
func (d *Database) CreatePoll(poll *entities.Poll) error {
	tuple := []interface{}{
		poll.PollId,
		poll.Question,
		poll.Options,
		poll.Voters,
		poll.Creator,
		poll.Closed,
	}

	reqPost := tarantool.NewInsertRequest(entities.PollsSpaceName).Tuple(tuple)
	if _, err := d.conn.Do(reqPost).Get(); err != nil {
		return err
	}

	return nil
}

// Vote регистрирует голос пользователя в опросе,
// в соответствии с выбранным вариантом и обновляет данные БД.
func (d *Database) Vote(voice *entities.Voice) (string, error) {
	reqGet := tarantool.NewSelectRequest(entities.PollsSpaceName).
		Index("primary").
		Iterator(tarantool.IterEq).
		Key([]interface{}{voice.PollId})
	data, err := d.conn.Do(reqGet).Get()
	if err != nil {
		return "", fmt.Errorf("failed to execute select request: %w", err)
	}

	poll, err := parseData(data)
	if err != nil {
		return "", err
	}

	if poll.Closed {
		return "", entities.NewUserError(fmt.Sprintf("*Poll*: `%s` **is closed!**", voice.PollId))
	}

	if err := storage.ValidateVoice(poll, voice); err != nil {
		return "", err
	}

	poll.Options[voice.Option]++
	poll.Voters[voice.UserId] = true

	reqUpd := tarantool.NewUpdateRequest(entities.PollsSpaceName).
		Key([]interface{}{voice.PollId}).
		Operations(tarantool.NewOperations().
			Assign(2, poll.Options).
			Assign(3, poll.Voters))
	if _, err = d.conn.Do(reqUpd).Get(); err != nil {
		return "", fmt.Errorf("failed to execute update request: %w", err)
	}

	return "**Voice recorded!**", nil
}

// GetPollResult получает результаты опроса из БД.
func (d *Database) GetPollResult(pollId string) (string, error) {
	reqGet := tarantool.NewSelectRequest(entities.PollsSpaceName).
		Index("primary").
		Iterator(tarantool.IterEq).
		Key([]interface{}{pollId})
	data, err := d.conn.Do(reqGet).Get()
	if err != nil {
		return "", fmt.Errorf("failed to execute select request: %w", err)
	}

	poll, err := parseData(data)
	if err != nil {
		return "", err
	}

	tbl := storage.PrintTable(poll)

	return tbl, nil
}

// ClosePoll закрывает опрос и обновляет данные в БД.
func (d *Database) ClosePoll(pollId, userId string) (string, error) {
	reqGet := tarantool.NewSelectRequest(entities.PollsSpaceName).
		Index("primary").
		Iterator(tarantool.IterEq).
		Key([]interface{}{pollId})
	data, err := d.conn.Do(reqGet).Get()
	if err != nil {
		return "", fmt.Errorf("failed to execute select request: %w", err)
	}

	poll, err := parseData(data)
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

	reqUpdate := tarantool.NewUpdateRequest(entities.PollsSpaceName).
		Key([]interface{}{pollId}).
		Operations(tarantool.NewOperations().
			Assign(5, poll.Closed))
	if _, err = d.conn.Do(reqUpdate).Get(); err != nil {
		return "", fmt.Errorf("failed to execute update request: %w", err)
	}

	return fmt.Sprintf("*Poll*: `%s` **has been successfully closed!**", pollId), nil
}

// DeletePoll удаляет опрос из БД.
func (d *Database) DeletePoll(pollId, userId string) (string, error) {
	reqGet := tarantool.NewSelectRequest(entities.PollsSpaceName).
		Index("primary").
		Iterator(tarantool.IterEq).
		Key([]interface{}{pollId})
	data, err := d.conn.Do(reqGet).Get()
	if err != nil {
		return "", fmt.Errorf("failed to execute select request: %w", err)
	}

	poll, err := parseData(data)
	if err != nil {
		return "", err
	}

	if poll.Creator != userId {
		return "", entities.NewUserError("**You don't have the permission to delete a vote!**")
	}

	reqUpd := tarantool.NewDeleteRequest(entities.PollsSpaceName).
		Key([]interface{}{pollId})
	if _, err = d.conn.Do(reqUpd).Get(); err != nil {
		return "", fmt.Errorf("failed to execute delete request: %w", err)
	}

	return fmt.Sprintf("*Poll*: `%s` **has been successfully delete!**", pollId), nil
}

func (d *Database) AddCmdToken(cmdPath, token string) error {
	tuple := []interface{}{cmdPath, token}
	reqPost := tarantool.NewInsertRequest(entities.TokensSpaceName).
		Tuple(tuple)

	if _, err := d.conn.Do(reqPost).Get(); err != nil {
		return fmt.Errorf("failed to execute upsert request: %w", err)
	}

	return nil
}

func (d *Database) ValidateCmdToken(cmdPath, token string) bool {
	reqGet := tarantool.NewSelectRequest(entities.TokensSpaceName).
		Index("primary").
		Iterator(tarantool.IterEq).
		Key([]interface{}{cmdPath})
	data, err := d.conn.Do(reqGet).Get()
	if err != nil {
		log.Printf("Error validating token: %v\n", err)
		return false
	}

	if len(data) == 0 {
		return false
	}

	tuple, ok := data[0].([]interface{})
	if !ok || len(tuple) < 2 {
		return false
	}

	validToken, ok := tuple[1].(string)
	if !ok {
		return false
	}

	return validToken == token
}
