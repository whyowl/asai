package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Telegram TelegramConfig
	Postgres PostgresConfig
	LLM      LLMConfig
	General  GeneralConfig
}

type TelegramConfig struct {
	Token     string
	WhiteList map[string]bool
}

type LLMConfig struct {
	Ollama       OllamaConfig
	GigaChat     GigaChatConfig
	ContextLimit int
}

type PostgresConfig struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}

type OllamaConfig struct {
	Model      string
	EmbedModel string
	Url        string
}

type GigaChatConfig struct {
	Secret      string
	Certificate string
	Model       string
	EmbedModel  string
	Scope       string
	ClientUrl   string
	TokenUrl    string
}

type GeneralConfig struct {
	LogLevel    int
	LLMProvider string
	UseTelegram bool
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	contextLimit, _ := strconv.Atoi(getEnv("LLM_CONTEXT_LIMIT", "12"))
	useTelegram, _ := strconv.ParseBool(getEnv("TELEGRAM_USE", "false"))

	AppConfig = &Config{
		Telegram: TelegramConfig{
			Token:     getEnv("TELEGRAM_TOKEN", ""),
			WhiteList: map[string]bool{getEnv("TELEGRAM_ID_OWNER", "0"): true},
		},
		Postgres: PostgresConfig{
			Host: getEnv("PG_HOST", "localhost"),
			Port: getEnv("PG_PORT", "5432"),
			User: getEnv("PG_USER", "postgres"),
			Pass: getEnv("PG_PASS", "password"),
			DB:   getEnv("PG_DB", "asai_db"),
		},
		LLM: LLMConfig{
			Ollama: OllamaConfig{
				Model:      getEnv("OLLAMA_MODEL", "llama3.1:8b"),
				EmbedModel: getEnv("OLLAMA_EMBED_MODEL", "nomic-embed-text:v1.5"),
				Url:        getEnv("OLLAMA_URI_BASE", "http://localhost:11434"),
			},
			GigaChat: GigaChatConfig{
				Secret:      getEnv("GIGACHAT_CLIENT_SECRET", ""),
				Certificate: getEnv("GIGACHAT_CLIENT_CERT", "russian_trusted_root_ca.cer"),
				Model:       getEnv("GIGACHAT_MODEL", "GigaChat"),
				EmbedModel:  getEnv("GIGACHAT_EMBED_MODEL", "Embeddings"),
				Scope:       getEnv("GIGACHAT_SCOPE", "GIGACHAT_API_PERS"),
				ClientUrl:   getEnv("GIGACHAT_CLIENT_URL", "https://gigachat.devices.sberbank.ru/"),
				TokenUrl:    getEnv("GIGACHAT_TOKEN_URL", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"),
			},
			ContextLimit: contextLimit,
		},
		General: GeneralConfig{
			LLMProvider: getEnv("LLM_PROVIDER", "ollama"),
			UseTelegram: useTelegram,
		},
	}
	log.Println("Config loaded")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
