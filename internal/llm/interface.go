package llm

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLM interface {
	Generate(prompt []Message) (string, error)
}
