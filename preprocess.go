package double_array

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

func compareInt32s(a, b []int32) int {
	shorter := len(a)
	if len(b) < shorter {
		shorter = len(b)
	}

	for i := 0; i < shorter; i++ {
		if a[i] != b[i] {
			return int(a[i] - b[i])
		}
	}

	return len(a) - len(b)
}

func sortLexicographicalOrder(data []idsValue) {
	fmt.Print("sort in lexicographical order...")
	sort.Slice(data, func(i, j int) bool {
		return compareInt32s(data[i].ids, data[j].ids) < 0
	})
	fmt.Println("done")
}

func sortTrieOrder(data []idsValue) {
	fmt.Print("sort in trie insertion order...")
	sort.Slice(data, func(i, j int) bool {
		cmp := data[i].compareContext(&data[j])
		if cmp == 0 {
			return data[i].getLeaf() < data[j].getLeaf()
		} else {
			return cmp < 0
		}
	})
	fmt.Println("done")
}

func diff(first, second []int32) int {
	shorter := len(first)
	if len(second) < shorter {
		shorter = len(second)
	}

	for i := 0; i < shorter; i++ {
		if first[i] != second[i] {
			return i
		}
	}

	return shorter
}

func fill(data []idsValue) ([]idsValue, error) {
	fmt.Print("fill lacking nodes")
	filled := make([]idsValue, 0, len(data)*2)

	previousWord := make([]int32, 0)
	for i := range data {
		word := data[i].ids

		if reflect.DeepEqual(word, previousWord) {
			return nil, errors.New("duplicate entry")
		}

		d := diff(previousWord, word)
		for j := d; j < len(word)-1; j++ {
			filled = append(filled, idsValue{word[:j+1], false})
		}
		filled = append(filled, data[i])

		previousWord = word

		if i%10000 == 0 {
			fmt.Print(".")
		}
	}
	fmt.Printf("done [count=%d, cap=%d]\n", len(filled), cap(filled))

	return filled, nil
}

func reverseKeys(data []Item, dict RuneID) []idsValue {
	ret := make([]idsValue, len(data))
	for i := range data {
		key := data[i]
		rev := make([]rune, len(key))
		for j := range key {
			rev[j] = key[len(key)-j-1]
		}

		ret[i] = Item(rev).toIDsValue(dict, true)
	}

	return ret
}

func preprocess(data []Item) ([]idsValue, RuneID, error) {
	dict := buildDict(data)

	ivs := reverseKeys(data, dict)
	sortLexicographicalOrder(ivs)
	ivs, err := fill(ivs)
	if err != nil {
		return nil, nil, err
	}
	sortTrieOrder(ivs)

	return ivs, dict, nil
}
