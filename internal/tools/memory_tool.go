package tools

import (
	"asai/internal/memory"
	"context"
	"fmt"
	"strconv"
)

type embedFunction interface {
	GetEmbed(string) ([]float32, error)
}

type dataMgr struct {
	embedFunc embedFunction
}

func (d *dataMgr) Execute(data map[string]string, userID int64) (string, error) {
	switch data["method"] {
	case "write":
		ctx := context.Background()
		embedText, err := d.embedFunc.GetEmbed(data["data"])
		if err != nil {
			return "", err
		}
		err = memory.InsertEmbedding(ctx, memory.DB, "memory", strconv.FormatInt(userID, 10), data["data"], embedText)
		if err != nil {
			return "", err
		}
		fmt.Println("Была вызвана запись в память", data)
		return fmt.Sprintf("Была вызвана запись в память '%s'", data), nil
	case "search":
		ctx := context.Background()
		embedText, err := d.embedFunc.GetEmbed(data["data"])
		if err != nil {
			return "", err
		}
		response, err := memory.QuerySimilarEmbeddings(ctx, memory.DB, "memory", strconv.FormatInt(userID, 10), embedText, 5)
		if err != nil {
			return "", err
		}
		var result string
		for i, s := range response {
			result += fmt.Sprintf("record %d: %s \n", i, s)
		}
		fmt.Println("Был вызван поиск в памяти", data["data"], response)
		return result, nil
	default:
		fmt.Println("Не верный параметр", data)
		return "Не верный параметр", nil
	}
}

func NewDataMgr(e embedFunction) *dataMgr {
	return &dataMgr{embedFunc: e}
}

func InitDataMgr(e embedFunction) {
	RegisterFunction(Function{
		Name:        "memory",
		Description: "The function for interacting with memory. It is possible to write to memory or search in memory on request. The database is vector-based. To search, you need to formulate a semantic query. Write down information about the user in the third person on Russian and search same. For example: Пользователя зовут Максим.",
		Handler:     NewDataMgr(e).Execute,
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
