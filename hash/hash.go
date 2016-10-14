package hash

import (
	"math/rand"
)

// Polinomial Hash Function
struct StringHash {
	// Coefficent which we can take randomally for Polimonial Family
	x int64
	// Cardinality of the hash function, should be more than m*L (L is Maximum lenght of the string)
	p int64
}

// Uniform family integer hash Function
struct IntegerHash {
	// Coefficient which we can take randomally for Universal Family
	a int64
	// Coefficient which we can take randomally for Universal Family
	b int64
	// Cardinality of the hash function
	m int64
	// Big prime number, must be more than function parameter
	p int64
}

// Complex Hash function used for hash tables, we calculate value which will less than m(size of hash table)
// Using thus fuction for hash table the average lenght of the longest chain c is O(1+alpha)
// where aplha=n/m is the load factor of the hash table
struct ComplexStringHash {
	// Cardinality of the hash function (size of hash table)
	m int64
	// Cardinality of the hash string function, should be more than m*L (L is Maximum lenght of the string)
	pString int64
	// Big prime number, must be more than function parameter i.e. more than pString 
        pInteger int64
	stringHash *StringHash
	integerHash *IntegerHash
}


// Calculate String Hash, using thus fuction for hash table  the average lenght of the longest chain c is O(1+alpha)
// where aplha=n/m is the load factor of the hash table
func (this *ComplexStringHash) CalculateHash(value string) int64 {
	resultStringHash := this.stringHash.CalculateHash(value)
	return this.integerHash(resultStringHash)
}

// Create new Compmplex String Function which can be used in hash tables
// m - Cardinality of the hash function (size of hash table)
// pString - Cardinality of the hash string function, should be more than m*L (L is Maximum lenght of the string)
// pInteger - Big prime number, must be more than function parameter i.e. more than pString
func NewComplexStringHash(m,pString,pInteger int64) *ComplexStringHash {
	return &ComplexStringHash{
			m:m,
			pString: pString,
			pInteger: pInteger,
			stringHash: NewStringHash(pString),
			integerHash: NewIntegerHash(pInteger,m),
	}


// Create new Polinomial function, x gives randomly from [1,p-1]
// p - Cardinality of the hash function, should be more than m*L (L is Maximum lenght of the string)
func NewStringHash(p int64) *StringHash {
	return &PolinomailHashing{x: getRandom(p),p: p,}
}

// Get random value from 1 to maxValue - 1
func getRandom(maxValue int64) int64 {
	return rand.Int63(maxValue-2) + 1
}

// Create new Univeral Family Hash integer hash function, a,b gives randomly form [1,p-1] 
// p - Big prime number, must be more than function parameter
// m - Cardinality of the hash function
func NewIntegerHash(p,m int int64) *IntegerHash {
	return &IntegerHashing{a: getRandom(p), b: getRandom(p), m: m,p: p,}



// Calculate Polinomial Hash function for string this function is from Universal Family 
func (this *PolinomialHashing) CalculateHash(value string) int64 {
	var coef := 1
	for _, character := range mixed {
		result += character*coef % this.p
		coef *= this.x
	}
	return result % this.p
}

// Calculate Universal Family Hash function for integers
func (this *IntegerHasing) CalculateHash(value int64) int64 {
	return ((this.a * value + this.b) % this.p) % this.m
}


