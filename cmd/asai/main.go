package main

import (
	"bufio"
	"fmt"
	"os"

	"asai/internal/core"
)

func main() {
	agent := core.NewAgent()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Asai: локальный ИИ-агент. Напиши команду:")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			continue
		}

		output, err := agent.Process(input)
		if err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}

		fmt.Println(output)
	}
}
