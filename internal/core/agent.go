package core

import (
	"asai/internal/config"
	"asai/internal/shared"
	"context"
	"log"
	"strings"
	"time"

	"asai/internal/llm"
	"asai/internal/tools"
)

type ContextManager interface {
	LoadContext(userID int64) *shared.MessageHistory
	SaveContext(userID int64, ctx *shared.MessageHistory)
}

type Agent struct {
	llm          llm.LLM
	memory       ContextManager
	systemPrompt string
}

func NewAgent(prompt string) *Agent {
	log.Printf("Created new Agent on %s provider (%s) with %d symbols context limit\n",
		config.AppConfig.General.LLMProvider,
		llm.Providers[config.AppConfig.General.LLMProvider],
		config.AppConfig.LLM.ContextLimit)
	return &Agent{
		llm:          llm.Providers[config.AppConfig.General.LLMProvider],
		memory:       shared.NewInMemoryContextManager(config.AppConfig.LLM.ContextLimit),
		systemPrompt: prompt,
	}

}

func (a *Agent) HandleInput(ctx context.Context, userID int64, input string) (string, error) {

	messsageContext := a.memory.LoadContext(userID)
	systemPrompt := buildSystemPrompt(a.systemPrompt, time.Now(), "Telegram", "")
	messages := messsageContext.WithNewUserInput(systemPrompt, input)

	response, err := a.llm.Generate(ctx, messages, tools.GetFunctionsForModel(), userID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	messsageContext.Messages = append(messsageContext.Messages, response...)
	a.memory.SaveContext(userID, messsageContext)
	return response[len(response)-1].Content, nil
}

func (a *Agent) GetDimensions(ctx context.Context) (int, error) {
	vector, err := a.llm.Embed(ctx, "test")
	if err != nil {
		return 0, err
	}
	return len(vector), nil
}

func (a *Agent) GetEmbed(ctx context.Context, n string) ([]float32, error) {
	embedText, err := a.llm.Embed(ctx, n)
	if err != nil {
		return []float32{}, err
	}
	return embedText, nil
}

func buildSystemPrompt(prompt string, data time.Time, mode string, userInfo string) string {
	tpl := strings.ReplaceAll(prompt, "{{TIME}}", data.String())
	tpl = strings.ReplaceAll(tpl, "{{MODE}}", mode)
	tpl = strings.ReplaceAll(tpl, "{{USER_INFO}}", userInfo)
	return tpl
}
