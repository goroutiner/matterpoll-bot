package database

import (
	"fmt"
	"matterpoll-bot/internal/entities"
)

func parseData(data []interface{}) (*entities.Poll, error) {
	if len(data) == 0 {
		return nil, entities.NewUserError("**Invalid Poll_ID or not exists!**")
	}

	row, ok := data[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for data: %v", data)
	}
	if len(row) != 6 {
		return nil, fmt.Errorf("unexpected data format")
	}

	pollId, ok := row[0].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected type for pollId: %v", row[0])
	}

	questions, ok := row[1].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected type for questions: %v", row[1])
	}

	optionsRow, ok := row[2].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for options: %v", row[2])
	}
	options, err := convertMapInterfaceToStringInt(optionsRow)
	if err != nil {
		return nil, err
	}

	votersRow, ok := row[3].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for voters: %v", row[3])
	}
	voters, err := convertMapInterfaceToStringBool(votersRow)
	if err != nil {
		return nil, err
	}

	creator, ok := row[4].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected type for creator: %v", row[4])
	}

	closed, ok := row[5].(bool)
	if !ok {
		return nil, fmt.Errorf("unexpected type for closed: %v", row[5])
	}

	return &entities.Poll{PollId: pollId, Question: questions, Options: options, Voters: voters, Creator: creator, Closed: closed}, nil
}
