package agentics

type State interface {
	GetMessages() []string
	AddMessages(messages []string)
}

type InputState struct {
	Messages []string
}

func (s *InputState) GetMessages() []string {
	return s.Messages
}

func (s *InputState) AddMessages(messages []string) {
	s.Messages = append(s.Messages, messages...)
}
