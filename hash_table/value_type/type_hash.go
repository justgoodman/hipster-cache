package value_type

type HashGetOperation struct {
	// Enter parameter
	Key string
	// Outout parameter
	Value string
	Err   error
}

func (this *HashGetOperation) GetValue(value interface{}) {
	switch hashTable := value.(type) {
	case *HashTable:
		this.Value = hashTable.GetValue(this.Key)
	default:
		this.Err = fmt.Errorf("Hash table is not type of the hash table")
	}
}

type HashSetOperation struct {
	// Enter parameter
	Key string
	Err error
}

func (this *HashSetOperation) SetValue(sourceValue, value interface{}) (valueByteSize int) {
	var stringValue string
	switch val := value.(type) {
	case string:
		stringValue := val
	default:
		this.Err = fmt.Errorf("hash key is not the string type")
		return
	}
	switch hashTable := sourceValue.(type) {
	case *HashTable:
		hashTable.SetValue(this.Key, stringValue)

		valueByteSize = hashTable.ByteSize
	case nil:
		hashTable = NewHashTable()
		hashTable.Set(this.Key, stringValue)

		valueByteSize = hashTable.ByteSize
	default:
		this.Err = fmt.Errorf("Hash table is not type of the hash table")
	}
}
