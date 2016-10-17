package hash_table

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
)

type HashTable struct {
	capacity int64
	countElements uint64
        bytesSize uint64
	maxChainLenght uint64
	chains []*Chain
	hashFunction *ComplexStringHash
	mutexChains sync.RWMutex
	chainsMutex []*sync.RWMutex
	mutexChainsMutex sync.RWMutex
	lruChain *LRUChain
	// We will use LRU after we acheve this value and we will not use reHaching	
	maximumBytesSize uint

	countElementsMetric prometheus.Gauge
	bytesSizeMetric prometheus.Gauge
        maxChainLenghtMetric prometheus.Gauge
	responseTimeMetric *prometheus.SummaryVec
	hitCountMetric prometheus.Counter
	missCountMetric prometheus.Counter
}


type IBaseOperation interface {
	GetError() error
	GetOperationName() string
}

// For implement get operation for all types data you can implement this interface 
type IGetterValue interface {
     IBaseOperation
     GetValue(value interface{})
}

// For implement set operation for all types data you can implemet this interface
type ISetterValue interface {
     IBaseOperation
     SetValue(sourceValue interface{}, newValue interface{}) (valueSizeBytes int)
}


const (
	namespace = "hispter_cache"
	maximumLoadFactor = 0.9
)

func NewHashTable(capacity int64) *HashTable {
	return &HashTable{
			capacity: capacity,
			chains: make([]*Chain,capacity,capacity),
			chainsMutex: make([]*sync.RWMutex,capacity, capacity),
		}
}

func (h *HashTable) initMerics() {
	h.countElementsMetric = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "elements_total",
			Help:      "Count of elements stored in hipster cache",
			Namespace: namespace,
		})


	h.bytesSizeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
                        Name:      "bytes_total",
                        Help:      "Size of cache in bytes",
                        Namespace: namespace,
                })

	h.maxChainLenghtMetric = prometheus.NewGauge(prometheus.GaugeOpts{
                        Name:      "max_chain_lenght_total",
                        Help:      "Maximum size of chein",
                        Namespace: namespace,
                })

	h.responseTimeMetric = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name:      "response_time_microseconds",
			Help:      "Cache reponse time",
			Namespace: namespace,
		}, []string{"operation","is_error"})

	h.hitCountMetric = prometheus.NewCounter(prometheus.CounterOpts{
                        Name:      "hit_total",
                        Help:      "Hit count",
                        Namespace: namespace,
                })

	h.missCountMetric = prometheus.NewCounter(prometheus.CounterOpts{
                        Name:      "miss_total",
                        Help:      "Miss count",
                        Namespace: namespace,
                })

}

//  get duration in microsecond for sub: time1-time2
func getDurationMicroseconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Microsecond)
}


func (h *HashTable) SetElement(key string, expDate time.Time, value interface{}, setterValue ISetterValue) *ChainElement {
	timeStart := time.Now()

	chainElement,hasAdded, chainLenght, deltaBytes := h.setElement(key, expDate, value, setterValue)
	if hasAdded {
		countElements := atomic.AddUint64(&h.countElements,1)
		h.countElementsMetric.Set(float64(countElements))
		// Maybe we need reHaching
		h.reHaching()
		h.lruChain.Add(chainElement)
		h.lruEviction()
	} else {
		h.lruChain.MoveToFront(chainElement)
	}

	bytesSize := atomic.AddUint64(&h.bytesSize,uint64(deltaBytes))
	h.bytesSizeMetric.Set(float64(bytesSize))
	maxChainLenght := atomic.LoadUint64(&h.maxChainLenght)
	if maxChainLenght < uint64(chainLenght) {
		atomic.StoreUint64(&h.maxChainLenght, uint64(chainLenght))
		h.maxChainLenghtMetric.Set(float64(chainLenght))
	}

	responseMetric := h.responseTimeMetric.WithLabelValues(setterValue.GetOperationName(), boolToString(setterValue.GetError() != nil))
	responseMetric.Observe(getDurationMicroseconds(time.Since(timeStart)))

	return chainElement
}


func boolToString(value bool) string {
	if value {
		return "1"
	}
	return "0"
}
// Please don't modificate return value, it is not safe, we need this pointer to elemen for implementation LRU
func (h *HashTable) GetElement(key string,value interface{}, getterValue IGetterValue) *ChainElement {
	timeStart := time.Now()

	isHit,chainElement := h.getElement(key,value, getterValue)

	if chainElement.expDate.Unix() < time.Now().Unix() {
		isHit = false
		h.RemoveElement(chainElement)
	}

	if isHit {
		h.hitCountMetric.Inc()
	} else {
		h.missCountMetric.Inc()
	}

	responseMetric := h.responseTimeMetric.WithLabelValues("get")
	responseMetric.Observe(getDurationMicroseconds(time.Since(timeStart)))

	return chainElement
}

func (h *HashTable) RemoveElement(element *ChainElement) (sizeBytes uint64) {
	element.chainMutex.Lock()
	// We try to delete first element
	if element.prev == nil {
		key := element.key
		element.chainMutex.Unlock()
		// For delete first element we need to know chain
		h.deleteElement(key)
		element.chainMutex.Lock()
	} else {
		element.prev.next = element.next
		if element.next != nil {
			element.next.prev = element.prev
		}
	}

	prevLRU := element.prevLRU
	nextLRU := element.nextLRU
	sizeBytes =  uint64(unsafe.Sizeof(element)) + uint64(element.valueByteSize)
	element.chainMutex.Unlock()
	h.lruChain.delete(element,prevLRU,nextLRU)
	return
}


func (h *HashTable) setElement(key string, expDate time.Time, value interface{},setOperationObject ISetterValue) (chainElement *ChainElement,hasAdded bool, chainLenght int, deltaBytes uint64) {
       var (
		chain *Chain
		chainMutex *sync.RWMutex
		ok	bool
	)
	hashKey := h.hashFunction.CalculateHash(key)

	h.mutexChains.RLock()
	chain = h.chains[hashKey]
	h.mutexChains.RUnlock()

	if chain == nil {
		chainElement = NewChainElement(key)
		chainElement.setValue(setOperationObject,value, expDate)
		chain := NewChain(chainElement)

		deltaBytes = uint64(chainElement.valueByteSize) + uint64(unsafe.Sizeof(chainElement)) + uint64(unsafe.Sizeof(chain))
		hasAdded = true
		chainLenght = 1

		h.mutexChains.Lock()
		h.chains[hashKey] = chain
		h.mutexChains.Unlock()

		return
	}

	h.mutexChainsMutex.RLock()
	chainMutex = h.chainsMutex[hashKey]
	h.mutexChainsMutex.RUnlock()


	// @Check
	if chainMutex == nil {
		chainMutex := &sync.RWMutex{}
		h.mutexChainsMutex.Lock()
		h.chainsMutex[hashKey] = chainMutex
		h.mutexChainsMutex.Unlock()
	}
	chainMutex.Lock()

	chainElement = chain.findElement(key)

	if chainElement != nil {
		deltaBytes = uint64(chainElement.setValue(setOperationObject,value, expDate))
		hasAdded = false
	} else {
		chainElement = NewChainElement(key)
		chainElement.setValue(setOperationObject,value, expDate)
		chain.addElement(chainElement)

		deltaBytes = uint64(chainElement.valueByteSize) + uint64(unsafe.Sizeof(chain)) + uint64(unsafe.Sizeof(chainElement))
		hasAdded = false
	}

	chainLenght = chain.countElements

	chainMutex.Unlock()

	return
}


func (h *HashTable) getElement(key string,  value interface{}, getOperationObject IGetterValue) (isHit bool, chainElement *ChainElement) {
       var (
		chain *Chain
		chainMutex *sync.RWMutex
		ok	bool
	)
	hashKey := h.hashFunction.CalculateHash(key)

	h.mutexChains.RLock()
	chain = h.chains[hashKey]
	h.mutexChains.RUnlock()

	if chain == nil {
		return
	}

	h.mutexChainsMutex.RLock()
	chainMutex = h.chainsMutex[hashKey]
	h.mutexChainsMutex.RUnlock()

	if chainMutex == nil {
		return
	}

	chainMutex.RLock()

	chainElement = chain.findElement(key)
	isHit = true

	if chainElement != nil {
		chainElement.getValue(getOperationObject)
	}

	chainMutex.RUnlock()
	return
}

func (h *HashTable) deleteElement(key string) (deletedBytes uint64) {
       var (
		chain *Chain
		chainMutex *sync.RWMutex
		ok	bool
	)
	hashKey := h.hashFunction.CalculateHash(key)

	h.mutexChains.RLock()
	chain = h.chains[hashKey]
	h.mutexChains.RUnlock()

	if chain == nil {
		return
	}

	h.mutexChainsMutex.RLock()
	chainMutex = h.chainsMutex[hashKey]
	h.mutexChainsMutex.RUnlock()

	if chainMutex == nil {
		return
	}

	chainMutex.RLock()

	chainElement := chain.findElement(key)

	if chainElement != nil {
		deletedBytes = uint64(chainElement.valueByteSize) + uint64(unsafe.Sizeof(chainElement))
		chain.deleteElement(chainElement)
	}

	chainMutex.RUnlock()
	return
}

func (h *HashTable) lruEviction() {
	sizeBytes := atomic.LoadUint64(&h.sizeBytes)
       if h.maximumSizeBytes <= sizeBytes {
		return
	}

	// We will free 10% of maximumSizeBytes
	needFreeSize = sizeBytes - (h.maxumumSizeByte - h.maximumSizeBytes/10)

	var freeBytes uint
	var element *ChainElement
	var countElements int
	for element = h.lruChain.lastElement; element != nil; element = element.prevLRU {
		freeBytes += h.deleteElement(element)
		countElements++
		if freeBytes >= needFreeSize {
			break
		}
	}
	atomic.AddUint64(&h.sizeBytes, -1*freeBytes)
	atomic.AddUint64(&h.countElements, -1*countElements)
}


func (h *HashTable) reHaching() {
	countChanks := atomic.LoadUint64(&h.capacity)
	countElements := atomic.LoadUint64(&h.countElements)
	sizeBytes := atomic.LoadUint64(&h.sizeBytes)

	if countElement/countChanks <= 0.9 {
		return
	}

	newSizeBytes := sizeBytes + (nosafe.Sizeof(*Chank) + nosafe.Sizeof(*sync.RWMutex)) * countChanks

	if h.maximumSizeBytes <= newSizeBytes {
		return
	}

	newCapacity := countChanks*2
	newChanks := [newCapacity]*Chank
	newMutexes := [newCapacity]*sync.RWMutex
	newHashFunction := NewHashFunction()

	h.mutexChain.Lock()

	var chainElement *ChainElement
	for _, chain := range(h.chainsMutex) {
		if chain == nil {
			continue
		}
		for chainElement := chain.FirstElement; chainElement != nill; chainElement = nextChainElement {
			nextChainElement = chainElement.next
			// We will use it element in the new chain
			chainElement.next = nil
			h.reHachingAddElement(newChains,newHashFunction,chainElement)
		}
	}

	h.capacity = newCapacity
	h.chainMutexes = newMutexes
	h.hasFunction = newHashFunction
	h.sizeBytes = newSizeBytes

	h.mutexChain.Unlock()

	// AddNewMetic
	h.rehashingCountMetric.Inc()
	h.countBytesMetric.Set(newSizeBytes)
}

func (h *HashTable) reHachingAddElement(chains []*Chain, hashFunction *hashFunction, chainElement *ChainElement) {
       var (
		chain *Chain
		chainMutex *sync.RWMutex
		ok	bool
	)
	hashKey := hashFunction.CalculateHash(chainElement.key)

	chain = chains[hashKey]

	if chain == nil {
		chain = NewChain(chainElement)
		h.chains[hashKey] = chain
		return
	}

	chain.AddElement(chainElement)
}
