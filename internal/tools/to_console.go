package tools

import (
	"context"
	"fmt"
)

func toConsoleEnter(ctx context.Context, data map[string]string, userID int64) (string, error) {
	fmt.Println(userID, data["text"])
	return fmt.Sprintf("'%s' was printed to console", data["text"]), nil
}

func init() {
	RegisterFunction(Function{
		Name:        "print_to_console",
		Description: "Print text to console on server where Asai run. Dont use this for chat with user, thats just tool for logs",
		Handler:     toConsoleEnter,
		Parameters: FunctionParameterSpec{
			Type: "object",
			Properties: map[string]FunctionParameter{
				"text": FunctionParameter{
					Type:        "string",
					Description: "text that must be print",
				},
			},
			Required: []string{"text"},
		},
	})
}
