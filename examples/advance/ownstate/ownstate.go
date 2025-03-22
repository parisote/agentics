package ownstate

type OwnState struct {
	Step     int
	Messages []string
}

func (s *OwnState) GetMessages() []string {
	return s.Messages
}

func (s *OwnState) AddMessages(messages []string) {
	s.Messages = append(s.Messages, messages...)
}
