package database_test

import (
	"context"
	"fmt"
	"log"
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/storage/database"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	ttC  testcontainers.Container // test Tarantool container
	conn *tarantool.Connection
	d    *database.Database
	poll = &entities.Poll{
		PollId:   "valid_id",
		Question: "test_question",
		Options:  map[string]int32{"opt1": 0, "opt2": 0},
		Voters:   map[string]bool{},
		Creator:  "creator_id",
		Closed:   false,
	}
)

// TestMain выполняет настройку окружения для интеграционных тестов.
func TestMain(m *testing.M) {
	var err error

	ctx := context.Background()
	buildContext, err := filepath.Abs("./docker")
	if err != nil {
		log.Fatalf("Failed to resolve absolute path: %v", err)
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: buildContext,
		},
		WaitingFor: wait.ForListeningPort("3301/tcp"),
	}

	ttC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	host, err := ttC.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	mappedPort, err := ttC.MappedPort(ctx, "3301")
	if err != nil {
		log.Fatal(err)
	}

	ttConf := &entities.TarantoolConfig{
		Address:  fmt.Sprintf("%s:%s", host, mappedPort.Port()),
		User:     "user",
		Password: "secret",
	}
	conn, err = database.NewDatabaseConection(ttConf)
	if err != nil {
		log.Fatal(err)
	}
	d = &database.Database{Conn: conn}

	code := m.Run()

	testcontainers.TerminateContainer(ttC)
	conn.CloseGraceful()
	os.Exit(code)
}

// TestCreatePoll проверяет успешное создание опроса в базе данных.
func TestCreatePoll(t *testing.T) {
	t.Cleanup(func() { truncateTable("polls", t) })
	err := d.CreatePoll(poll)
	require.NoError(t, err)

	actualPoll, err := getPoll(poll.PollId)
	require.NoError(t, err)
	require.Equal(t, poll, actualPoll)
}

// TestVote проверяет различные сценарии голосования.
func TestVote(t *testing.T) {
	t.Run("successful vote", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		voice := &entities.Voice{
			PollId: "valid_id",
			UserId: "user_id_1",
			Option: "opt1",
		}

		msg, err := d.Vote(voice)
		require.NoError(t, err)
		require.Equal(t, "**Voice recorded!**", msg)

		updatedPoll, err := getPoll(poll.PollId)
		require.NoError(t, err)
		require.Equal(t, int32(1), updatedPoll.Options["opt1"])
		require.True(t, updatedPoll.Voters["user_id_1"])
	})

	t.Run("invalid poll_id", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		voice := &entities.Voice{
			PollId: "invalid_id",
			UserId: "user_id_1",
			Option: "opt1",
		}

		msg, err := d.Vote(voice)
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
		require.Empty(t, msg)
	})

	t.Run("vote again", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		firstVoice := &entities.Voice{
			PollId: "valid_id",
			UserId: "user_id_2",
			Option: "opt1",
		}

		msg, err := d.Vote(firstVoice)
		require.NoError(t, err)
		require.NotEmpty(t, msg)

		repeatVoice := &entities.Voice{
			PollId: "valid_id",
			UserId: "user_id_2",
			Option: "opt1",
		}

		msg, err = d.Vote(repeatVoice)
		require.Error(t, err)
		require.Equal(t, "**You can't vote again!**", err.Error())
		require.Empty(t, msg)
	})

	t.Run("invalid option", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		voice := &entities.Voice{
			PollId: "valid_id",
			UserId: "user_id_3",
			Option: "invalid_opt",
		}

		msg, err := d.Vote(voice)
		require.Error(t, err)
		require.Equal(t, "**Invalid option!**", err.Error())
		require.Empty(t, msg)
	})

	t.Run("poll is already closed", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		msg, err := d.ClosePoll(poll.PollId, poll.Creator)
		require.NoError(t, err)
		require.NotEmpty(t, msg)

		voice := &entities.Voice{
			PollId: "valid_id",
			UserId: "user_id_4",
			Option: "opt1",
		}

		msg, err = d.Vote(voice)
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **is already closed!**", voice.PollId), err.Error())
		require.Empty(t, msg)
	})
}

// TestClosePoll проверяет различные сценарии закрытия голосования.
func TestClosePoll(t *testing.T) {
	t.Run("successful closed", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		pollId := "valid_id"
		userId := "creator_id"

		msg, err := d.ClosePoll(pollId, userId)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has been successfully closed!**", pollId), msg)

		updatedPoll, err := getPoll(poll.PollId)
		require.NoError(t, err)
		require.True(t, updatedPoll.Closed)
	})

	t.Run("invalid poll_id", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		pollId := "invalid_id"
		userId := "creator_id"

		msg, err := d.ClosePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
		require.Empty(t, msg)
	})

	t.Run("already closed poll", func(t *testing.T) {
		closedPoll := *poll
		closedPoll.Closed = true
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(&closedPoll, t)

		pollId := "valid_id"
		userId := "creator_id"

		msg, err := d.ClosePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **is already closed!**", pollId), err.Error())
		require.Empty(t, msg)
	})

	t.Run("don't have  permission", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		pollId := "valid_id"
		userId := "not_creator_id"

		msg, err := d.ClosePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, "**You don't have the permission to close a vote!**", err.Error())
		require.Empty(t, msg)
	})
}

// TestDeletePoll проверяет различные сценарии удаления голосования.
func TestDeletePoll(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		msg, err := d.DeletePoll(poll.PollId, "creator_id")
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("*Poll*: `%s` **has been successfully delete!**", poll.PollId), msg)

		deletedPoll, err := getPoll(poll.PollId)
		require.Error(t, err)
		require.Nil(t, deletedPoll)
	})

	t.Run("invalid poll_id", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		pollId := "invalid_id"
		userId := "creator_id"

		msg, err := d.DeletePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
		require.Empty(t, msg)
	})

	t.Run("don't have  permission", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		pollId := "valid_id"
		userId := "not_creator_id"

		msg, err := d.DeletePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, "**You don't have the permission to delete a vote!**", err.Error())
		require.Empty(t, msg)
	})

	t.Run("already deleted poll", func(t *testing.T) {
		t.Cleanup(func() { truncateTable("polls", t) })
		createTestPoll(poll, t)

		pollId := "valid_id"
		userId := "creator_id"

		_, err := d.DeletePoll(pollId, userId)
		require.NoError(t, err)

		msg, err := d.ClosePoll(pollId, userId)
		require.Error(t, err)
		require.Equal(t, "**Invalid Poll_ID or not exists!**", err.Error())
		require.Empty(t, msg)
	})
}

// TestAddCmdToken проверяет успешное добавления команды и токена в базу данных.
func TestAddCmdToken(t *testing.T) {
	t.Cleanup(func() { truncateTable("cmd_tokens", t) })

	cmdPath := "/test/command"
	token := "test_token"

	err := d.AddCmdToken(cmdPath, token)
	require.NoError(t, err)

	err = d.AddCmdToken(cmdPath, token)
	require.Error(t, err)
}

// TestValidateCmdToken проверяет токена команды.
func TestValidateCmdToken(t *testing.T) {
	t.Cleanup(func() { truncateTable("cmd_tokens", t) })

	cmdPath := "/test/command"
	token := "test_token"

	err := d.AddCmdToken(cmdPath, token)
	require.NoError(t, err)

	isValid := d.ValidateCmdToken(cmdPath, token)
	require.True(t, isValid)

	isInvalid := d.ValidateCmdToken(cmdPath, "invalid_token")
	require.False(t, isInvalid)
}

// createTestPoll создает тестовый опрос в базе данных.
func createTestPoll(poll *entities.Poll, t *testing.T) {
	err := d.CreatePoll(poll)
	require.NoError(t, err)
}

// getPoll получает опрос из базы данных по его идентификатору.
func getPoll(pollId string) (*entities.Poll, error) {
	reqGet := tarantool.NewSelectRequest(entities.PollsSpaceName).
		Index("primary").
		Iterator(tarantool.IterEq).
		Key([]interface{}{pollId})
	data, err := d.Conn.Do(reqGet).Get()
	if err != nil {
		return nil, fmt.Errorf("failed to execute select request: %w", err)
	}

	poll, err := database.ParseData(data)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

// truncateTable удаляет все записи из указанной таблицы в базе данных Tarantool.
func truncateTable(spaceName string, t *testing.T) {
	req := tarantool.NewCallRequest("box.space." + spaceName + ":truncate")
	_, err := d.Conn.Do(req).Get()
	require.NoError(t, err, "Failed to truncate table: %s", spaceName)
}
