package hash

type Chain struct {
        lastElement *ChainElement
        countElements int
}

type LRUChain struct {
	lastElement *ChainElement
	firstElement *ChainElement
	countElements int
	mutex sync.Mutex
}

func NewLRUChain() *LRUChain {
	retrun &LRUChain{}
}

func (this *LRUChain) Add(element *ChainElement) {
	this.mutex.Lock()
	if this.firstElement != nil {
		this.firstElement.mutexChain.Lock()
		this.firstElement.prevLRU = element
		this.firstElement.mutexChain.Unlock()
	}
	this.element.mutexChain.Lock()
	element.nextLRU = this.firstElement
	this.element.mutexChain.Unlock()
	this.firstElement = element

	if this.lastElement == nil {
		this.lastElement = this.firstElement
	}
	this.countElements++
	this.mutex.Unlock()
}

func (this *LRUChain) delete(element *ChainElement, prevLRU *ChainElement, nextLRU *ChainElement) {
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

	this.mutex.Lock()
	if this.firstElement == element {
		this.firstElement = nextLRU
	}
	if this.lastElement == element {
		this.lastElement = prevLRU
	}
	this.mutex.Unlock()
}

func (this *LRUChain) MoveToFront(element *ChainElement) {
	this.mutex.Lock()
	if this.lastElement == element {
		this.mutex.Unlock()
		return
	}
	this.mutex.Unlock()

	element.chainMutex.Lock()
	prevLRU := element.prevLRU
        nextLRU := element.nextLRU
	element.chainMutex.Unlock()
	this.delete(element,prevLRU,nextLRU)
	this.Add(element)
}


type ChainElement struct {
	next *ChainElement
	prev *ChainElement
	nextLRU *ChainElement
	prevLRU *ChainElement
	// Need for working with LRU(because in LRU element can be from another chain)
	chainMutex *sync.Mutex
	// Experation Date
	expDate time.Time
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

func (this *Chain) deleteElemenet(element *ChainElement) {
	if this.firstElement == element {
		this.firstElement = element.next
		if element.next != nil {
			element.next.prev = nil
		}
		return
	}
	element.prev.next = element.next
	if elemenet.next != nil {
		element.next.prev = element.prev
	}
	this.lenght--
}


func (this *ChainElement) setValue(setterValue ISetterValue,value interface, expDate time.time) (deltaBytes int) {
	this.expDate = expDate
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
	element.PrevElement :=  this.lastElement
        this.lastElement.next = element
        this.lastElement = element
        this.lenght += 1
}
