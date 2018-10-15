package double_array

type Item []rune

func (item Item) toIDsValue(dict RuneID, exists bool) idsValue {
	ids := make([]int32, len(item))
	for i := range item {
		ids[i] = dict[item[i]]
	}

	return idsValue{
		ids:    ids,
		exists: exists,
	}
}

type idsValue struct {
	ids    []int32
	exists bool
}

func (iv *idsValue) getContext() []int32 {
	return iv.ids[:len(iv.ids)-1]
}

func (iv *idsValue) compareContext(other *idsValue) int {
	context1 := iv.getContext()
	context2 := other.getContext()
	return compareInt32s(context1, context2)
}

func (iv *idsValue) getLeaf() int32 {
	// if kv is a leaf node
	return iv.ids[len(iv.ids)-1]
}
