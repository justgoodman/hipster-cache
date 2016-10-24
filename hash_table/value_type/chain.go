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
	next *ChainElement
}

func NewChain(firstElement *ChainElement) *Chain {
        chain := &Chain{firstElement: firstElement, lastElement: firstElement, countElements: 1}
	return chain
}

func NewChainElement(value string) *ChainElement {
        chainElement := &ChainElement{value: value}
        chainElement.bytesSize = int(unsafe.Sizeof(chainElement)) + len(value)
        return chainElement
}

func (this *Chain) findElement(index int) *ChainElement {
        i := 0
        for element := this.firstElement; element != nil; element = element.next {
                if i == index {
                        return element
                }
                i += 1
        }
	return nil
}

func (this *Chain) GetRangeValues(indexStart,indexEnd int) []string {
        values := []string{}
        i := 0
        for element := this.firstElement; element != nil; element = element.next {
                if i < indexStart || i > indexEnd {
                      return values
                }

                if i >= indexStart && i <= indexEnd {
                        values = append(values,element.value)
                }
        }
        return values
}

func (this *Chain) addElement(element *ChainElement) {
        this.lastElement.next = element
        this.lastElement = element
        this.countElements += 1
}
