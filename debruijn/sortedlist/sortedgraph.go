package sortedgraph

import (
	"sort"
	"velour/debruijn"
)

// ===================================
// SortedGraph
// ===================================

type SortedGraph struct {
	nodes			[]*debruijn.GraphNode
	newNode			debruijn.NodeGenerator
}

func NewGraph(newNode debruijn.NodeGenerator) debruijn.Graph {
	var graph debruijn.Graph = &SortedGraph{make([]*debruijn.GraphNode, 0, 3000000), newNode}
	return graph
}

// ===================================
// SortedGraph Functions
// ===================================

func (graph *SortedGraph) Len() int {
	return len(graph.nodes)
}

func (graph *SortedGraph) GetFrequencies() []int {
	freqs := make([]int, graph.Len())

	for i := range freqs {
		freqs[i] = (*graph.nodes[i]).GetFrequency()
	}

	return freqs
}

func (graph *SortedGraph) GetNumNodesSeen() int {
	num_seen := 0

	for _, freq := range graph.GetFrequencies() {
		num_seen += freq
	}

	return num_seen
}

func (graph *SortedGraph) NewNode(kmer debruijn.Kmer) debruijn.GraphNode {
	return graph.newNode(kmer)
}

func (graph *SortedGraph) GetNode(kmer debruijn.Kmer) (int, debruijn.GraphNode, bool) {
	var node debruijn.GraphNode

	n := graph.Len()
	i := sort.Search(n, func (i int) bool {
		other_kmer := (*graph.nodes[i]).GetKmer()
		return kmer.GetValue() <= other_kmer.GetValue()
	})

	if i == n {
		return -1, node, false
	} else if node = (*graph.nodes[i]); node.GetKmer() == kmer {
		return i, node, true
	} else {
		return i, node, false
	}
}

func (graph *SortedGraph) SetNode(kmer debruijn.Kmer, node debruijn.GraphNode) int {
	n := graph.Len()

	i := sort.Search(n, func (i int) bool {
		other_kmer := (*graph.nodes[i]).GetKmer()
		return kmer.GetValue() <= other_kmer.GetValue()
	})

	graph.SetNodeAtIndex(i, node)

	return i
}

func (graph *SortedGraph) SetNodeAtIndex(i int, node debruijn.GraphNode) int {
	if i >= 0 && i < graph.Len() {
		graph.nodes = append(graph.nodes, graph.nodes[graph.Len() - 1])
		copy(graph.nodes[i + 1:], graph.nodes[i : graph.Len() - 2])
		graph.nodes[i] = &node
	} else {
		i = graph.Len()
		graph.nodes = append(graph.nodes, &node)
	}

	return i
}

func (graph *SortedGraph) ConnectNodeToGraph(kmer debruijn.Kmer, kmer_ind int, node debruijn.GraphNode) {
	nts := [4]byte{'A', 'C', 'G', 'T'}

	for i, nt := range nts {
		prec_kmer := kmer.GeneratePredecessor(nt)

		if _, prec_node, ok := graph.GetNode(prec_kmer); ok {
			node.AddPredecessor(i)
			prec_node.AddSuccessor(kmer.GetLastNucleotide())
			break
		}
	}
}

func (graph *SortedGraph) AddNode(kmer debruijn.Kmer) int {
	var kmer_ind int
	var node debruijn.GraphNode
	var ok bool

	if kmer_ind, node, ok = graph.GetNode(kmer); ok {
		node.IncrementFrequency()
	} else {
		node = graph.newNode(kmer)
		kmer_ind = graph.SetNodeAtIndex(kmer_ind, node)
		// graph.ConnectNodeToGraph(kmer, kmer_ind, node)
	}

	return kmer_ind
}

func (graph *SortedGraph) AddNodes(kmers []debruijn.Kmer) []int {
	ids := make([]int, 0)

	for _, kmer := range kmers {
		ids = append(ids, graph.AddNode(kmer))
	}

	return ids
}
