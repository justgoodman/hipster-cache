package hash_table

import (
	"fmt"
	"math/rand"
)

// Polinomial Hash Function
type StringHash struct {
	// Coefficent which we can take randomally for Polimonial Family
	x uint64
	// Cardinality of the hash function, should be more than m*L (L is Maximum lenght of the string)
	p uint64
}

// Uniform family integer hash Function
type IntegerHash struct {
	// Coefficient which we can take randomally for Universal Family
	a uint64
	// Coefficient which we can take randomally for Universal Family
	b uint64
	// Cardinality of the hash function
	m uint64
	// Big prime number, must be more than function parameter
	p uint64
}

// Complex Hash function used for hash tables, we calculate value which will less than m(size of hash table)
// Using thus fuction for hash table the average lenght of the longest chain c is O(1+alpha)
// where aplha=n/m is the load factor of the hash table
type ComplexStringHash struct {
	// Cardinality of the hash function (size of hash table)
	m uint64
	// Cardinality of the hash string function, should be more than m*L (L is Maximum lenght of the string)
	pString uint64
	// Big prime number, must be more than function parameter i.e. more than pString
	pInteger    uint64
	stringHash  *StringHash
	integerHash *IntegerHash
}

// Calculate String Hash, using thus fuction for hash table  the average lenght of the longest chain c is O(1+alpha)
// where aplha=n/m is the load factor of the hash table
func (c *ComplexStringHash) CalculateHash(value string) uint64 {
	fmt.Printf(`Enter value: "%s"`, value)
	fmt.Printf(`String hash function: "%#v"`, c.stringHash)
	resultStringHash := c.stringHash.CalculateHash(value)
	fmt.Printf(`StringHash: "%d"`, resultStringHash)
	fmt.Printf("\n IntegerHash: %#v \n", c.integerHash)
	return c.integerHash.CalculateHash(resultStringHash)
}

// Create new Compmplex String Function which can be used in hash tables
// m - Cardinality of the hash function (size of hash table)
// pString - Cardinality of the hash string function, should be more than m*L (L is Maximum lenght of the string)
// pInteger - Big prime number, must be more than function parameter i.e. more than pString
func NewComplexStringHash(m, pString, pInteger uint64) *ComplexStringHash {
	return &ComplexStringHash{
		m:           m,
		pString:     pString,
		pInteger:    pInteger,
		stringHash:  NewStringHash(pString),
		integerHash: NewIntegerHash(pInteger, m),
	}
}

// Create new Polinomial function, x gives randomly from [1,p-1]
// p - Cardinality of the hash function, should be more than m*L (L is Maximum lenght of the string)
func NewStringHash(p uint64) *StringHash {
	return &StringHash{x: getRandom(p), p: p}
}

// Get random value from 1 to maxValue - 1
func getRandom(maxValue uint64) uint64 {
	return uint64(rand.Int63n(int64(maxValue-2))) + 1
}

// Create new Univeral Family Hash integer hash function, a,b gives randomly form [1,p-1]
// p - Big prime number, must be more than function parameter
// m - Cardinality of the hash function
func NewIntegerHash(p, m uint64) *IntegerHash {
	return &IntegerHash{a: getRandom(p), b: getRandom(p), m: m, p: p}
}

// Calculate Polinomial Hash function for string this function is from Universal Family
func (s *StringHash) CalculateHash(value string) uint64 {
	var result uint64
	var coef uint64
	coef = 1
	for _, character := range value {
		result += uint64(character) * coef % s.p
		fmt.Printf(`"\n" Result: "%d"`, result)
		coef *= s.x
	}
	return result % s.p
}

// Calculate Universal Family Hash function for integers
func (i *IntegerHash) CalculateHash(value uint64) uint64 {
	return ((i.a*value + i.b) % i.p) % i.m
}
