package value_type

import (
	"unsafe"
)

type Chain struct {
	firstElement  *ChainElement
	lastElement   *ChainElement
	countElements int
	bytesSize     int
}

type ChainElement struct {
	bytesSize int
	value     string
	key       string
	next      *ChainElement
}

func (e *ChainElement) updateValue(value string) {
	lastString := e.value
	e.value = value
	e.bytesSize += len(lastString) - len(value)
}

func NewChain(firstElement *ChainElement) *Chain {
	chain := &Chain{firstElement: firstElement, lastElement: firstElement, countElements: 1}
	chain.bytesSize = int(unsafe.Sizeof(chain)) + firstElement.bytesSize
	return chain
}

func NewChainElement(key, value string) *ChainElement {
	chainElement := &ChainElement{key: key, value: value}
	chainElement.bytesSize = int(unsafe.Sizeof(chainElement)) + len(value)
	return chainElement
}

func (c *Chain) findElement(key string) *ChainElement {
	for element := c.firstElement; element != nil; element = element.next {
		if element.key == key {
			return element
		}
	}
	return nil
}

func (c *Chain) getAllValues() []*DictElement {
	values := []*DictElement{}
	for element := c.firstElement; element != nil; element = element.next {
		values = append(values, NewDictElement(element.key, element.value))
	}
	return values
}

func (c *Chain) addElement(element *ChainElement) {
	c.lastElement.next = element
	c.lastElement = element
	c.countElements += 1
	c.bytesSize += element.bytesSize
}
