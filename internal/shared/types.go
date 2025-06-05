package shared

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MessageHistory struct {
	Messages []Message
	Limit    int
}

func (h *MessageHistory) WithNewUserInput(systemPrompt string, input string) []Message {
	msgs := trimToLimit(h.Messages, h.Limit)
	h.Messages = append(msgs, Message{
		Role:    "user",
		Content: input,
	})
	return append([]Message{{Role: "system", Content: systemPrompt}}, h.Messages...)
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
	return &MessageHistory{Messages: []Message{}, Limit: m.limit}
}

func (m *InMemoryContextManager) SaveContext(userID int64, ctx *MessageHistory) {
	m.data[userID] = ctx
}

func trimToLimit(messages []Message, limit int) []Message {
	var total int
	var trimmed []Message

	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		total += len(msg.Content)
		if total > limit {
			break
		}
		trimmed = append([]Message{msg}, trimmed...)
	}

	return trimmed
}
