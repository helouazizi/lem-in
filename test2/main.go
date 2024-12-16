// test2/main.go
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
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
	StartRoom, EndRoom string
	Ants               int
	FileSize           int64
}

// CreateAdjacencyMatrix creates an adjacency matrix from the Farm struct
func CreateAdjacencyMatrix(farm *Farm) ([][]bool, []string, int, int) {
	size := len(farm.Rooms)
	startIndex := 0
	endIndex := 0
	roomNames := make([]string, 0, size)
	for name := range farm.Rooms {
		// lets delete this room from the map
		// after we've added it to the matrix
		// we dont need it any more
		delete(farm.Rooms, name)
		roomNames = append(roomNames, name)
	}

	adjacencyMatrix := make([][]bool, size)
	for i := range adjacencyMatrix {
		adjacencyMatrix[i] = make([]bool, size)
	}

	for i, roomName := range roomNames {
		if roomName == farm.StartRoom {
			startIndex = i
		}
		if roomName == farm.EndRoom {
			endIndex = i
		}
		for _, linkedRoom := range farm.Links[roomName] {
			for j, linkedRoomName := range roomNames {
				if linkedRoomName == linkedRoom {
					adjacencyMatrix[i][j] = true
					break
				}
			}
		}
		// lets delete this room with it links from the link map
		// after we have processed it
		delete(farm.Links, roomName)
	}
	// lets clear the maps we dont need them any more
	//farm.Links = nil
	//farm.Rooms = nil
	return adjacencyMatrix, roomNames, startIndex, endIndex
}

func FindAllPaths(adjacencyMatrix [][]bool, start int, end int, roomNames []string, ants int) [][]string {
	var paths [][]string
	var currentPath []string

	visited := make([]bool, len(roomNames))
	//foundEnoughPaths := false // Flag to indicate if we found enough paths

	var dfs func(vertex int)
	dfs = func(vertex int) {
		visited[vertex] = true
		currentPath = append(currentPath, roomNames[vertex])

		if vertex == end {
			// Found a path to the end
			//fmt.Println(currentPath)
			paths = append(paths, append([]string(nil), currentPath...)) // Store a copy of the current path
			//fmt.Println("Found path:", currentPath)                      // Debug print
			//paths = append(paths, currentPath)
			// Check if we have found enough paths
			//if len(paths) == ants {
			//foundEnoughPaths = true // Set the flag to true
			//return                  // Stop searching if we have enough paths
			//}
		} else {
			// Explore neighbors
			for i := 0; i < len(adjacencyMatrix[vertex]); i++ {
				if adjacencyMatrix[vertex][i] && !visited[i] {
					dfs(i) // Recur for the next vertex
					//if foundEnoughPaths { // Check again after returning from recursion
					//return // Stop searching if we have enough paths
					//}
				}
			}
		}

		// Backtrack
		visited[vertex] = false
		currentPath = currentPath[:len(currentPath)-1]
	}

	dfs(start)
	return paths
}
func (F *Farm) ReadFile(fileName string) error {
	// open the file
	var err error

	fileinfo, err := os.Stat("test.txt")
	if err != nil {
		return err
	}
	F.FileSize = int64(fileinfo.Size())
	exist := 0
	file, err := os.Open(fileName)
	if err != nil {
		return err
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
				return err
			}
			if F.Ants <= 0 {
				return errors.New("invalid ants number")
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
				return errors.New("found Duplicated rooms")
			}

		} else if len(check) == 1 {
			link := strings.Split(line, "-")
			if len(link) != 2 {
				fmt.Println(line)
				return errors.New("no valid link found")

			}
			_, exist := F.Rooms[link[0]]
			if !exist {
				fmt.Println(line)
				return errors.New("no valid link found")
			}
			_, exist1 := F.Rooms[link[1]]
			if !exist1 {
				fmt.Println(line)
				return errors.New("no valid link found")
			}
			F.Links[link[0]] = append(F.Links[link[0]], link[1])
			F.Links[link[1]] = append(F.Links[link[1]], link[0])

			//graph.Add_Edges(link[1],link[0])

		} else {
			continue
		}

	}
	if exist != 3 {
		return errors.New("no start or end room")
	}

	return nil
}

func main() {
	// Parse the input data
	farm := &Farm{
		Rooms: make(map[string]*Room),
		Links: make(map[string][]string),
	}
	err := farm.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Create the adjacency matrix
	adjMatrix, roomNames, startIndex, endIndex := CreateAdjacencyMatrix(farm)

	// Find all paths
	paths := FindAllPaths(adjMatrix, startIndex, endIndex, roomNames, farm.Ants)

	// Sort paths by length
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

}
