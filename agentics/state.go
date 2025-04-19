package agentics

type State interface {
	GetMessages() []string
	AddMessages(messages []string)
	GetAllAtts() map[string]interface{}
	GetAtt(name string) interface{}
}

type InputState struct {
	Att      map[string]interface{}
	Messages []string
}

func (s *InputState) GetMessages() []string {
	return s.Messages
}

func (s *InputState) AddMessages(messages []string) {
	s.Messages = append(s.Messages, messages...)
}

func (s *InputState) GetAtt(name string) interface{} {
	return s.Att[name]
}

func (s *InputState) GetAllAtts() map[string]interface{} {
	return s.Att
}
