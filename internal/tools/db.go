package tools

import "fmt"

type dataMgr struct {
	data string
}

func (*dataMgr) Execute(data map[string]FunctionParameter) (string, error) {
	fmt.Println(data)
	return "data", nil
}

func NewDataMgr() *dataMgr {
	return &dataMgr{}
}

func init() {
	RegisterFunction(Tool{
		Type: "function",
		Function: Function{
			Name:        "get_current_weather",
			Description: "Get the current weather for a location",
			//Handler:     NewDataMgr().Execute,
			Parameters: FunctionParameterSpec{
				Type: "object",
				Properties: map[string]FunctionParameter{
					"location": FunctionParameter{
						Type:        "string",
						Description: "City and country, e.g. Paris, FR",
					},
					"format": FunctionParameter{
						Type:        "string",
						Description: "Units of measurement",
						Enum:        []string{"celsius", "fahrenheit"},
					},
				},
				Required: []string{"location", "format"},
			},
		},
	})
}
