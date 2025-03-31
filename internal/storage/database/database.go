package database

import (
	"context"
	"fmt"
	"matterpoll-bot/config"
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/storage"
	"sync"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	_ "github.com/tarantool/go-tarantool/v2/datetime"
	_ "github.com/tarantool/go-tarantool/v2/decimal"
	_ "github.com/tarantool/go-tarantool/v2/uuid"
)

type Database struct {
	conn *tarantool.Connection
	mu   sync.RWMutex
}

// NewDatabaseConection возвращает структуру соединения с БД
func NewDatabaseStore(conn *tarantool.Connection) *Database {
	return &Database{conn: conn}
}

func NewDatabaseConection() (*tarantool.Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  config.DbSocket,
		User:     "user",
		Password: "secret",
	}
	opts := tarantool.Opts{
		Timeout: 5 * time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		return nil, fmt.Errorf("connection refused: %v", err)
	}

	return conn, err
}

func (d *Database) CreatePoll(poll *entities.Poll) error {
	tuple := []interface{}{
		poll.PollId,
		poll.Question,
		poll.Options,
		poll.Voters,
		poll.Creator,
		poll.Closed,
	}

	req := tarantool.NewInsertRequest(entities.SpaceName).Tuple(tuple)
	if _, err := d.conn.Do(req).Get(); err != nil {
		return err
	}

	return nil
}

func (d *Database) Vote(voice *entities.Voice) (string, error) {
	reqGet := tarantool.NewSelectRequest(entities.SpaceName).
		Index("poll_id_index").
		Limit(1).
		Iterator(tarantool.IterEq).
		Key([]interface{}{voice.PollId})
	data, err := d.conn.Do(reqGet).Get()
	if err != nil {
		return "", err
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	info, err := parseData(data)
	if err != nil {
		return "", err
	}

	if info.Closed {
		return "", entities.NewUserError(fmt.Sprintf("*Poll*: `%s` **is closed!**", voice.PollId))
	}

	if err := storage.ValidateVoice(info, voice); err != nil {
		return "", err
	}

	info.Options[voice.Option]++
	info.Voters[voice.UserId] = true

	reqUpd := tarantool.NewUpdateRequest(entities.SpaceName).
		Key([]interface{}{voice.PollId}).
		Operations(tarantool.NewOperations().
			Assign(2, info.Options).
			Assign(3, info.Voters))
	if _, err = d.conn.Do(reqUpd).Get(); err != nil {
		return "", err
	}

	return "**Voice recorded!**", nil
}

func (d *Database) GetPollResult(pollId string) (string, error) {
	reqGet := tarantool.NewSelectRequest(entities.SpaceName).
		Index("poll_id_index").
		Limit(1).
		Iterator(tarantool.IterEq).
		Key([]interface{}{pollId})
	data, err := d.conn.Do(reqGet).Get()
	if err != nil {
		return "", err
	}

	poll, err := parseData(data)
	if err != nil {
		return "", err
	}

	tbl := storage.PrintTable(poll)

	return tbl, nil
}

func (d *Database) ClosePoll(pollId, userId string) (string, error) {
	reqGet := tarantool.NewSelectRequest(entities.SpaceName).
		Index("poll_id_index").
		Limit(1).
		Iterator(tarantool.IterEq).
		Key([]interface{}{pollId})
	data, err := d.conn.Do(reqGet).Get()
	if err != nil {
		return "", err
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

	reqUpdate := tarantool.NewUpdateRequest(entities.SpaceName).
		Key([]interface{}{pollId}).
		Operations(tarantool.NewOperations().
			Assign(5, poll.Closed))
	if _, err = d.conn.Do(reqUpdate).Get(); err != nil {
		return "", err
	}

	return fmt.Sprintf("*Poll*: `%s` **has been successfully closed!**", pollId), nil
}

func (d *Database) DeletePoll(pollId, userId string) (string, error) {
	reqGet := tarantool.NewSelectRequest(entities.SpaceName).
		Index("poll_id_index").
		Limit(1).
		Iterator(tarantool.IterEq).
		Key([]interface{}{pollId})
	data, err := d.conn.Do(reqGet).Get()
	if err != nil {
		return "", err
	}

	poll, err := parseData(data)
	if err != nil {
		return "", err
	}

	if poll.Creator != userId {
		return "", entities.NewUserError("**You don't have the permission to delete a vote!**")
	}

	reqUpd := tarantool.NewDeleteRequest(entities.SpaceName).
		Key([]interface{}{pollId})
	if _, err = d.conn.Do(reqUpd).Get(); err != nil {
		return "", err
	}

	return fmt.Sprintf("*Poll*: `%s` **has been successfully delete!**", pollId), nil
}
