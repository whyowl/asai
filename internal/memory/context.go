package memory

import "asai/internal/llm"

type ContextManager interface {
	LoadContext(userID int64) *MessageHistory
	SaveContext(userID string, ctx *MessageHistory)
}

type MessageHistory struct {
	Messages []llm.Message
	Limit    int
}

func (h *MessageHistory) WithNewUserInput(systemPrompt string, input string) []llm.Message {
	//msgs := trimToLimit(h.Messages, h.Limit)
	h.Messages = append(h.Messages, llm.Message{
		Role:    "user",
		Content: input,
	})
	msgs := h.Messages
	return append([]llm.Message{{Role: "system", Content: systemPrompt}}, msgs...)
}

type InMemoryContextManager struct {
	data map[int64]*MessageHistory
}

func NewInMemoryContextManager() *InMemoryContextManager {
	return &InMemoryContextManager{
		data: make(map[int64]*MessageHistory),
	}
}

func (m *InMemoryContextManager) LoadContext(userID int64) *MessageHistory {
	if ctx, ok := m.data[userID]; ok {
		return ctx
	}
	return &MessageHistory{Messages: []llm.Message{}, Limit: 20}
}

func (m *InMemoryContextManager) SaveContext(userID int64, ctx *MessageHistory) {
	m.data[userID] = ctx
}
