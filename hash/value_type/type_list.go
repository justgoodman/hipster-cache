package value_type

import (
	"fmt"
	"unsafe"
)

type ListPushOperation struct {
	elementIndex int
}

func (this *ListPushOperation) SetValue(sourceValue, value interface{}) (valueSizeBytes int) {
	var stringValue string
	switch setVal :=  value.(type) {
		case string:
			stringValue = setValue
		default:
			this.Err = fmt.Sprintf("Error: list value in not the string")
			return
	}
	switch chain := sourceValue.(type) {
		case *Chain:
			element := NewChainElement(stringValue)
			chain.addElement(element)
			chain.sizeBytes += elem.sizeBytes
			this.elementIndex = chain.lenght - 1

			valueSizeBytes = chain.sizeBytes
		case nil:
			element := NewChainElement(stringValue)
			chain = NewChain(element)
			this.elementIndex = 0

			valueSizeBytes := unsafe.Sizeof(chain) + element.sizeBytes
		default:
			this.Err = fmt.Sprintf("Error: list type in not the Chain")
	}
}

type ListLenghtOperation struct {
	Length int
	Err error
}

func (this *ListLenghtOperation) GetValue(value interface{}) {
	switch chain := value.(type) {
		case *Chain:
			this.Lenght = chain.Lenght
		default:
			this.Err = fmt.Errorf("Error: list type in not the Chain")
	}
}

type ListRangeOperation struct {
	IndexStart int
	IndexEnd int
	Values []string
	Err error
}

func (this *ListLenghtOperation) GetValue(value interface{}) {
	switch chain := value.(type) {
		case *Chain:
			this.Values = chain.GetRangeValues(this.IndexStart, this.IndexEnd)
		default:
			this.Err = fmt.Errorf("Error: list type in not the Chain")
	}
}

type ListSetOperation struct {
	Index int
	Err string
}

func (this *ListSetOperation) SetValue(sourceValue,value interface{}) (valueSizeBytes int) {
	var stringValue string
	switch setVal :=  value.(type) {
		case string:
			stringValue = setValue
		default:
			this.Err = fmt.Sprintf("Error: list value in not the string")
			return
	}
	switch chain := sourceValue.(type) {
		case *Chain:
			element = chain.FindElement(this.Index)
			if element == nil {
				this.Err = fmt.Errorf(`Index "%d" not found in List`)
			        return
			}
			deltaBytes := len(stringValue) - len(element.value)
			element.value := stringValue
			chain.byteSize -=  deltaBytes

			valueSizeBytes = chain.byteSize
		default:
			this.Err = fmt.Errorf("Error: list type in not the Chain")
	}
}
