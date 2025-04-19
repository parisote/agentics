package ownstate

type OwnState struct {
	Att      map[string]interface{}
	Step     int
	Messages []string
}

func (s *OwnState) GetMessages() []string {
	return s.Messages
}

func (s *OwnState) AddMessages(messages []string) {
	s.Messages = append(s.Messages, messages...)
}

func (s *OwnState) GetAtt(name string) interface{} {
	return s.Att[name]
}

func (s *OwnState) GetAllAtts() map[string]interface{} {
	return s.Att
}
