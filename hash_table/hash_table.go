package hash_table

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type HashTable struct {
	capacity         int64
	countElements    int64
	bytesSize        int64
	maxChainLenght   int64
	chains           []*Chain
	hashFunction     *ComplexStringHash
	mutexChains      sync.RWMutex
	chainsMutex      []*sync.RWMutex
	mutexChainsMutex sync.RWMutex
	reHashingMutex   sync.RWMutex
	lruChain         *LRUChain
	// We will use LRU after we acheve this value and we will not use reHaching
	maxBytesSize int64

	// Coefficient for hash function
	coefPString uint64
	// Coefficient for hash function
	coefPInt uint64

	countElementsMetric  prometheus.Gauge
	bytesSizeMetric      prometheus.Gauge
	maxChainLenghtMetric prometheus.Gauge
	responseTimeMetric   *prometheus.SummaryVec
	reHashingTimeMetric  prometheus.Summary
	hitCountMetric       prometheus.Counter
	missCountMetric      prometheus.Counter
}

type IBaseOperation interface {
	GetError() error
	GetCommandName() string
}

// For implement get operation for all types data you can implement this interface
type IGetterValue interface {
	IBaseOperation
	GetValue(value interface{})
}

// For implement set operation for all types data you can implemet this interface
type ISetterValue interface {
	IBaseOperation
	SetValue(sourceValue *interface{}, newValue interface{}) (valueSizeBytes int)
}

const (
	namespace         = "hispter_cache"
	maximumLoadFactor = 0.9
)

func NewHashTable(initCapacity int64, maxKeyLenght int64, maxBytesSize int64) *HashTable {
	sizeOfOneChain := int64(unsafe.Sizeof(&Chain{})) + int64(unsafe.Sizeof(&ChainElement{})) + int64(unsafe.Sizeof(Chain{}))
	maxCapacity := int64(maxBytesSize / sizeOfOneChain)
	// CoefP must be more than maximCapacity*maxKeyLenght
	coefP := uint64(maxCapacity*maxKeyLenght + (rand.Int63n(math.MaxInt64-int64(maxCapacity*maxKeyLenght))))
	return &HashTable{
		capacity:     initCapacity,
		maxBytesSize: maxBytesSize,
		chains:       make([]*Chain, initCapacity, initCapacity),
		chainsMutex:  make([]*sync.RWMutex, initCapacity, initCapacity),
		hashFunction: NewComplexStringHash(uint64(initCapacity), coefP, coefP),
		coefPString: coefP,
		coefPInt: coefP,
		lruChain: NewLRUChain(),
	}
}

func (h *HashTable) InitMetrics() {
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
	}, []string{"operation", "is_error"})

	h.reHashingTimeMetric = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "rehashing_time_microseconds",
		Help:      "Rehashing duration",
		Namespace: namespace,
	})
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
	fmt.Printf("\n Set element:%#v", value)
	timeStart := time.Now()

	chainElement, hasAdded, chainLenght, deltaBytes := h.setElement(key, expDate, value, setterValue)
	if hasAdded {
		countElements := atomic.AddInt64(&h.countElements, 1)
		h.countElementsMetric.Set(float64(countElements))
		// Maybe we need reHaching
		h.reHashing()
		h.lruChain.Add(chainElement)
		h.lruEviction()
	} else {
		h.lruChain.MoveToFront(chainElement)
	}

	bytesSize := atomic.AddInt64(&h.bytesSize, int64(deltaBytes))
	h.bytesSizeMetric.Set(float64(bytesSize))
	maxChainLenght := atomic.LoadInt64(&h.maxChainLenght)
	if maxChainLenght < int64(chainLenght) {
		atomic.StoreInt64(&h.maxChainLenght, int64(chainLenght))
		h.maxChainLenghtMetric.Set(float64(chainLenght))
	}

	responseMetric := h.responseTimeMetric.WithLabelValues(setterValue.GetCommandName(), boolToString(setterValue.GetError() != nil))
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
func (h *HashTable) GetElement(key string, getterValue IGetterValue) {
	timeStart := time.Now()

	isHit, chainElement := h.getElement(key, getterValue)

	if !isHit {
		return
	}

	fmt.Printf("\n Find Element:%#v \n", chainElement)
	if chainElement.expDate.Unix() < time.Now().Unix() {
		isHit = false
		h.removeElement(chainElement)
	}

	if isHit {
		h.hitCountMetric.Inc()
	} else {
		h.missCountMetric.Inc()
	}

	responseMetric := h.responseTimeMetric.WithLabelValues("get", boolToString(getterValue.GetError() != nil))
	responseMetric.Observe(getDurationMicroseconds(time.Since(timeStart)))
	return
}

func (h *HashTable) removeElement(element *ChainElement) (sizeBytes int64) {
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
	sizeBytes = int64(unsafe.Sizeof(element)) + int64(element.valueByteSize)
	element.chainMutex.Unlock()
	h.lruChain.delete(element, prevLRU, nextLRU)
	return
}

func (h *HashTable) setElement(key string, expDate time.Time, value interface{}, setOperationObject ISetterValue) (chainElement *ChainElement, hasAdded bool, chainLenght int, deltaBytes int64) {
	var (
		chain      *Chain
		chainMutex *sync.RWMutex
	)
	h.reHashingMutex.RLock()
	defer h.reHashingMutex.RUnlock()


	hashKey := h.hashFunction.CalculateHash(key)

	fmt.Printf("\n SET hash key:%d", hashKey)
	h.mutexChains.RLock()
	chain = h.chains[hashKey]
	h.mutexChains.RUnlock()

	if chain == nil {
		chain := NewChain(&sync.RWMutex{})

		chainElement = NewChainElement(key)
		chainElement.setValue(setOperationObject, value, expDate)

		chain.AddElement(chainElement)

		deltaBytes = int64(chainElement.valueByteSize) + int64(unsafe.Sizeof(chainElement)) + int64(unsafe.Sizeof(chain))
		hasAdded = true
		chainLenght = 1

		h.mutexChains.Lock()
		h.chains[hashKey] = chain
		h.mutexChains.Unlock()

		h.mutexChainsMutex.Lock()
		h.chainsMutex[hashKey] = chain.mutex
		h.mutexChainsMutex.Unlock()

		fmt.Printf("\n Setter chain Element:%#v", chainElement)
		fmt.Printf("\n Add chain:%#v \n", chain)
		return
	}

	h.mutexChainsMutex.RLock()
	chainMutex = h.chainsMutex[hashKey]
	h.mutexChainsMutex.RUnlock()

	if chainMutex == nil {
		return
	}

	chainMutex.Lock()


	chainElement = chain.findElement(key)

	if chainElement != nil {
		deltaBytes = int64(chainElement.setValue(setOperationObject, value, expDate))
		hasAdded = false
	} else {
		chainElement = NewChainElement(key)
		chainElement.setValue(setOperationObject, value, expDate)
		chain.AddElement(chainElement)

		deltaBytes = int64(chainElement.valueByteSize) + int64(unsafe.Sizeof(chain)) + int64(unsafe.Sizeof(chainElement))
		hasAdded = false
	}

	chainLenght = chain.countElements

	chainMutex.Unlock()

	fmt.Printf("\n Setter chain Element:%#v", chainElement)
	return
}

func (h *HashTable) getElement(key string, getOperationObject IGetterValue) (isHit bool, chainElement *ChainElement) {
	var (
		chain      *Chain
		chainMutex *sync.RWMutex
	)
	h.reHashingMutex.RLock()
	defer h.reHashingMutex.RUnlock()

	fmt.Printf(`Hash function "%#v`, h.hashFunction)
	hashKey := h.hashFunction.CalculateHash(key)

	fmt.Printf(`GET Hash key: "%d"`, hashKey)
	h.mutexChains.RLock()
	chain = h.chains[hashKey]
	h.mutexChains.RUnlock()

	fmt.Printf("\n GET CHAIN:%#v \n", chain)

	if chain == nil {
		fmt.Printf("\n Can't find chain")
		return
	}

	h.mutexChainsMutex.RLock()
	chainMutex = h.chainsMutex[hashKey]
	h.mutexChainsMutex.RUnlock()

	if chainMutex == nil {
		fmt.Printf("\n Can't find chain mutex")
		return
	}

	chainMutex.RLock()

	fmt.Printf("\n Go throw elements chain %#v", chain)

	chainElement = chain.findElement(key)

	if chainElement != nil {
		chainElement.getValue(getOperationObject)
		isHit = true
	}

	chainMutex.RUnlock()
	return
}

func (h *HashTable) deleteElement(key string) (deletedBytes int64) {
	var (
		chain      *Chain
		chainMutex *sync.RWMutex
	)
	h.reHashingMutex.RLock()
	defer h.reHashingMutex.RUnlock()

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
		deletedBytes = int64(chainElement.valueByteSize) + int64(unsafe.Sizeof(chainElement))
		chain.deleteElement(chainElement)
	}

	chainMutex.RUnlock()
	return
}

func (h *HashTable) lruEviction() {
	bytesSize := atomic.LoadInt64(&h.bytesSize)
	fmt.Printf("\n bytesSize:%d maxBytesSize:%d \n", bytesSize, h.maxBytesSize)
	if h.maxBytesSize > bytesSize {
		return
	}

	h.reHashingMutex.RLock()
	defer h.reHashingMutex.RUnlock()

	// We will free 10% of maximumSizeBytes
	needFreeSize := bytesSize - (h.maxBytesSize - h.maxBytesSize/10)

	var freeBytes int64
	var element *ChainElement
	var countElements int64
	for element = h.lruChain.lastElement; element != nil; element = element.prevLRU {
		freeBytes += h.removeElement(element)
		countElements++
		if freeBytes >= needFreeSize {
			break
		}
	}

	atomic.AddInt64(&h.bytesSize, -1*freeBytes)
	atomic.AddInt64(&h.countElements, -1*countElements)
}

func (h *HashTable) reHashing() {
	countChains := atomic.LoadInt64(&h.capacity)
	countElements := atomic.LoadInt64(&h.countElements)
	bytesSize := atomic.LoadInt64(&h.bytesSize)

	if float64(countElements)/float64(countChains) <= 0.9 {
		return
	}

	newBytesSize := bytesSize + (int64(unsafe.Sizeof(&Chain{}))+int64(unsafe.Sizeof(&sync.RWMutex{})))*countChains

	// If we can't use more memory
	if h.maxBytesSize <= newBytesSize {
		return
	}

	timeStart := time.Now()
	h.reHashingMutex.Lock()

	newCapacity := countChains * 2
	newChains := make([]*Chain, newCapacity)
	newMutexes := make([]*sync.RWMutex, newCapacity)
	newHashFunction := NewComplexStringHash(uint64(newCapacity), h.coefPString, h.coefPInt)

	var nextChainElement *ChainElement
	for _, chain := range h.chains {
		if chain == nil {
			continue
		}
		for chainElement := chain.firstElement; chainElement != nil; chainElement = nextChainElement {
			nextChainElement = chainElement.next
			// We will use thist element in the new chain
			chainElement.next = nil
			h.reHachingAddElement(newChains, newHashFunction, chainElement)
		}
	}

	h.capacity = newCapacity
	h.chainsMutex = newMutexes
	h.hashFunction = newHashFunction
	h.bytesSize = newBytesSize

	h.reHashingTimeMetric.Observe(getDurationMicroseconds(time.Since(timeStart)))
	h.bytesSizeMetric.Set(float64(newBytesSize))

	h.reHashingMutex.Unlock()
}

func (h *HashTable) reHachingAddElement(chains []*Chain, hashFunction *ComplexStringHash, chainElement *ChainElement) {
	var (
		chain *Chain
	)
	hashKey := hashFunction.CalculateHash(chainElement.key)

	chain = chains[hashKey]

	if chain == nil {
		chain = NewChain(&sync.RWMutex{})
		chain.AddElement(chainElement)
		h.chains[hashKey] = chain
		return
	}

	chain.AddElement(chainElement)
}
