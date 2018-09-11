package chatserver

import (
	"fmt"
	"time"
)

type Message struct {
	SenderName string
	Text       string
	Time       time.Time
}

func (m Message) String() string {
	return fmt.Sprintf("[%s] %s: %s\n", m.Time.Format("15:04:05"), m.SenderName, m.Text)
}

func NewMessage(name, msg string) Message {
	message := Message{name, msg, time.Now()}
	return message
}
