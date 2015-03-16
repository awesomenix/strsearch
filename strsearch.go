// Substring search algorithm
// By Leonid Volnitsky
//paper reference http://volnitsky.com/project/str_search/

package strsearch

import (
	"bytes"
	"errors"
	"math"
)

//searches substring in string str

func Search(str, substr []byte) (int, error) {
	strlen := len(str)

	if strlen <= 0 {
		return -1, errors.New("Invalid String Specified")
	}

	substrlen := len(substr)

	if substrlen <= 0 {
		return -1, errors.New("Invalid Search String Specified")
	}

	wordSize := 2 // skip 2 bytes at a time
	step := substrlen - wordSize + 1

	// args sizes check
	// see startup costs;  algo limit: 2*SS_size
	if strlen < 1000 ||
		substrlen < (2*wordSize-1) ||
		substrlen >= math.MaxInt32 {
		return bytes.Index(str, substr), nil //fallback to default search
	}

	// prefill hash table with substr step values
	hashSize := 64 * 1024
	hashTable := make([]int, hashSize)
	for i := substrlen - wordSize; i >= 0; i-- {
		h := int(substr[i]) % hashSize
		for hashTable[h] != 0 {
			h = (h + 1) % hashSize // find free cell
		}
		hashTable[h] = i + 1
	}

	// step through text
	probeVal := substrlen - wordSize
	for ; probeVal <= strlen-substrlen; probeVal += step {
		//skip substr chars specified by hash table value
		for h := int(str[probeVal]) % hashSize; hashTable[h] != 0; h = (h + 1) % hashSize {
			subStrHit := str[probeVal-(hashTable[h]-1):]
			match := true
			//Scan through the substrlen
			for i := 0; i < substrlen; i++ {
				if subStrHit[i] != substr[i] {
					match = false
					break
				}
			}
			if match {
				return probeVal - (hashTable[h] - 1), nil // found
			}
		}
	}

	// check tail for possible match, may happen if the step size is large
	return bytes.Index(str[probeVal-step+1:], substr), nil
}
