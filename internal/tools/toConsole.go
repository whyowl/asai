package tools

import "fmt"

func toConsoleEnter(data map[string]string) (string, error) {
	fmt.Println(data["text"])
	return "Text printed to console", nil
}

func init() {
	RegisterFunction(Tool{
		Type: "function",
		Function: Function{
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
		},
	})
}
