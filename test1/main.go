// test/main.go
package main

import (
	"container/list"
	"fmt"
)

// Graph represents an undirected graph
type Graph struct {
	adjacencyList map[string]*list.List
}

// NewGraph creates a new empty graph
func NewGraph() *Graph {
	return &Graph{
		adjacencyList: make(map[string]*list.List),
	}
}

// AddVertex adds a new vertex to the graph
func (g *Graph) AddVertex(v string) {
	g.adjacencyList[v] = list.New()
}

// AddEdge adds an edge between two vertices
func (g *Graph) AddEdge(v1, v2 string) {
	if g.adjacencyList[v1] == nil {
		g.adjacencyList[v1] = list.New()
	}
	if g.adjacencyList[v2] == nil {
		g.adjacencyList[v2] = list.New()
	}
	g.adjacencyList[v1].PushBack(v2)
	g.adjacencyList[v2].PushBack(v1)
}

// BFS performs a breadth-first search starting from the given vertex
func BFS(g *Graph, start, end string) [][]string {
	visited := make(map[string]bool)
	result := []string{}
	test := [][]string{}
	queue := list.New()

	// Mark the start node as visited and enqueue it
	visited[start] = true
	queue.PushFront(start)

	for queue.Len() > 0 {
		current := queue.Front().Value.(string)
		queue.Remove(queue.Front())

		//if len(result) == 0 || current != result[len(result)-1] {
		result = append(result, current)
		//}
		if current == end {
			test = append(test, result)
			queue.Init()
			//queue.PushFront(start)
			visited[start] = false
			visited[end] = false
			result = []string{}
			fmt.Println("here")
			current = start

		}

		// Dequeue all nodes adjacent to this node
		for neighborElement := g.adjacencyList[current].Front(); neighborElement != nil; neighborElement = neighborElement.Next() {
			neighbor := neighborElement.Value.(string)
			if !visited[neighbor] {
				visited[neighbor] = true
				queue.PushFront(neighbor)
			}
		}
	}

	return test
}

func main() {
	graph := NewGraph()
	graph.AddVertex("0")
	graph.AddVertex("1")
	graph.AddVertex("2")
	graph.AddVertex("3")
	graph.AddVertex("4")
	graph.AddVertex("5")

	graph.AddEdge("0", "1")
	graph.AddEdge("0", "2")
	graph.AddEdge("1", "3")
	graph.AddEdge("2", "4")
	graph.AddEdge("4", "5")
	graph.AddEdge("3", "5")

	fmt.Println("BFS traversal:", BFS(graph, "0", "5"))
}
