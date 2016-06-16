package main

import (
	"errors"
	"strconv"
	"time"
)

type OP int

const (
	// string
	SET OP = iota
	INCR
	GET
	GETSET
	EXPIRE
)

type Command struct {
	OP        OP
	Key       string
	Arguments []string
	TX        *Transaction
}

func (cmd Command) Execute(state State) (interface{}, error) {
	switch cmd.OP {
	case SET:
		state[cmd.Key] = &Value{Val: cmd.Arguments[0]}
		return "OK", nil
	case INCR:
		if _, ok := state[cmd.Key]; !ok {
			state[cmd.Key] = &Value{Val: "0"}
		}
		i, err := strconv.ParseInt(state[cmd.Key].Val.(string), 10, 64)
		if err != nil {
			return nil, err
		}
		//newValue := Value{Val: strconv.FormatInt(i+1, 10), Expire: state[cmd.Key].Expire, WillExpire: state[cmd.Key].WillExpire}
		//state[cmd.Key] = newValue
		state[cmd.Key].UpdateVal(strconv.FormatInt(i+1, 10))

		return state[cmd.Key].Val, nil
	case GET:
		return state[cmd.Key].Val, nil
	case GETSET:
		if _, ok := state[cmd.Key].Val.(string); !ok {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		oldValue := state[cmd.Key].Val
		state[cmd.Key] = &Value{Val: cmd.Arguments[0]}
		return oldValue, nil
	case EXPIRE:
		seconds, err := strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return nil, err
		}
		state[cmd.Key].UpdateExpire(cmd.TX.Header.Time.Add(time.Duration(seconds) * time.Second))

		return "OK", nil
	}

	return nil, nil
}

func NewCommand(op OP, key string, arguments ...string) Command {
	return Command{op, key, arguments, nil}
}
