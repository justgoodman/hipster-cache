package value_type

import (
	"fmt"
	"unsafe"
	"strconv"
)

const (
	PushListCmdName = "LPUSH"
	RangeListCmdName = "LRANGE"
	SetListCmdName = "LSET"
	LenghtListCmdName = "LLEN"
)

type PushListOperation struct {
	baseOperation
	elementIndex int
}

func NewPushListOperation() *PushListOperation {
	return &PushListOperation{baseOperation: baseOperation{commandName: PushListCmdName}}
}

func (l *PushListOperation) SetValue(sourceValue *interface{}, value interface{}) (valueSizeBytes int) { var stringValue string
	switch setValue :=  value.(type) {
		case string:
			stringValue = setValue
		default:
			l.err = fmt.Errorf("Error: list value in not the string")
			return
	}

	element := NewChainElement(stringValue)
	switch chain := (*sourceValue).(type) {
		case *Chain:
			chain.addElement(element)
			chain.bytesSize += element.bytesSize
			l.elementIndex = chain.countElements - 1

			valueSizeBytes = chain.bytesSize
		// It is like nil, default value
		case string:
			if chain != "" {
				l.err = fmt.Errorf("Error: list type is string")
			}
			newChain := NewChain(element)
			l.elementIndex = 0
			*sourceValue = interface{}(newChain)

			valueSizeBytes = int(unsafe.Sizeof(chain)) + element.bytesSize
		default:
			l.err = fmt.Errorf("Error: list type is not the Chain")
	}
	return
}

func (o *PushListOperation) GetResult() (string,error) {
	if o.err == nil {
		return "OK",o.err }
	return "",o.err
}



type LenghtListOperation struct {
	baseOperation
	Lenght int
}

func NewLenghtListOperation() *LenghtListOperation {
	return  &LenghtListOperation{baseOperation: baseOperation{commandName: LenghtListCmdName}}
}

func (l *LenghtListOperation) GetResult() (string,error) {
	return strconv.Itoa(l.Lenght), l.err
}

func (l *LenghtListOperation) GetValue(value interface{}) {
	switch chain := value.(type) {
		case *Chain:
			l.Lenght = chain.countElements
		default:
			l.err = fmt.Errorf("Error: list type in not the Chain")
	}
}

type RangeListOperation struct {
	baseOperation
	indexStart int
	indexEnd int
	Values []string
}

func NewRangeListOperation(indexStart,indexEnd int) *RangeListOperation {
	return &RangeListOperation{baseOperation: baseOperation{commandName: RangeListCmdName},indexStart: indexStart, indexEnd: indexEnd,}
}

func (l *RangeListOperation) GetValue(value interface{}) {
	switch chain := value.(type) {
		case *Chain:
			l.Values = chain.GetRangeValues(l.indexStart, l.indexEnd)
		default:
			l.err = fmt.Errorf("Error: list type in not the Chain")
	}
}

func (l *RangeListOperation) GetResult() (string,error) {
	result := ""
	isFirst := true
	for _, value := range l.Values {
		if isFirst {
			isFirst = false
		} else {
		result += "\n"
		}
		result +=  fmt.Sprintf(`"%s"`,value)
	}
	if len(l.Values) == 0 {
		result = "(empty list)"
	}
	return result, l.err
}

type SetListOperation struct {
	baseOperation
	index int
}

func NewSetListOperation(index int) *SetListOperation {
	return &SetListOperation{baseOperation: baseOperation{commandName: SetListCmdName}, index: index,}
}

func (l *SetListOperation) SetValue(sourceValue *interface{},value interface{}) (valueSizeBytes int) {
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
			element := chain.findElement(l.index)
			if element == nil {
				l.err = fmt.Errorf(`Error: Index "%d" not found in List`, l.index)
			        return
			}
			deltaBytes := len(stringValue) - len(element.value)
			element.value = stringValue
			chain.bytesSize -=  deltaBytes

			valueSizeBytes = chain.bytesSize
		default:
			l.err = fmt.Errorf("Error: list type in not the Chain")
	}
	return
}

func (l *SetListOperation) GetResult() (string,error) {
	if l.err == nil {
		return "OK",l.err }
	return "",l.err
}
