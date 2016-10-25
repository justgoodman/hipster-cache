package value_type

import (
	"fmt"
	"strconv"
	"unsafe"
)

const (
	PushListCmdName = "LPUSH"
	RangeListCmdName = "LRANGE"
	SetListCmdName = "LSET"
	LenghtListCmdName = "LLEN"
)

type sliceString struct {
	slice []string
	bytesSize int
}

func NewSliceString() *sliceString {
	slice := &sliceString{}
	slice.bytesSize = int(unsafe.Sizeof(slice))
	return slice
}

func (s *sliceString) addElement(value string) {
	s.slice = append(s.slice, value)
	s.bytesSize += len(value) + int(unsafe.Sizeof(value))
}

func (s *sliceString) setElement(index int, value string) error {
	if len(s.slice) <= index || index < 0 {
		return fmt.Errorf(`Error: Index "%d" not found in List`, index)
	}
	lastString := s.slice[index]
	s.slice[index] = value
	s.bytesSize += len(value) - len(lastString)
	return nil
}

func (s *sliceString) getRangeValues(indexStart,indexEnd int) []string {
        values := []string{}

	for i,value := range s.slice {
                if i > indexEnd {
                     break
                }

                if i >= indexStart && i <= indexEnd {
                        values = append(values,value)
                }
        }
        return values
}

func (s *sliceString) getLenght() int {
	return len(s.slice)
}

type PushListOperation struct {
	baseOperation
	elementIndex int
}

func NewPushListOperation() *PushListOperation {
	return &PushListOperation{baseOperation: baseOperation{commandName: PushListCmdName}}
}

func (l *PushListOperation) SetValue(sourceValue *interface{}, value interface{}) (valueBytesSize int) { var stringValue string
	switch setValue :=  value.(type) {
		case string:
			stringValue = setValue
		default:
			l.err = fmt.Errorf("Error: list value in not the string")
			return
	}

	switch slice := (*sourceValue).(type) {
		case *sliceString:
			slice.addElement(stringValue)
			valueBytesSize = slice.bytesSize
		// It is like nil, default value
		case string:
			if slice != "" {
				l.err = fmt.Errorf("Error: list type is string")
			}
			newSlice := NewSliceString()
			newSlice.addElement(stringValue)
			valueBytesSize = newSlice.bytesSize
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
		case *sliceString:
			l.Lenght = slice.getLenght()
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


func (l *RangeListOperation) GetValue(value interface{}) {
	switch slice := value.(type) {
		case *sliceString:
			l.Values = slice.getRangeValues(l.indexStart, l.indexEnd)
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

func (l *SetListOperation) SetValue(sourceValue *interface{},value interface{}) (valueBytesSize int) {
	var stringValue string
	switch setValue :=  value.(type) {
		case string:
			stringValue = setValue
		default:
			l.err = fmt.Errorf("Error: list value in not the string")
			return
	}

	switch slice := (*sourceValue).(type) {
		case *sliceString:
			l.err = slice.setElement(l.index,stringValue)
			valueBytesSize = slice.bytesSize
		default:
			l.err = fmt.Errorf("Error: list type in not the []string")
	}
	return
}

func (l *SetListOperation) GetResult() (string,error) {
	if l.err == nil {
		return "OK",l.err
	}
	return "",l.err
}
