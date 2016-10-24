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

func (o *SetStringOperation) GetResult() (string,error) {
	if o.err == nil {
		return "OK",o.err }
	return "",o.err
}

func (o *SetStringOperation) SetValue2(sourceValue *string, value *string) {
	fmt.Printf("\n SetValue inner:%#v",value)
	*sourceValue = *value
	fmt.Printf("\n New value:%#v", sourceValue)
	return
}


func (o *SetStringOperation) SetValue(sourceValue *interface{}, value interface{}) (valueSizeBytes int) {
	fmt.Printf("\n SetValue inner:%#v",value)
	switch v := value.(type) {
	case string:
		valueSizeBytes = int(unsafe.Sizeof(sourceValue)) + len(v)
	default:
		o.err = fmt.Errorf(`Incorrect value type: type is not string`)
		return
	}
	*sourceValue = value
	fmt.Printf("\n New value:%#v", sourceValue)
	return
}

type GetStringOperation struct {
	baseOperation
	value string
	isHit bool
}

func (o *GetStringOperation) GetResult() (string, error) {
	if !o.isHit {
		return "(nil)", o.err
	}
	return fmt.Sprintf(`"%s"`, o.value), o.err
}

func NewGetStringOperation() *GetStringOperation {
	return &GetStringOperation{baseOperation: baseOperation{commandName: GetStringCmdName}}
}

func (o *GetStringOperation) GetValue(value interface{}) {
	switch v := value.(type) {
	case string:
		o.value = v
		o.isHit = true
	default:
		o.err = fmt.Errorf(`Incorrect value type: type is not string`)
	}
}
