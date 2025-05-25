package tools

type dataMgr struct {
	data string
}

func (*dataMgr) Execute(data map[string]string) (string, error) {
	return "Эта заглушка, сервис погоды пока не работает. Предупредите пользователя, что временно не можешь дать данные", nil
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
			Handler:     NewDataMgr().Execute,
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
