package hash

type Chain struct {
        firstElement *ChainElement
        lastElement *ChainElement
        countElements int
}

type ChainElement struct {
        key string
        byteSize int
        value interface
}


func NewChain(firstElement *ChainElement) &Chain {
        return &Chain{firstElement: firstElement, lastElement: firstElement, countElements: 1}
}


func NewChainElement(key string, ) *Chain {
        chain := &Chain(key: key)
        chain.setValue(value)
        return chain
}

func (this *ChainElement) getNewByteSize(value interface) int {
        byteSize := unsage.Sizeof(this) + len(key)
        switch v := value.(type) {
                case string:
                        byteSize += unsafe.Sizeof(this) + len(v)
        }
}

func (this *ChainElement) setValue(value interface) int {
        newByteSize := getNewByteSize(value)
        deltaByteSize := newByteSize - this.byteSize
        this.byteSize = getNewByteSize(value)
        this.value = value
        return deltaByteSize
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
