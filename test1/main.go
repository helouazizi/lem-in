// test/main.go
package main

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	Name string
	X, Y string
}
type Farm struct {
	Rooms              map[string]*Room
	Links              map[string][]string
	StartNeighbots     []string
	StartRoom, EndRoom string
	Ants               int
}

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
		result = append(result, current)
		if current == end {
			test = append(test, result)
			queue.Init()
			// queue.PushFront(start)
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

func (F *Farm) ReadFile(fileName string) (*Graph, error) {
	// open the file
	var err error
	graph := NewGraph()
	exist := 0
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// read the file by using the buffio pkg
	// that can give us convenient way to read input from a file
	// line by line using the  function scan()
	// befor looping lets inisialise our maps
	if F.Rooms == nil {
		F.Rooms = make(map[string]*Room)
	}
	if F.Links == nil {
		F.Links = make(map[string][]string)
	}
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {

		line := scanner.Text()
		line = strings.TrimSpace(line)
		// lets check if the first line is the valid number off ants
		if i == 0 {
			F.Ants, err = strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
			if F.Ants <= 0 {
				return nil, errors.New("invalid ants number")
			}
			i++
			continue
		}
		if i == 2 {
			check := strings.Split(line, " ")
			F.StartRoom = check[0]

			i = 1

		}

		if i == 3 {
			check := strings.Split(line, " ")
			F.EndRoom = check[0]

			i = 1

		}

		if line == "##start" {
			i = 2
			exist++
			///F.Rooms["##start"] = Room{Name: "", X: "", Y: ""}
			continue
		}
		if line == "##end" {
			i = 3
			exist += 2
			// F.Rooms["##end"] = Room{Name: "", X: "", Y: ""}
			continue
		}
		if line == "" || (line[0] == '#' && line != "##start" && line != "##end") {
			continue
		}
		check := strings.Split(line, " ")
		if len(check) == 3 {
			_, exist := F.Rooms[check[0]]
			if !exist {
				F.Rooms[check[0]] = &Room{X: check[1], Y: check[2]}
			} else {
				return nil, errors.New("found Duplicated rooms")
			}

		} else if len(check) == 1 {
			link := strings.Split(line, "-")
			if len(link) != 2 {
				fmt.Println(line)
				return nil, errors.New("no valid link found")

			} 
			_, exist := F.Rooms[link[0]]
			if !exist {
				fmt.Println(line)
				return nil, errors.New("no valid link found")
			}
			_, exist1 := F.Rooms[link[1]]
			if !exist1 {
				fmt.Println(line)
				return nil, errors.New("no valid link found")
			}
			F.Links[link[0]] = append(F.Links[link[0]], link[1])
			F.Links[link[1]] = append(F.Links[link[1]], link[0])
			graph.adjacencyList[link[0]] = graph.AddEdge(link[0],link[1])

			// graph.Add_Edges(link[1],link[0])

		} else {
			continue
		}

	}
	if exist != 3 {
		return nil, errors.New("no start or end room")
	}

	return graph, nil
}

func () createGraph() {

}

func main() {
	farm := &Farm{}
	graph, err := farm.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
}
