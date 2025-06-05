package tools

import "context"

var functionRegistry = map[string]Function{}

type Function struct {
	Name        string                                                          `json:"name"`
	Description string                                                          `json:"description"`
	Handler     func(context.Context, map[string]string, int64) (string, error) `json:"-"`
	Parameters  FunctionParameterSpec                                           `json:"parameters"`
}

type FunctionParameterSpec struct {
	Type       string                       `json:"type"` // "object"
	Properties map[string]FunctionParameter `json:"properties"`
	Required   []string                     `json:"required"`
}

type FunctionParameter struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type FunctionCall struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments"`
}

func RegisterFunction(f Function) {
	functionRegistry[f.Name] = f
}

func GetFunctionsForModel() []Function {
	var functions []Function
	for _, f := range functionRegistry {
		functions = append(functions, f)
	}
	return functions
}

func CallFunctionsByModel(ctx context.Context, name string, arg map[string]string, userID int64) (string, error) {
	response, err := functionRegistry[name].Handler(ctx, arg, userID)
	if err != nil {
		return "", err
	}
	return response, nil
}
