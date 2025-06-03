package tools

type weather struct {
	data string
}

func (*weather) Execute(data map[string]string, userID int64) (string, error) {
	return "Эта заглушка, сервис погоды пока не работает. Предупредите пользователя, что временно не можешь дать данные", nil
}

func NewWeather() *weather {
	return &weather{}
}

func init() {
	RegisterFunction(Function{
		Name:        "get_current_weather",
		Description: "Get the current weather for a location",
		Handler:     NewWeather().Execute,
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
	})
}
