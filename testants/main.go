// testants/main.go
// lem-in/main.go
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
	roomNames := make([]string, 0, size)
	roomIndex := make(map[string]int)

	for name := range farm.Rooms {
		roomNames = append(roomNames, name)
	}
	sort.Strings(roomNames) // Sort to ensure consistency

	for i, name := range roomNames {
		roomIndex[name] = i
	}

	adjacencyMatrix := make([][]bool, size)
	for i := range adjacencyMatrix {
		adjacencyMatrix[i] = make([]bool, size)
	}

	for room, neighbors := range farm.Links {
		for _, neighbor := range neighbors {
			i, j := roomIndex[room], roomIndex[neighbor]
			adjacencyMatrix[i][j] = true
			adjacencyMatrix[j][i] = true
		}
	}

	startIndex := roomIndex[farm.StartRoom]
	endIndex := roomIndex[farm.EndRoom]

	return adjacencyMatrix, roomNames, startIndex, endIndex
}

func FindAllPaths(adjacencyMatrix [][]bool, start, end int, roomNames []string) [][]string {
	var (
		paths       [][]string
		currentPath []string
		visited     = make([]bool, len(roomNames))
	)

	var dfs func(int)
	dfs = func(vertex int) {
		visited[vertex] = true
		currentPath = append(currentPath, roomNames[vertex])

		if vertex == end {
			paths = append(paths, append([]string(nil), currentPath...))
		} else {
			for i, isConnected := range adjacencyMatrix[vertex] {
				if isConnected && !visited[i] {
					dfs(i)
				}
			}
		}

		visited[vertex] = false
		currentPath = currentPath[:len(currentPath)-1]
	}

	dfs(start)
	return paths
}

func DistributeAnts(paths [][]string, ants int) []string {
	var result []string

	// Track progress and path assignments
	antPositions := make(map[int]int)      // Tracks which room each ant is in (progress along the path)
	antPaths := make(map[int]int)          // Tracks which path each ant is assigned to
	occupiedRooms := make(map[string]bool) // Tracks rooms occupied during a step

	// Assign ants to paths in a round-robin way
	for i := 1; i <= ants; i++ {
		antPaths[i] = (i - 1) % len(paths)
		antPositions[i] = -1 // Ant starts outside the paths
	}

	moving := true
	for moving {
		moving = false
		step := []string{} // Stores the movements for this step
		occupiedRooms = make(map[string]bool)

		// Move each ant if possible
		for ant := 1; ant <= ants; ant++ {
			pathIdx := antPaths[ant]
			position := antPositions[ant]

			// Calculate the next position for the ant
			nextPosition := position + 1
			if nextPosition < len(paths[pathIdx]) {
				nextRoom := paths[pathIdx][nextPosition]

				// Move the ant if the room is unoccupied
				if !occupiedRooms[nextRoom] {
					antPositions[ant] = nextPosition
					occupiedRooms[nextRoom] = true

					// Add the move to the step, but skip the start room
					if nextPosition > 0 {
						step = append(step, fmt.Sprintf("L%d-%s", ant, nextRoom))
					}

					moving = true // Movement still occurring
				}
			}
		}

		// Add the step's movements to the result
		if len(step) > 0 {
			result = append(result, strings.Join(step, " "))
		}
	}

	return result
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
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <filename>")
	}

	farm := &Farm{}
	if err := farm.ReadFile(os.Args[1]); err != nil {
		log.Fatal(err)
	}

	adjMatrix, roomNames, startIndex, endIndex := CreateAdjacencyMatrix(farm)
	paths := FindAllPaths(adjMatrix, startIndex, endIndex, roomNames)

	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	antMoves := DistributeAnts(paths, farm.Ants)

	for _, move := range antMoves {
		fmt.Println(move)
	}
}
