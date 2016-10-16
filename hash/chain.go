package hash

type Chain struct {
        firstElement *ChainElement
        lastElement *ChainElement
        countElements int
}

type ChainElement struct {
        key string
        valueByteSize int
        value interface
}


func NewChain(firstElement *ChainElement) &Chain {
        return &Chain{firstElement: firstElement, lastElement: firstElement, countElements: 1}
}


func NewChainElement(key string) *Chain {
        chain := &ChainElement(key: key)
	chain.byteSize := unsage.Sizeof(this) + len(key)
        return chain
}

func (this *ChainElement) getNewByteSize(value interface) int {
        byteSize := unsage.Sizeof(this) + len(key)
        switch v := value.(type) {
                case string:
                        byteSize += unsafe.Sizeof(this) + len(v)
        }
	return byteSize
}

func (this *ChainElement) setValue(setterValue ISetterValue,value interface) (deltaBytes int) {
	newValueByteSize := setterValue.setValue(this.value, value)

        deltaByteSize = newValueByteSize - this.valueByteSize
        this.valueByteSize = valueByteSize
        return deltaByteSize
}

func (this *ChainElement) getValue(getterValue IGetterValue) {
	getterValue.getValue(this.value)
}

func (this *Chain) findElement(key string) *ChainElement {
        for element := this.firstElement; element != nil; element = element.next {
                if element.key == key {
                        return element
                }
        }
}

func (this *Chain) addElement(element *ChainElement) {
        this.lastElement.next = element
        this.lastElement = element
        this.lenght += 1
}
