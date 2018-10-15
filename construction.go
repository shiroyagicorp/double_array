package double_array

import (
	"errors"
	"fmt"
)

func findSmallestCandidate(da *doubleArray, branch int32) int32 {
	candidate := -da.nodes[0].Check

	// base > 0
	for candidate-branch <= 0 {
		candidate = -da.nodes[candidate].Check
	}

	return candidate
}

func getContextIndex(da *doubleArray, context []int32) (int32, error) {
	index := int32(0)
	for _, c := range context {
		index = da.traverse(index, c)
		if index < 0 {
			return 0, errors.New("double array may be broken")
		}
	}

	return index, nil
}

func popNodeFromDoublyLinkedList(da *doubleArray, i int32) *doubleArray {
	pre := -da.nodes[i].Base
	post := -da.nodes[i].Check

	da.nodes[pre].Check = -post
	if int32(len(da.nodes)) <= post {
		da = extendDoubleArray(da)
	}
	da.nodes[post].Base = -pre

	return da
}

func insert(da *doubleArray, context []int32, targets []int32, exists []bool, n, tail int32) (*doubleArray, int32, error) {
	for i := int32(1); i < n; i++ {
		if targets[i-1] >= targets[i] {
			return nil, 0, errors.New("data is not sorted")
		}
	}

	index, err := getContextIndex(da, context)
	if err != nil {
		return nil, 0, err
	}

	candidate := findSmallestCandidate(da, targets[0])
	base := candidate - targets[0]

	i := int32(1)
	for {
		next := base + targets[i]

		if next >= int32(len(da.nodes)) {
			break
		}

		if candidate >= int32(len(da.nodes)) {
			break
		}

		if da.nodes[next].Check >= 0 {
			candidate = -da.nodes[candidate].Check
			base = candidate - targets[0]
			i = 1
		} else {
			i++
			if i >= n {
				break
			}
		}
	}

	maxIndex := base + targets[n-1]
	if maxIndex >= int32(len(da.nodes)) {
		da = extendDoubleArray(da)
	}

	if maxIndex > tail {
		tail = maxIndex
	}

	if da.nodes[index].Base < 0 {
		da.nodes[index].Base = -base
	} else {
		da.nodes[index].Base = base
	}

	for i := int32(0); i < n; i++ {
		next := base + targets[i]

		da = popNodeFromDoublyLinkedList(da, next)

		da.nodes[next].Check = index

		if exists[i] {
			da.nodes[next].Base = -1
		} else {
			da.nodes[next].Base = 0
		}
	}

	return da, tail, nil
}

func extendDoubleArray(da *doubleArray) *doubleArray {
	max := len(da.nodes)
	da.nodes = append(da.nodes, make([]node, len(da.nodes))...)
	initDoubleArray(da, int32(max))
	return da
}

func initDoubleArray(da *doubleArray, after int32) {
	if after == 0 {
		// set a root node
		da.nodes[0].Base = 0
		da.nodes[0].Check = -1
		after = 1
	}

	// build doubly linked list
	for i := after; i < int32(len(da.nodes)); i++ {
		// base refers pre-empty node index
		da.nodes[i].Base = -(i - 1)

		// check refers post-empty node index
		da.nodes[i].Check = -(i + 1)
	}
}

func constructDA(da *doubleArray, data []idsValue, maxRuneID int) (int32, error) {
	initDoubleArray(da, 0)

	var err error
	var previous *idsValue
	var targets = make([]int32, maxRuneID)
	var exists = make([]bool, maxRuneID)
	tail := int32(0)
	head := int32(0)
	nSingleNode := 0

	fmt.Println("...entering build loop...")
	progress := 0.1
	for i := range data {
		doneRatio := float64(i) / float64(len(data))
		if doneRatio >= progress {
			fmt.Printf("%d%% [%d / %d][%d]\n", int(doneRatio*100), i, len(data), tail+1)
			progress += 0.1
		}

		if previous != nil && previous.compareContext(&data[i]) != 0 {
			if head == 1 {
				nSingleNode++
			}

			da, tail, err = insert(da, previous.getContext(), targets, exists, head, tail)
			if err != nil {
				return 0, err
			}

			head = 0
		}

		target := data[i].getLeaf()
		if target == 0 {
			return 0, errors.New(fmt.Sprintf("invalid target string on index %d", i))
		}

		targets[head] = target
		exists[head] = data[i].exists
		head++
		previous = &data[i]
	}

	da, tail, err = insert(da, previous.getContext(), targets, exists, head, tail)
	if err != nil {
		return 0, err
	}

	fmt.Println("...done...")
	fmt.Printf("# of valid nodes: %d, length of DA: %d, # of single node: %d\n", len(data), tail+1, nSingleNode)

	return tail, nil
}
