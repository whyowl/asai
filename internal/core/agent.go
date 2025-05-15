package core

import (
	"strings"

	"asai/internal/llm"
	"asai/internal/tools"
)

type Agent struct {
	llm   llm.LLM
	tools map[string]Tool
}

type Tool interface {
	Execute(input string) (string, error)
}

func NewAgent() *Agent {
	return &Agent{
		llm: llm.NewLlamaClient(), // универсальный
		tools: map[string]Tool{
			"bitwarden": tools.NewBitwardenTool(),
		},
	}
}

func (a *Agent) Process(input string) (string, error) {
	// Если ключевое слово — bitwarden, вызываем инструмент напрямую
	if strings.Contains(strings.ToLower(input), "bitwarden") {
		return a.tools["bitwarden"].Execute(input)
	}

	// Иначе отправляем в LLaMA
	response, err := a.llm.Generate(input)
	if err != nil {
		return "", err
	}

	return response, nil
}
