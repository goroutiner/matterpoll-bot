package database

import (
	"fmt"
	"matterpoll-bot/internal/entities"
)

// ParseData преобразовывает слайс интерфейсов к ожидаемым типам.
//   - `pollId`, `questions`, `creator` — строки.
//   - `options` и `voters` — карты, которые преобразуются с помощью вспомогательных функций.
//   - `closed` — булево значение.
func ParseData(data []interface{}) (*entities.Poll, error) {
	if len(data) == 0 {
		return nil, entities.NewUserError("**Invalid Poll_ID or not exists!**")
	}

	tuple, ok := data[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for data: %v", data)
	}
	if len(tuple) != 6 {
		return nil, fmt.Errorf("unexpected data format")
	}

	pollId, ok := tuple[0].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected type for pollId: %v", tuple[0])
	}

	questions, ok := tuple[1].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected type for questions: %v", tuple[1])
	}

	optionsRow, ok := tuple[2].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for options: %v", tuple[2])
	}
	options, err := convertMapInterfaceToStringInt(optionsRow)
	if err != nil {
		return nil, err
	}

	votersRow, ok := tuple[3].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for voters: %v", tuple[3])
	}
	voters, err := convertMapInterfaceToStringBool(votersRow)
	if err != nil {
		return nil, err
	}

	creator, ok := tuple[4].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected type for creator: %v", tuple[4])
	}

	closed, ok := tuple[5].(bool)
	if !ok {
		return nil, fmt.Errorf("unexpected type for closed: %v", tuple[5])
	}

	return &entities.Poll{PollId: pollId, Question: questions, Options: options, Voters: voters, Creator: creator, Closed: closed}, nil
}
