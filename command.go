// LWW register
package main

type OP int

const (
	// string
	SET OP = iota
)

type Command struct {
	OP        OP
	Key       string
	Arguments []string
}

func NewCommand(op OP, key string, arguments ...string) Command {
	return Command{op, key, arguments}
}
