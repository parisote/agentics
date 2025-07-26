package agentics

import (
	"sync"
)

type Message struct {
	Role       string // "system", "user", "assistant", "tool"
	Content    string
	ToolCallID string
	ToolCalls  []ToolCall // For assistant messages with tool calls
}

type Memory interface {
	Add(role string, content string, toolCallID ...string)
	LastN(n int) []Message
	All() []Message
	Len() int
	ToArrayString() []string
	AddBytes(string, []byte)
}

type SliceMemory struct {
	mu   sync.RWMutex
	max  int
	data []Message
}

func NewSliceMemory(max int) *SliceMemory {
	return &SliceMemory{
		max:  max,
		data: []Message{},
	}
}

func (m *SliceMemory) Add(role string, content string, toolCallID ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(toolCallID) > 0 {
		m.data = append(m.data, Message{
			Role:       role,
			Content:    content,
			ToolCallID: toolCallID[0],
		})
	} else {
		m.data = append(m.data, Message{
			Role:    role,
			Content: content,
		})
	}

	if len(m.data) > m.max {
		m.data = m.data[len(m.data)-m.max:]
	}
}

func (m *SliceMemory) AddMessage(message Message) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = append(m.data, message)

	if len(m.data) > m.max {
		m.data = m.data[len(m.data)-m.max:]
	}
}

func (m *SliceMemory) LastN(n int) []Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if n > len(m.data) {
		return m.data
	}

	return m.data[len(m.data)-n:]
}

func (m *SliceMemory) All() []Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data
}

func (m *SliceMemory) ToArrayString() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var arr []string
	for _, message := range m.data {
		arr = append(arr, message.Content)
	}

	return arr
}

func (m *SliceMemory) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.data)
}

func (m *SliceMemory) AddBytes(role string, content []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.data) > m.max {
		m.data = m.data[len(m.data)-m.max:]
	}
}
