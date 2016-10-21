package value_type

import (
	"fmt"
	"unsafe"
)

const (
	GetStringCmdName = "GET"
	SetStringCmdName = "SET"
)

type SetStringOperation struct {
	baseOperation
}

func NewSetStringOperation() *SetStringOperation {
	return &SetStringOperation{baseOperation{commandName: SetStringCmdName}}
}

func (o *SetStringOperation) GetResult() error {
	return o.err
}

func (o *SetStringOperation) SetValue(sourceValue, value interface{}) (valueSizeBytes int) {
	switch v := value.(type) {
	case string:
		valueSizeBytes = int(unsafe.Sizeof(sourceValue)) + len(v)
	default:
		o.err = fmt.Errorf(`Incorrect value type: type is not string`)
		return
	}
	sourceValue = value
	return
}

type GetStringOperation struct {
	baseOperation
	value string
}

func (o *GetStringOperation) GetResult() (string, error) {
	return o.value, o.err
}

func NewGetStringOperation() *GetStringOperation {
	return &GetStringOperation{baseOperation: baseOperation{commandName: GetStringCmdName}}
}

func (o *GetStringOperation) GetValue(value interface{}) {
	switch v := value.(type) {
	case string:
		o.value = v
	default:
		o.err = fmt.Errorf(`Incorrect value type: type is not string`)
	}
}
