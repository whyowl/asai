package tools

import "fmt"

type BitwardenTool struct {
}

func NewBitwardenTool() *BitwardenTool {
	return &BitwardenTool{}
}

func (b *BitwardenTool) Execute(input string) (string, error) {
	return fmt.Sprintf("Bitwarden получил запрос: %s", input), nil
}
