package sortedgraph

import (
	"velour/debruijn"
)

// ===================================
// SortedNode
// ===================================

type SortedNode struct {
	kmer			debruijn.Kmer
	frequency		int
	neighbors		uint8
}

func NewNode(kmer debruijn.Kmer) debruijn.GraphNode {
	var node debruijn.GraphNode = &SortedNode{kmer, 1, 0}
	return node
}

// ===================================
// SortedNode Functions
// ===================================

func (node *SortedNode) GetKmer() debruijn.Kmer {
	return node.kmer
}

func (node *SortedNode) GetFrequency() int {
	return node.frequency
}

func (node *SortedNode) SetFrequency(new_frequency int) {
	node.frequency = new_frequency
}

func (node *SortedNode) IncrementFrequency() {
	freq := node.GetFrequency()
	if freq != 255 {
		node.SetFrequency(freq + 1)
	}
}

func (node *SortedNode) GetPredecessors() []int {
	preds := make([]int, 4)

	var i uint8 = 3
	for int(i) - 3 < len(preds) {
		bit := 0x01 & (node.neighbors >> i)
		if bit == 1 {
			preds[i - 3] = 1
		}

		i++
	}

	return preds
}

func (node *SortedNode) IsAPredecessor(nt int) bool {
	bit := (node.neighbors >> uint8(3 + nt)) & 0x01

	if bit == 1 {
		return true
	} else {
		return false
	}
}

func (node *SortedNode) AddPredecessor(nt int) {
	var bit uint8

	if nt == 0 {
		bit = 0x10
	} else if nt == 1 {
		bit = 0x20
	} else if nt == 2 {
		bit = 0x40
	} else if nt == 3 {
		bit = 0x80
	}

	node.neighbors = node.neighbors | bit
}

func (node *SortedNode) GetSuccessors() []int {
	succs := make([]int, 4)

	var i uint8 = 0
	for int(i) < len(succs) {
		bit := 0x01 & (node.neighbors >> i)
		if bit == 1 {
			succs[i] = 1
		}

		i++
	}

	return succs
}

func (node *SortedNode) IsASuccessor(nt int) bool {
	bit := (node.neighbors >> uint8(nt)) & 0x01

	if bit == 1 {
		return true
	} else {
		return false
	}
}

func (node *SortedNode) AddSuccessor(nt int) {
	var bit uint8

	if nt == 0 {
		bit = 0x01
	} else if nt == 1 {
		bit = 0x02
	} else if nt == 2 {
		bit = 0x04
	} else if nt == 3 {
		bit = 0x08
	}

	node.neighbors = node.neighbors | bit
}

func (node *SortedNode) Merge(other_entry debruijn.GraphNode) {
	kmer_a := node.GetKmer()
	kmer_b := other_entry.GetKmer()

	if kmer_a.Equals(kmer_b) {
		node.SetFrequency(node.GetFrequency() + other_entry.GetFrequency())

		for _, nt := range other_entry.GetPredecessors() {
			node.AddPredecessor(nt)
		}

		for _, nt := range other_entry.GetSuccessors() {
			node.AddSuccessor(nt)
		}
	}
}
