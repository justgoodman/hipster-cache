package value_type

import (
	"fmt"
	"unsafe"
)

const (
	SetDictCmdName = "DSET"
	GetDictCmdName= "DGET"
	GetAllDictCmdName = "DGETALL"
)
type SetDictOperation struct {
	baseOperation
	key string
}

func NewSetDictOperation(key string) *SetDictOperation {
	return &SetDictOperation{baseOperation: baseOperation{commandName: SetDictCmdName}, key: key}
}

func (l *SetDictOperation) SetValue(sourceValue *interface{}, value interface{}) (valueSizeBytes int) {
	var stringValue string
	switch setValue :=  value.(type) {
		case string:
			stringValue = setValue
		default:
			l.err = fmt.Errorf("Error: list value in not the string")
			return
	}

	switch chain := (*sourceValue).(type) {
		case *Chain:
			element := chain.findElement(l.key)
			if element != nil {
				element.value = stringValue
				break
			}

			element = NewChainElement(l.key, stringValue)
			chain.addElement(element)
			chain.bytesSize += element.bytesSize

			valueSizeBytes = chain.bytesSize
		// It is like nil, default value
		case string:
			if chain != "" {
				l.err = fmt.Errorf("Error: list type is string")
			}
			element := NewChainElement(l.key, stringValue)
			newChain := NewChain(element)
			*sourceValue = interface{}(newChain)

			valueSizeBytes = int(unsafe.Sizeof(chain)) + element.bytesSize
		default:
			l.err = fmt.Errorf("Error: list type is not the Chain")
	}
	return
}

func (o *SetDictOperation) GetResult() (string,error) {
	if o.err == nil {
		return "OK",o.err }
	return "",o.err
}


type DictElement struct {
	key string
	value string
}

func NewDictElement(key, value string) *DictElement {
	return &DictElement{key: key, value: value}
}

type GetAllDictOperation struct {
	baseOperation
	Elements []*DictElement
}

func NewGetAllDictOperation() *GetAllDictOperation {
	return &GetAllDictOperation{baseOperation: baseOperation{commandName: GetAllDictCmdName}}
}

func (l *GetAllDictOperation) GetValue(value interface{}) {
	switch chain := value.(type) {
		case *Chain:
			l.Elements = chain.getAllValues()
		default:
			l.err = fmt.Errorf("Error: list type in not the Chain")
	}
}

func (l *GetAllDictOperation) GetResult() (string,error) {
	result := ""
	isFirst := true
	for _, element := range l.Elements {
		if isFirst {
			isFirst = false
		} else {
		result += "\n"
		}
		result +=  fmt.Sprintf(`"%s"`,element.key)
		result +=  "\n" + fmt.Sprintf(`"%s"`,element.value)
	}
	if len(l.Elements) == 0 {
		result = "(empty dictionary)"
	}
	return result, l.err
}

type GetDictOperation struct {
	baseOperation
	key string
	value string
}

func NewGetDictOperation(key string) *GetDictOperation {
	return &GetDictOperation{baseOperation: baseOperation{commandName: SetDictCmdName}, key: key,}
}

func (l *GetDictOperation) GetValue(value interface{}) {

	switch chain := value.(type) {
		case *Chain:
			element := chain.findElement(l.key)
			if element == nil {
				l.value = "(nil)"
			        return
			}
			l.value = fmt.Sprintf(`"%s"`,element.value)
		default:
			l.err = fmt.Errorf("Error: dict type in not the Chain")
	}
	return
}

func (l *GetDictOperation) GetResult() (string,error) {
	return l.value,l.err
}
