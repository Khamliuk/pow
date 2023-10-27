package entity

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Quit = iota
	RequestChallengeMessageType
	ResponseChallengeMessageType
	RequestResourceMessageType
	ResponseResourceMessageType
)

type Message struct {
	Type    int
	Payload string
}

func (m *Message) String() string {
	return fmt.Sprintf("%d|%s", m.Type, m.Payload)
}

// ParseMessage - parses Message from str, checks header and payload
func ParseMessage(str string) (*Message, error) {
	str = strings.TrimSpace(str)
	var msgType int
	// message has view as 1|payload (payload is optional)
	parts := strings.Split(str, "|")
	if len(parts) < 1 || len(parts) > 2 { //only 1 or 2 parts allowed
		return nil, fmt.Errorf("message doesn't match protocol")
	}
	// try to parse header
	msgType, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("cannot parse header")
	}
	msg := Message{
		Type: msgType,
	}
	// last part after | is payload
	if len(parts) == 2 {
		msg.Payload = parts[1]
	}
	return &msg, nil
}
