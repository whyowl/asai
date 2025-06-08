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

type llmModel interface {
	Generate(ctx context.Context, prompt []shared.Message, tools []tools.Function, userId int64) ([]shared.Message, error)
}

type embedModel interface {
	Embed(ctx context.Context, input string) ([]float32, error)
}

type Agent struct {
	llm          llmModel
	embed        embedModel
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
		embed:        llm.Providers[config.AppConfig.General.EmbedProvider],
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
	vector, err := a.embed.Embed(ctx, "test")
	if err != nil {
		return 0, err
	}
	return len(vector), nil
}

func (a *Agent) GetEmbed(ctx context.Context, n string) ([]float32, error) {
	embedText, err := a.embed.Embed(ctx, n)
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
