package hash_table

import (
	"sync"
	"time"
	"fmt"
)

type Chain struct {
	firstElement  *ChainElement
	lastElement   *ChainElement
	countElements int
	mutex *sync.RWMutex
}

type LRUChain struct {
	lastElement   *ChainElement
	firstElement  *ChainElement
	countElements int
	mutex         sync.Mutex
}

func NewLRUChain() *LRUChain {
	return &LRUChain{}
}

func (c *LRUChain) Add(element *ChainElement) {
	c.mutex.Lock()
	if c.firstElement != nil {
		c.firstElement.chainMutex.Lock()
		c.firstElement.prevLRU = element
		c.firstElement.chainMutex.Unlock()
	}
	element.chainMutex.Lock()
	element.nextLRU = c.firstElement
	element.chainMutex.Unlock()

	c.firstElement = element

	if c.lastElement == nil {
		c.lastElement = c.firstElement
	}
	c.countElements++
	c.mutex.Unlock()
}

func (c *LRUChain) delete(element *ChainElement, prevLRU *ChainElement, nextLRU *ChainElement) {
	if prevLRU != nil {
		prevLRU.chainMutex.Lock()
		prevLRU.nextLRU = nextLRU
		prevLRU.chainMutex.Unlock()
	}
	if nextLRU != nil {
		nextLRU.chainMutex.Lock()
		nextLRU.prevLRU = prevLRU
		nextLRU.chainMutex.Unlock()
	}

	c.mutex.Lock()
	if c.firstElement == element {
		c.firstElement = nextLRU
	}
	if c.lastElement == element {
		c.lastElement = prevLRU
	}
	c.mutex.Unlock()
}

func (c *LRUChain) MoveToFront(element *ChainElement) {
	c.mutex.Lock()
	if c.lastElement == element {
		c.mutex.Unlock()
		return
	}
	c.mutex.Unlock()

	element.chainMutex.Lock()
	prevLRU := element.prevLRU
	nextLRU := element.nextLRU
	element.chainMutex.Unlock()
	c.delete(element, prevLRU, nextLRU)
	c.Add(element)
}

type ChainElement struct {
	next    *ChainElement
	prev    *ChainElement
	nextLRU *ChainElement
	prevLRU *ChainElement
	// Need for working with LRU(because in LRU element can be from another chain)
	chainMutex *sync.RWMutex
	// Experation Date
	expDate       time.Time
	key           string
	valueByteSize int
	value         interface{}
}

func NewChain(mutex *sync.RWMutex) *Chain {
	return &Chain{mutex: mutex}
}

func NewChainElement(key string)  *ChainElement {
	chainElement := &ChainElement{key: key,value: interface{}("")}
	return chainElement
}

func (c *Chain) deleteElement(element *ChainElement) {
	fmt.Printf("\n Delete element \n")
	if c.firstElement == element {
		c.firstElement = element.next
		if element.next != nil {
			element.next.prev = nil
		}
		if c.lastElement == element {
			c.lastElement = nil
		}
		return
	}
	element.prev.next = element.next
	if element.next != nil {
		element.next.prev = element.prev
	}
	if c.lastElement == element {
		c.lastElement = element.prev
	}
	c.countElements--
}

func (e *ChainElement) setValue(setterValue ISetterValue, value interface{}, expDate time.Time) (deltaBytes int) {
	e.expDate = expDate
	newValueByteSize := setterValue.SetValue(&e.value, value)

	deltaByteSize := newValueByteSize - e.valueByteSize
	e.valueByteSize = newValueByteSize
	return deltaByteSize
}

func (e *ChainElement) getValue(getterValue IGetterValue) {
	getterValue.GetValue(e.value)
}

func (c *Chain) findElement(key string) *ChainElement {
	fmt.Printf("\n FIND ELEMENT \n")
	for element := c.firstElement; element != nil; element = element.next {
		fmt.Printf("\n Go thow element:%#v", element)
		if element.key == key {
			fmt.Printf("\n CHAIN:%#v \n", c)
			return element
		}
	}
	return nil
}

func (c *Chain) AddElement(element *ChainElement) {
	element.chainMutex = c.mutex
	if c.lastElement == nil {
		c.firstElement = element
		c.lastElement = element
		element.prev = nil
		element.next = nil
		c.countElements++
		return
	}
	element.prev = c.lastElement
	c.lastElement.next = element
	c.lastElement = element
	c.countElements++
}
