package tools

import "fmt"

type dataMgr struct {
	data string
}

func (*dataMgr) Execute(data string) (string, error) {
	fmt.Println(data)
	return data, nil
}

func NewDataMgr() *dataMgr {
	return &dataMgr{}
}
