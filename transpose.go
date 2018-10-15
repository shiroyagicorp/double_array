package double_array

type serializableDA struct {
	Base  []int32
	Check []int32
	Dict  RuneID
}

func da2serializableDA(da *doubleArray) *serializableDA {
	data := serializableDA{
		Base:  make([]int32, len(da.nodes)),
		Check: make([]int32, len(da.nodes)),
		Dict:  da.dict,
	}

	for i, n := range da.nodes {
		data.Base[i] = n.Base
		data.Check[i] = n.Check
	}

	return &data
}

func serializableDA2DA(tda *serializableDA) *doubleArray {
	da := &doubleArray{
		nodes: make([]node, len(tda.Base)),
		dict: tda.Dict,
	}
	for i := range tda.Base {
		da.nodes[i] = node{
			Base:  tda.Base[i],
			Check: tda.Check[i],
		}
	}
	return da
}
