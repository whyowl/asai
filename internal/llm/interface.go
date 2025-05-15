package llm

type LLM interface {
	Generate(prompt string) (string, error)
}
