package hash

import (
	"sync"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
)

type HashTable struct {
	capacity int64
	countElements uint64
        bytesSize uint64
	maxChainLenght int
	chains []*Chain
	hashFunction *ComplexStringHash
	mutexChains sync.RWMutex
	chainsMutex []*sync.RWMutex
	mutexChainsMutex sync.RWMutex
	lruChain *Chain

	countElementsMetric *prometheus.Gauge
	bytesSizeMetric *prometheus.Gauge
        maxChainLenghtMetic *prometheus.Gauge
	responseTimeMetric *prometheus.SummaryVec
}


type IBaseOperation interface {
	func GetError() error
	func GetOperationName() string
}

// For implement get operation for all types data you can implement this interface 
type IGetterValue interface {
     IBaseOperation
     func GetValue(value interface{})
}

// For implement set operation for all types data you can implemet this interface
type ISetterValue interface {
     IBaseOperation
     func SetValue(sourceValue interface{}, newValue interface{}) (valueSizeBytes int)
}

const (
	namespace = "hispter_cache"
	maximumLoadFactor = 0.9
)

func NewHashTable(capacity int) *HashTable {
	return &HashTable{
			capacity: capacity,
			chains: [capacity]*Chain,
			chainsMutex: [capacity]*sync.RWMutex
		}
}

func (this *HashTable) initMerics() {
	this.countElementsMetric := prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "elements_total",
			Help:      "Count of elements stored in hipster cache",
			Namespace: namespace,
		})


	this.bytesSizeMetric := prometheus.NewGauge(prometheus.GaugeOpts{
                        Name:      "bytes_total",
                        Help:      "Size of cache in bytes",
                        Namespace: namespace,
                })

	this.maxCheinLenghtMetric := prometheus.NewGauge(prometheus.GaugeOpts{
                        Name:      "max_chain_lenght_total",
                        Help:      "Maximum size of chein",
                        Namespace: namespace,
                })

	this.responseTimeMetic = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name:      "response_time_microseconds",
			Help:      "Cache reponse time",
			Namespace: NAMESPACE,
		}, []string{"operation","is_error"})

	this.hitCountMetric := prometheus.NewCount(prometheus.CountOpts{
                        Name:      "hit_total",
                        Help:      "Hit count",
                        Namespace: namespace,
                })

	this.missCountMetric := prometheus.NewCount(prometheus.CountOpts{
                        Name:      "miss_total",
                        Help:      "Miss count",
                        Namespace: namespace,
                })

}

//  get duration in microsecond for sub: time1-time2
func getDurationMicroseconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Microsecond)
}


// Please don't modificate return value, it is not safe, we need this pointer to elemen for implementation LRU
func (this *HashTable) SetElement(key string, expDate time.Time, value interface{}, setterValue ISetterValue) *ChainElement {
	timeStart := time.Now()

	chainElement,hasAdded, chainLenght, deltaBytes := this.setElement(key, expDate, value, setterValue)
	if hadAdded {
		countElements := atomic.AddUint64(&this.countElements,1)
		this.countElementsMetric.Set(float64(countElemenets))
		// Maybe we need reHaching
		this.reHaching()
		this.lruChain.Add(chainElement)
	} else {
		this.lruChain.MoveToFront(chainElement)
	}

	bytesSize := atomic.AddUint64(&this.bytesSize,deltaBytes)
	this.bytesSizeMetric.Set(float64(butesSize))
	maxChainLenght = atomic.LoadUint64(&this.maxChainLenght)
	if maxChainLenght < chainLenght {
		atomic.StoreUint64(*this.maxChainLenght, chainLenght)
		this.maxChainLenghtMetric.Set(float64(chainLenght))
	}

	responseMetric := this.responseTimeMetric.WithLabelValues(setterValue.GetOpationName(), setterValue.GetError() != nil)
	responseMetric.Observe(getDurationMicroseconds(time.Since(timeStart))

	return chainElement
}

// Please don't modificate return value, it is not safe, we need this pointer to elemen for implementation LRU
func (this *HashTable) GetElement(key string,value interface{}, getterValue IGetterValue) *ChainElement {
	timeStart := time.Now()

	isHit,chainElement := this.getElement(key,getterValue)

	if chainElement.expDate < time.Now() {
		isHit = false
		this.RemoveElement(chainElement)
	}
	if isHit {
		this.hitMetric.Inc()
	}
	else {
		this.missMetric.Inc()
	}

	responseMetric := this.responseTimeMetric.WithLabelValues("get")
	responseMetric.Observe(getDurationMicroseconds(time.Since(timeStart))
}

func (this *HashTable) RemoveElement(element *ChainElement) {
	element.chainMutex.Lock()
	// We try to delete first element
	if element.prev = nil {
		key := element.key
		element.chainMutex.Unlock()
		// For delete first element we need to know chain
		this.removeElement(key)
		element.chainMutex.Lock()
	} else {
		element.prev.next := element.next
		if element.next != nil {
			element.next.prev = element.prev
	}

	prevLRU := element.prevLRU
	nextLRU := element.nextLRU
	sizeBytes := element.sizeByte
	element.chainMutex.Unlock()
	this.lruChain.delete(element,prevLRU,nextLRU)
}


func (this *HashTable) SetElement(key string, expDate time.Time, setOpertionObject ISetOperation, value interface) (chainElement *ChainElement,hasAdded bool, chainLenght int, deltaBytes int) {
       var (
		chain *Chain
		chainMutex *sync.RWMutex
		ok	bool
	)
	hashKey := hashFunction.CalculateHash(key)

	this.mutexChain.RLock()
	chain, ok := this.chains[hashKey]
	this.mutexChain.RULock()

	if !ok {
		chainElement = NewChainElement(key)
		chainElement.SetValue(setOperationObject,value, expDate)
		chain := NewChain(chainElement)

		deltaBytes = chainElement.byteSize + unsafe.Sizeof(chain)
		hasAdded = true
		chainLeight = 1

		this.mutexChain.Lock()
		this.chains[hashKey] := chain
		this.mutexChain.Unlock()

		// For safe work with object
		chainElement = chainElement.copy()
		return
	}

	this.mutexChainMutexes.RLock()
	chainMutex,ok := this.chainsMutex[hashKey]
	this.mutexChainMutexed.RULock()

	if !ok {
		chainMutex := &sync.Mutex{}
		this.mutexChainMutexes.Lock()
		this.chainsMutex[hasKey] = chainMutex
		this.mutexChainMutexes.Unlock()
	}

	chainMutex.Lock()

	chainElement = chain.findElement(key)

	if chainElement != nil {
		deltaBytes = element.setValue(value)
		hasAdded = false
	} else {
		chainElement = NewChainElement(key)
		chainElement.setValue(setOperationObject,value, expDate)
		chain.AddElement(chainElement)

		deltaBytes = chainElement.sizeButes
		hasAdded = false
	}

	chainLenght = chain.lenght

	// For safe work with object
	chainElement := chainElement.copy()
	chainMutex.Unlock()
}


func (this *HashTable) getElement(key string, getOpertionObject IGetOperation, value interface) (isHit bool, chainElement *ChainElement) {
       var (
		chain *Chain
		chainMutex *sync.RWMutex
		ok	bool
	)
	hashKey := hashFunction.CalculateHash(key)

	this.mutexChain.RLock()
	chain, ok := this.chains[hashKey]
	this.mutexChain.RULock()

	if !ok {
		return
	}

	this.mutexChainMutexes.RLock()
	chainMutex,ok := this.chainsMutex[hashKey]
	this.mutexChainMutexed.RULock()

	if !ok {
		return
	}

	chainMutex.RLock()

	chainElement = chain.findElement(key)
	isHit = true

	if chainElement != nil {
		deltaBytes = chainElement.getValue(getOperationObject)
	}

	chainLenght = chain.lenght

	// For safe return Element
	chainElement = chanElement.copy()

	chainMutex.RUnlock()
}

func (this *HashTable) deleteElement(key string) (deletedBytes int) {
       var (
		chain *Chain
		chainMutex *sync.RWMutex
		ok	bool
	)
	hashKey := hashFunction.CalculateHash(key)

	this.mutexChain.RLock()
	chain, ok := this.chains[hashKey]
	this.mutexChain.RULock()

	if !ok {
		return
	}

	this.mutexChainMutexes.RLock()
	chainMutex,ok := this.chainsMutex[hashKey]
	this.mutexChainMutexed.RULock()

	if !ok {
		return
	}

	chainMutex.RLock()

	chainElement := chain.findElement(key)

	if chainElement != nil {
		deltaBytes = chainElement.sizeBytes
		chain.RemoveElement(chainElement)
	}

	chainLenght = chain.lenght

	chainMutex.RUnlock()
}



func (this *HashTable) reHaching() {
	countChanks := atomic.LoadUint64(this.&capacity)
	countElements := atomic.LoadUint64(this.&countElements)

	if countElement/countChanks <= 0.9 {
		return
	}

	newCapacity := countChanks*2
	newChanks := [newCapacity]*Chank
	newMutexes := [newCapacity]
	newHashFunction := NewHashFunction()
	newSizeBytes := this.sizeBytes + (nosafe.Sizeof(*Chank) + nosafe.Sizeof(*sync.RWMutex)) * countChanks

	this.mutexChain.Lock()

	var chainElement *ChainElement
	for _, chain := range(this.chainsMutex) {
		if chain == nil {
			continue
		}
		for chainElement := chain.FirstElement, chainElement != nill, chainElement = nextChainElement) {
			nextChainElement = chainElement.next
			// We will use it element in the new chain
			chainElement.next = nil
			this.reHachingAddElement(newChains,newHashFunction,chainElement)
		}
	}

	this.capacity = newCapacity
	this.chainMutexes = newMutexes
	this.hasFunction = newHashFunction
	this.sizeBytes = newSizeBytes

	this.mutexChain.Unlock()

	// AddNewMetic
	this.rehashingCountMetric.Inc()
	this.countBytesMetric.Set(newSizeBytes)
}

func (this *HashTable) reHachingAddElement(chains []*Chain, hashFunction *hashFunction, chainElement *ChainElement) {
       var (
		chain *Chain
		chainMutex *sync.RWMutex
		ok	bool
	)
	hashKey := hashFunction.CalculateHash(chainElement.key)

	chain, ok := chains[hashKey]

	if !ok {
		chain := NewChain(chainElement)
		this.chains[hashKey] := chain
		return
	}

	chain.AddElement(chainElement)
}
