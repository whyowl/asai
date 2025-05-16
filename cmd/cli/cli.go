package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"asai/internal/core"
)

func Run(ctx context.Context, a *core.Agent) {
	reader := bufio.NewScanner(os.Stdin)
	chatID := int64(0) // CLI = один пользователь

	fmt.Println("Asai CLI — введите текст:")

	for reader.Scan() {
		input := reader.Text()
		if input == "exit" {
			break
		}

		resp, err := a.HandleInput(chatID, input)
		if err != nil {
			fmt.Println("⚠️ Ошибка:", err)
			continue
		}

		fmt.Println("🤖:", resp)
	}
}
