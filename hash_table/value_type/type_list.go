package value_type

import (
	"fmt"
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

	switch slice := (*sourceValue).(type) {
		case []string:
			slice = append(slice, stringValue)
			*sourceValue = interface{}(slice)
		// It is like nil, default value
		case string:
			if slice != "" {
				l.err = fmt.Errorf("Error: list type is string")
			}
			newSlice := []string{}
			newSlice = append(newSlice,stringValue)
			*sourceValue = interface{}(newSlice)
		default:
			l.err = fmt.Errorf("Error: list type is not the []string")
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
	switch slice := value.(type) {
		case []string:
			l.Lenght = len(slice)
		default:
			l.err = fmt.Errorf("Error: list type in not the []string")
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

func GetRangeValues(slice []string, indexStart,indexEnd int) []string {
        values := []string{}

	for i,value := range slice {
                if i > indexEnd {
                     break
                }

                if i >= indexStart && i <= indexEnd {
                        values = append(values,value)
                }
        }
        return values
}

func (l *RangeListOperation) GetValue(value interface{}) {
	switch slice := value.(type) {
		case []string:
			l.Values = GetRangeValues(slice,l.indexStart, l.indexEnd)
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

	switch slice := (*sourceValue).(type) {
		case []string:
			if len(slice) <= l.index || l.index < 0 {
				l.err = fmt.Errorf(`Error: Index "%d" not found in List`, l.index)
				return
			}
			slice[l.index] = stringValue
		default:
			l.err = fmt.Errorf("Error: list type in not the []string")
	}
	return
}

func (l *SetListOperation) GetResult() (string,error) {
	if l.err == nil {
		return "OK",l.err }
	return "",l.err
}
