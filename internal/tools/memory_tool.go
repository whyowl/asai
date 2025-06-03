package tools

import "fmt"

type dataMgr struct {
	data string
}

func (*dataMgr) Execute(data map[string]string, userID int64) (string, error) {
	switch data["method"] {
	case "write":
		fmt.Println("Была вызвана запись в память", data)
		return "Была вызвана запись в память", nil
	case "search":
		fmt.Println("Был вызван поиск в памяти", data)
		return "Был вызван поиск в памяти", nil
	default:
		fmt.Println("Не верный параметр", data)
		return "Не верный параметр", nil
	}
}

func NewDataMgr() *dataMgr {
	return &dataMgr{}
}

func init() {
	RegisterFunction(Function{
		Name:        "memory",
		Description: "The function for interacting with memory.  It is possible to write to memory or search in memory on request. The database is vector-based. To search, you need to formulate a semantic query.",
		Handler:     NewDataMgr().Execute,
		Parameters: FunctionParameterSpec{
			Type: "object",
			Properties: map[string]FunctionParameter{
				"data": FunctionParameter{
					Type:        "string",
					Description: "The text that needs to be stored in memory or searched for.",
				},
				"method": FunctionParameter{
					Type:        "string",
					Description: "The method for interacting with memory is writing or searching with reading",
					Enum:        []string{"write", "search"},
				},
			},
			Required: []string{"data", "method"},
		},
	})
}
