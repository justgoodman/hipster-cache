package value_type

import(
	"fmt"
	"unsafe"
)

type OperationSetString struct {
	Err string
}

func (this *OperationSetString) SetValue(sourveValue, value interface{}) (valueSizeBytes int) {
	 switch v := value.(type) {
                case string:
                        byteSize += unsafe.Sizeof(sourvevalue) + len(v)
		default:
			this.Err := fmt.Sprintf(`Incorrect value type: type is not string`)
			return
        }
	sourceValue := value
}

type OperationGetString struct {
	value string
	Err string
}


func (this *OperationGetString) GetValue(value interface{}) {
	switch v := value.(type) {
                case string:
                       this.value = value
                default:
                        this.Err := fmt.Sprintf(`Incorrect value type: type is not string`)
        }
}
