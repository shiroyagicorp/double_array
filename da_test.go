package double_array

import (
	"reflect"
	"testing"
)

func getData() []Item {
	data := []string{
		"abc",
		"abcd",
		"文字",
		"全角",
		"全角文字",
		"x",
		"y",
		"z",
		"xyzabc",
		"good",
		"漢字",
	}

	ret := make([]Item, len(data))

	for i, item := range data {
		ret[i] = Item(item)
	}

	return ret
}

func TestDoubleArray_Lookup(t *testing.T) {
	data := getData()

	da, err := NewDoubleArray(data)
	if err != nil {
		t.Error(err)
	}

	testData := []struct {
		item   string
		exists bool
	}{
		{
			item:   "ab",
			exists: false,
		},
		{
			item:   "bc",
			exists: false,
		},
		{
			item:   "abc",
			exists: true,
		},
		{
			item:   "abc",
			exists: true,
		},
		{
			item:   "漢字",
			exists: true,
		},
		{
			item:   "ひらがな",
			exists: false,
		},
	}

	inverse := ToInverseID(da)

	for _, tt := range testData {
		itemID := da.Lookup([]rune(tt.item))

		if !tt.exists {
			if itemID != ItemNotFound {
				t.Errorf("Item wrongly found: %s", tt.item)
			}
			continue
		}

		if itemID == ItemNotFound {
			t.Errorf("Item not found: %s", tt.item)
			continue
		}

		deserialized := Deserialize(da, itemID, inverse)
		if deserialized != tt.item {
			t.Errorf("Deserialization failed: expected %s, actual: %s", tt.item, deserialized)
		}
	}
}

func TestDoubleArray_Scan(t *testing.T) {
	data := getData()

	da, err := NewDoubleArray(data)
	if err != nil {
		t.Error(err)
	}

	testData := []struct {
		text     string
		expected []string
	}{
		{
			text:     "昔は全角文字が表示できないコンピューターも多かった。",
			expected: []string{"全角文字"},
		},
		{
			text:     "昔は全角文字が表示できないコンピューターも多かった。文字の形を保存しておくためのメモリが不足していたためだ。",
			expected: []string{"文字", "全角文字"},
		},
	}

	inverse := ToInverseID(da)

	for _, tt := range testData {
		actual := make([]string, 0)
		da.Scan([]rune(tt.text), func(i, j int, id ItemID) {
			deserialized := Deserialize(da, id, inverse)
			actual = append(actual, deserialized)
		})

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("failed to extract entries: expected: %v, actual %v", tt.expected, actual)
		}
	}
}
