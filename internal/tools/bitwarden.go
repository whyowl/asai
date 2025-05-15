package tools

import "fmt"

type BitwardenTool struct {
	// тут можно хранить токены, пути к bw CLI
}

func NewBitwardenTool() *BitwardenTool {
	return &BitwardenTool{}
}

func (b *BitwardenTool) Execute(input string) (string, error) {
	// тут будет разбор и вызов bw CLI — пока заглушка
	return fmt.Sprintf("Bitwarden получил запрос: %s", input), nil
}
