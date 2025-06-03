# Asai — Personal AI Agent

**Asai** is a local or API-connected AI agent designed to coordinate secure access to tools such as a password manager, crypto wallets, and encrypted vector memory. It prioritizes **security**, **extensibility**, and **user privacy**.

The project is intended as a personal assistant, with optional support for access delegation in case the user becomes incapacitated — enabling secure transfer of knowledge and permissions to a trusted party.

---

## Features

* Connect to LLMs via API (supports local models via Ollama, GigaChat \[WIP], and more)
* Encrypted vector memory (in development)
* Integration with external tools (e.g., Bitwarden, crypto wallets)
* CLI interface + optional Telegram bot
* Persistent memory storage using PostgreSQL with `pgvector` support

---

## Project Structure

```bash
.
├── cmd/               # Entry points and interfaces
├── core/              # Core agent logic
├── llm/               # LLM API integrations (Ollama, GigaChat, etc.)
├── memory/            # Vector memory and context management
├── tools/             # External tools accessible by the agent (e.g., Bitwarden)
└── config/            # Configuration and environment setup
```

---

## Usage

Choose your preferred interface:

* CLI: `go run main.go cli`
* Telegram Bot: `go run main.go telegram`

---

## Supported LLMs

* Local models (via [Ollama](https://ollama.com) or LM Studio)
* GigaChat *(work in progress)*
* OpenAI (`https://api.openai.com`) *(work in progress)*

---

## Roadmap

* [ ] Role-based access and rights delegation
* [ ] Smart contract support
* [ ] External API calls via gRPC
* [ ] Custom tool plugin system
* [ ] Web-based UI
* [ ] Multi-user profiles and agent cooperation

