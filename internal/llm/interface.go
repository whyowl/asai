package llm

type Message struct {
	Role    string
	Content string
}

type LLM interface {
	Generate(prompt []Message) (string, error)
}
