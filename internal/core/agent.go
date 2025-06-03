package core

import (
	"asai/internal/config"
	"asai/internal/memory"
	"log"
	"strings"
	"time"

	"asai/internal/llm"
	"asai/internal/tools"
)

type Agent struct {
	llm          llm.LLM
	memory       *memory.InMemoryContextManager
	systemPrompt string
}

func NewAgent(prompt string) *Agent {
	log.Printf("Created new Agent on %s provider (%s) with %d symbols context limit\n",
		config.AppConfig.General.LLMProvider,
		llm.Providers[config.AppConfig.General.LLMProvider],
		config.AppConfig.LLM.ContextLimit)
	return &Agent{
		llm:          llm.Providers[config.AppConfig.General.LLMProvider],
		memory:       memory.NewInMemoryContextManager(config.AppConfig.LLM.ContextLimit),
		systemPrompt: prompt,
	}

}

func (a *Agent) HandleInput(userID int64, input string) (string, error) {

	ctx := a.memory.LoadContext(userID)
	systemPrompt := buildSystemPrompt(a.systemPrompt, time.Now(), "Telegram", "")
	messages := ctx.WithNewUserInput(systemPrompt, input)

	response, err := a.llm.Generate(messages, tools.GetFunctionsForModel(), userID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	ctx.Messages = append(ctx.Messages, response...)
	a.memory.SaveContext(userID, ctx)
	return response[len(response)-1].Content, nil
}

func (a *Agent) GetDimensions() (int, error) {
	vector, err := a.llm.Embed("test")
	if err != nil {
		return 0, err
	}
	return len(vector), nil
}

func buildSystemPrompt(prompt string, data time.Time, mode string, userInfo string) string {
	tpl := strings.ReplaceAll(prompt, "{{TIME}}", data.String())
	tpl = strings.ReplaceAll(tpl, "{{MODE}}", mode)
	tpl = strings.ReplaceAll(tpl, "{{USER_INFO}}", userInfo)
	return tpl
}
