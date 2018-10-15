package double_array

import (
	"sort"
)

type RuneID map[rune]int32
type InverseID map[int32]rune

func buildDict(data []Item) RuneID {
	counts := make(map[rune]int32)
	for _, item := range data {
		for _, r := range item {
			if _, ok := counts[r]; ok {
				counts[r]++
			} else {
				counts[r] = 1
			}
		}
	}

	type pair struct {
		r rune
		c int32
	}

	pairs := make([]pair, len(counts))
	i := 0
	for r, c := range counts {
		pairs[i] = pair{r, c}
		i++
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].c > pairs[j].c
	})

	dict := make(RuneID)
	for i := range pairs {
		dict[pairs[i].r] = int32(i + 1)
	}

	return dict
}

func ToInverseID(da DoubleArray) InverseID {
	dict := da.(*doubleArray).dict
	inverse := make(InverseID)
	for r, id := range dict {
		inverse[id] = r
	}

	return inverse
}
