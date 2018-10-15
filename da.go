package double_array

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/vmihailenco/msgpack"
)

type ItemID int32

func Deserialize(obj DoubleArray, id ItemID, inverse InverseID) string {
	da := obj.(*doubleArray)

	rs := make([]rune, 0)
	index := int32(id)
	for index != 0 {
		next := da.nodes[index].Check
		base := da.nodes[next].Base
		if base < 0 {
			base = -base
		}

		runeID := index - base
		r := inverse[runeID]
		rs = append(rs, r)

		index = next
	}

	return string(rs)
}

const (
	ItemNotFound = ItemID(-1)
)

type node struct {
	Base  int32
	Check int32
}

type doubleArray struct {
	nodes []node
	dict  RuneID
}

func (da *doubleArray) traverse(index int32, branch int32) int32 {
	base := da.nodes[index].Base
	if base < 0 {
		base = -base
	}
	next := base + branch
	if next >= int32(len(da.nodes)) || da.nodes[next].Check != index {
		return -1
	} else {
		return next
	}
}

func (da *doubleArray) searchMaximumLengthMatchFromSuffix(text []rune) (ItemID, int) {
	index := int32(0)
	currentMaximumMatch := ItemNotFound
	lastPos := -1
	for i := len(text) - 1; i >= 0; i-- {
		branch, ok := da.dict[text[i]]
		if !ok {
			return currentMaximumMatch, lastPos
		}

		next := da.traverse(index, branch)
		if next < 0 {
			return currentMaximumMatch, lastPos
		}

		if da.nodes[next].Base < 0 {
			currentMaximumMatch = ItemID(next)
			lastPos = i
		}

		index = next
	}

	return currentMaximumMatch, lastPos
}

func (da *doubleArray) Lookup(key []rune) ItemID {
	index := int32(0)
	for i := len(key) - 1; i >= 0; i-- {
		branch, ok := da.dict[key[i]]
		if !ok {
			return ItemNotFound
		}

		next := da.traverse(index, branch)
		if next < 0 {
			return ItemNotFound
		}

		index = next
	}

	if da.nodes[index].Base < 0 {
		return ItemID(index)
	}

	return ItemNotFound
}

func (da *doubleArray) Scan(text []rune, callback func(i, j int, id ItemID)) {
	for i := len(text) - 1; i > 0; i-- {
		id, pos := da.searchMaximumLengthMatchFromSuffix(text[:i])
		if id >= 0 {
			callback(pos, i, id)
			i = pos
		}
	}
}

func (da *doubleArray) Serialize() ([]byte, error) {
	serializable := da2serializableDA(da)
	data, err := msgpack.Marshal(serializable)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(data)
	writer.Flush()
	writer.Close()
	return buf.Bytes(), nil
}

type DoubleArray interface {
	Serialize() ([]byte, error)
	Lookup(key []rune) ItemID
	Scan(text []rune, callback func(i, j int, id ItemID))
}

func NewDoubleArrayFromBytes(data []byte) (DoubleArray, error) {
	buf := bytes.NewBuffer(data)
	reader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	rawData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var tda serializableDA
	err = msgpack.Unmarshal(rawData, &tda)
	if err != nil {
		return nil, err
	}

	da := serializableDA2DA(&tda)
	return da, err
}

func NewDoubleArray(data []Item) (DoubleArray, error) {
	ivs, dict, err := preprocess(data)
	if err != nil {
		return nil, err
	}

	doubleArray := &doubleArray{
		nodes: make([]node, len(data) * 2),
		dict:  dict,
	}
	tail, err := constructDA(doubleArray, ivs, len(dict)+2)
	if err != nil {
		return nil, err
	}

	doubleArray.nodes = doubleArray.nodes[:tail+1]

	return doubleArray, nil
}
