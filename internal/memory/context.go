package memory

import "asai/internal/llm"

type ContextManager interface {
	LoadContext(userID int64) *MessageHistory
	SaveContext(userID int64, ctx *MessageHistory)
}

type MessageHistory struct {
	Messages []llm.Message
	Limit    int
}

func (h *MessageHistory) WithNewUserInput(systemPrompt string, input string) []llm.Message {
	msgs := trimToLimit(h.Messages, h.Limit)
	h.Messages = append(msgs, llm.Message{
		Role:    "user",
		Content: input,
	})
	return append([]llm.Message{{Role: "system", Content: systemPrompt}}, h.Messages...)
}

type InMemoryContextManager struct {
	data  map[int64]*MessageHistory
	limit int
}

func NewInMemoryContextManager(limit int) *InMemoryContextManager {
	return &InMemoryContextManager{
		data:  make(map[int64]*MessageHistory),
		limit: limit,
	}
}

func (m *InMemoryContextManager) LoadContext(userID int64) *MessageHistory {
	if ctx, ok := m.data[userID]; ok {
		return ctx
	}
	return &MessageHistory{Messages: []llm.Message{}, Limit: m.limit}
}

func (m *InMemoryContextManager) SaveContext(userID int64, ctx *MessageHistory) {
	m.data[userID] = ctx
}

func trimToLimit(messages []llm.Message, limit int) []llm.Message {
	var total int
	var trimmed []llm.Message

	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		total += len(msg.Content)
		if total > limit {
			break
		}
		trimmed = append([]llm.Message{msg}, trimmed...)
	}

	return trimmed
}
