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
}

//  get duration in microsecond for sub: time1-time2
func getDurationMicroseconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Microsecond)
}

func (this *HashTable) SetElement(key string,value interface) {
	timeStart := time.Now()

	hasAdded, chainLenght, deltaBytes := this.setElement(key,value)
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

func (this *HashTable) appendList(key, value string) error {
	// Find Element
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
		return fmt.Errorf(`Can't find element by key`)
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
		return fmt.Errorf(`Can't find element by key`)
	}

	chainLenght = chain.lenght



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
		chainElement := NewChainElement(key,value)
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
		chainElement := NewChainElement(key,value)
		chain.AddElement(chainElement)

		deltaBytes = chainElement.sizeButes
		hasAdded = false
	}

	chainLenght = chain.lenght

	chainMutex.Unlock()

	return
}




