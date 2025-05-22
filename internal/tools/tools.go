package tools

var functionRegistry = map[string]Tool{}

type Tool struct {
	Type     string   `json:"type"` // "function"
	Function Function `json:"function"`
}

type Function struct {
	Name        string                                  `json:"name"`
	Description string                                  `json:"description"`
	Handler     func(map[string]string) (string, error) `json:"-"`
	Parameters  FunctionParameterSpec                   `json:"parameters"`
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

func RegisterFunction(tool Tool) {
	functionRegistry[tool.Function.Name] = tool
}

func GetToolsForModel() []Tool {
	var tools []Tool
	for _, tool := range functionRegistry {
		tools = append(tools, tool)
	}
	return tools
}

func CallFunctionsByModel(name string, arg map[string]string) (string, error) {
	response, err := functionRegistry[name].Function.Handler(arg)
	if err != nil {
		return "", err
	}
	return response, nil
}
