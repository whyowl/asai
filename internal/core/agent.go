package core

import (
	"asai/internal/memory"
	"fmt"
	"os"
	"strings"
	"time"

	"asai/internal/llm"
	"asai/internal/tools"
)

type Agent struct {
	llm    llm.LLM
	tools  map[string]Tool
	memory *memory.InMemoryContextManager
}

type Tool interface {
	Execute(input string) (string, error)
}

func NewAgent() *Agent {
	return &Agent{
		llm: llm.NewLlamaClient(os.Getenv("ASAI_LLM_URI_BASE"), os.Getenv("ASAI_LLM_MODEL")),
		tools: map[string]Tool{
			"bitwarden": tools.NewBitwardenTool(),
			"savedata":  tools.NewDataMgr(),
		},
		memory: memory.NewInMemoryContextManager(),
	}
}

func (a *Agent) HandleInput(userID int64, input string) (string, error) {

	ctx := a.memory.LoadContext(userID)
	systemPrompt := buildSystemPrompt(a.tools, time.Now(), "Telegram")
	messages := ctx.WithNewUserInput(systemPrompt, input)

	response, err := a.llm.Generate(messages)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	ctx.Messages = append(ctx.Messages, response)
	a.memory.SaveContext(userID, ctx)
	return response.Content, nil
}

func buildSystemPrompt(tools map[string]Tool, data time.Time, mode string) string {
	const SYSTEM_PROMPT = `
Ты — Asai, персональный ИИ-агент. Ты работаешь напрямую на одного пользователя и строго соблюдаешь приватность.

Всегда выбирай, когда уместно вызвать инструмент, а когда ответить сам. Если не уверен — уточни.
Никогда не выдумывай данные. Не сохраняй ничего без указания пользователя.
Ты не человек и не изображаешь его. Ты — приватный помощник.

Текущий режим работы: {{MODE}}
Дата и время: {{TIME}}
`
	tpl := strings.ReplaceAll(SYSTEM_PROMPT, "{{TIME}}", data.String())
	tpl = strings.ReplaceAll(tpl, "{{MODE}}", mode)
	return tpl
}

