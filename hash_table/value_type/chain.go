package value_type

import (
	"unsafe"
)

type Chain struct {
        firstElement *ChainElement
        lastElement *ChainElement
        countElements int
        bytesSize int
}

type ChainElement struct {
        bytesSize int
        value string
	key string
	next *ChainElement
}

func NewChain(firstElement *ChainElement) *Chain {
        chain := &Chain{firstElement: firstElement, lastElement: firstElement, countElements: 1}
	return chain
}

func NewChainElement(key,value string) *ChainElement {
        chainElement := &ChainElement{key:key, value: value}
        chainElement.bytesSize = int(unsafe.Sizeof(chainElement)) + len(value)
        return chainElement
}

func (this *Chain) findElement(key string) *ChainElement {
        for element := this.firstElement; element != nil; element = element.next {
                if element.key == key {
                        return element
                }
        }
	return nil
}

func (this *Chain) getAllValues() []*DictElement {
        values := []*DictElement{}
        for element := this.firstElement; element != nil; element = element.next {
                  values = append(values,NewDictElement(element.key, element.value))
        }
        return values
}

func (this *Chain) addElement(element *ChainElement) {
        this.lastElement.next = element
        this.lastElement = element
        this.countElements += 1
}
