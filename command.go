package main

import (
	"errors"
	"strconv"
)

type OP int

const (
	// string
	SET OP = iota
	INCR
	GET
	GETSET
)

type Command struct {
	OP        OP
	Key       string
	Arguments []string
}

func (cmd Command) Execute(state map[string]interface{}) (interface{}, error) {
	switch cmd.OP {
	case SET:
		state[cmd.Key] = cmd.Arguments[0]
	case INCR:
		if _, ok := state[cmd.Key]; !ok {
			state[cmd.Key] = "0"
		}
		i, err := strconv.ParseInt(state[cmd.Key].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		state[cmd.Key] = strconv.FormatInt(i+1, 10)
	case GET:
		return state[cmd.Key], nil
	case GETSET:
		if _, ok := state[cmd.Key].(string); !ok {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		oldValue := state[cmd.Key]
		state[cmd.Key] = cmd.Arguments[0]
		return oldValue, nil
	}

	return nil, nil
}

func NewCommand(op OP, key string, arguments ...string) Command {
	return Command{op, key, arguments}
}
