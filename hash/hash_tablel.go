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
	mutexChains *sync.RWMutex
	chainsMutex []*sync.Mutex
	mutexChainsMutex *sync.RWMutex

	countElementsMetric *prometheus.Gauge
	bytesSizeMetric *prometheus.Gauge
        maxChainLenghtMetic *prometheus.Gauge
	responseTimeMetric *prometheus.SummaryVec
}

// For implement get operation for all types data you can implement this interface 
type IGetterValue interface {
     func GetValue(value interface{})
}

// For implement set operation for all types data you can implemet this interface
type ISetterValue interface {
     func SetValue(sourceValue interface{}, newValue interface{}) (valueSizeBytes int)
}

const (
	namespace = "hispter_cache"
	maximumLoadFactor = 0.9
)

func NewHashTable() *HashTable {
	return &HashTable{}
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
		}, []string{"operation"})

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

func (this *HashTable) SetElement(key string,value interface{}, setterValue ISetterValue) {
	timeStart := time.Now()

	hasAdded, chainLenght, deltaBytes := this.setElement(key,value, setterValue)
	if hadAdded {
		countElements := atomic.AddUint64(&this.countElements,1)
		this.countElementsMetric.Set(float64(countElemenets))
	}
	bytesSize := atomic.AddUint64(&this.bytesSize,deltaBytes)
	this.bytesSizeMetric.Set(float64(butesSize))
	maxChainLenght = atomic.LoadUint64(&this.maxChainLenght)
	if maxChainLenght < chainLenght {
		atomic.StoreUint64(*this.maxChainLenght, chainLenght)
		this.maxChainLenghtMetric.Set(float64(chainLenght))
	}

	responseMetric := this.responseTimeMetric.WithLabelValues("set")
	responseMetric.Observe(getDurationMicroseconds(time.Since(timeStart))
}

func (this *HashTable) GetElement(key string,value interface{}, getterValue IGetterValue) {
	timeStart := time.Now()

	this.getElement(key,getterValue)

	if this.getElement(key,getterValue) {
		this.hitMetric.Inc()
	}
	else {
		this.missMetric.Inc()
	}

	responseMetric := this.responseTimeMetric.WithLabelValues("get")
	responseMetric.Observe(getDurationMicroseconds(time.Since(timeStart))
}


// Return delta byte size
func (this *HashTable) setElement(key string, setOpertionObject ISetOperation, value interface) (hasAdded bool, chainLenght int, deltaBytes int) {
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
		chainElement := NewChainElement(key)
		chainElement.SetValue(setOperationObject,value)
		chain := NewChain(chainElement)

		deltaBytes = chainElement.byteSize + unsafe.Sizeof(chain)
		hasAdded = true
		chainLeight = 1

		this.mutexChain.Lock()
		this.chains[hashKey] := chain
		this.mutexChain.Unlock()
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

	element := chain.findElement(key)

	if element != nil {
		deltaBytes = element.setValue(value)
		hasAdded = false
	} else {
		chainElement := NewChainElement(setOperationObject,value)
		chain.AddElement(chainElement)

		deltaBytes = chainElement.sizeButes
		hasAdded = false
	}

	chainLenght = chain.lenght

	chainMutex.Unlock()

	return
}


func (this *HashTable) getElement(key string, getOpertionObject IGetOperation, value interface) (isHit bool) {
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

	element := chain.findElement(key)
	isHit = true

	if element != nil {
		deltaBytes = element.getValue(getOperationObject)
	}

	chainLenght = chain.lenght

	chainMutex.RUnlock()

	return
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
