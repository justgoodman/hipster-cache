package tcp

import (
	"fmt"
	"strings"
)

type ClientMessage struct {
	command string
	params  []string
}

func NewClientMessage() *ClientMessage {
	return &ClientMessage{}
}

func (m *ClientMessage) Init(value string) error {
	words := m.splitMessageBySpaceAndQuates(value)
	fmt.Printf(`\n Words :"%#v"`, words)
	if len(words) == 0 {
		return fmt.Errorf(`Error: you don't set the command`)
	}
	m.command = strings.ToUpper(words[0])
	m.params = words[1:]
	return nil
}

func (m *ClientMessage) splitMessageBySpaceAndQuates(message string) []string {
	words := []string{}
	var word string
	var character string
	delimeter := ""
	for _, characterCode := range message {
		character = string(characterCode)
		switch character {
		case ` `:
			if delimeter == "" {
				words = append(words, word)
				word = ""
				break
			}
			word += character
		case `"`:
			if delimeter == character {
				delimeter = ""
				break
			}
			if delimeter == "" {
				delimeter = character
				break
			}
		case "\n", "\r":
		default:
			word += character

		}
	}
	if word != "" {
		words = append(words, word)
	}
	return words
}
