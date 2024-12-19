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
	StartNeighbots     []string
	badrooms           [][]string
	StartRoom, EndRoom string
	Ants               int
}

var f Farm

func constructor() {
	for i := 0; i < len(f.StartNeighbots); i++ {
		Paths = append(Paths, []string{"", "", "", "", "", "", "", "", "", ""})
	}
}

var Paths = make([][]string, len(f.StartNeighbots))

// CreateAdjacencyMatrix creates an adjacency matrix from the Farm struct
func CreateAdjacencyMatrix(farm *Farm) ([][]bool, []string, int, int) {
	farm.StartNeighbots = farm.Links[farm.StartRoom]
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
		bad := []string{roomName}
		if roomName == farm.StartRoom {
			startIndex = i
		}
		if roomName == farm.EndRoom {
			endIndex = i
		}
		for _, linkedRoom := range farm.Links[roomName] {
			bad = append(bad, linkedRoom)
			for j, linkedRoomName := range roomNames {
				if linkedRoomName == linkedRoom {
					adjacencyMatrix[i][j] = true
					break
				}
			}
		}
		farm.badrooms = append(farm.badrooms, bad)
		delete(farm.Links, roomName)
	}

	return adjacencyMatrix, roomNames, startIndex, endIndex
}

func (F *Farm) FindAllPaths(adjacencyMatrix [][]bool, start int, end int, roomNames []string, ants int, badrooms []string) ([][]string, [][]string) {
	constructor()
	var badpaths [][]string
	var currentPath []string
	f = *F
	// pathstructure := &Paths{
	// 	Paths: make([][]string, len(F.StartNeighbots)),
	// }
	visited := make([]bool, len(roomNames))
	enough := false
	var dfs func(vertex int)
	dfs = func(vertex int) {
		visited[vertex] = true
		currentPath = append(currentPath, roomNames[vertex])

		if vertex == end {
			if len(Paths) == ants {
				enough = true
			}
			if !contains(currentPath, badrooms) {
				index, best := bestone(currentPath)
				if best && index == -1 {
					Paths = append(Paths, currentPath)
				} else if best && index != -1 {
					Paths[index] = currentPath
				}

			} else {
				badpaths = append(badpaths, append([]string(nil), currentPath...))
			}
		} else {
			// Explore neighbors
			for i := 0; i < len(adjacencyMatrix[vertex]); i++ {
				if adjacencyMatrix[vertex][i] && !visited[i] {
					if enough {
						return
					}
					dfs(i)
				}
			}
		}

		// Backtrack
		visited[vertex] = false
		currentPath = currentPath[:len(currentPath)-1]
	}

	dfs(start)
	return Paths, badpaths
}

func bestone(cureentpath []string) (int, bool) {
	for i, path := range Paths {
		if path == nil {
			return -1, true
		} else {
			oldneighbor := path[0]
			neighbor := cureentpath[0]
			if oldneighbor == neighbor {
				if len(cureentpath) < len(path) {
					return i, true
				}
			}
		}
	}
	if len(Paths) == 0 {
		return -1, true
	}
	return -1, false
}

func contains(s []string, e []string) bool {
	for _, a := range s {
		for _, b := range e {
			if a == b {
				return true
			}
		}
	}
	return false
}

func Filter(goodpaths, badpaths [][]string) [][]string {
	result := make([][]string, 0)
	for i, badpath := range badpaths {
		count := 0
		for _, goodpath := range goodpaths {
			if !contains(badpath[1:len(badpath)-1], goodpath[1:len(goodpath)-1]) {
				count++
			} else {
				// delete this path
				badpaths[i] = nil
			}
		}
		if count == len(goodpaths) || len(goodpaths) == 0 {
			result = append(result, badpath)
		}

	}
	return result
}

func (F *Farm) ReadFile(fileName string) error {
	// open the file
	var err error

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
			// _, exist2 := F.Links[link[0]]
			// if exist2 {
			// 	return errors.New("no valid link found")
			// }
			// _, exist3 := F.Links[link[1]]
			// if exist3 {
			// 	return errors.New("no valid link found")
			// }

			F.Links[link[0]] = append(F.Links[link[0]], link[1])
			F.Links[link[1]] = append(F.Links[link[1]], link[0])

			// graph.Add_Edges(link[1],link[0])

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
	// Parse the input datacd
	farm := &Farm{
		Rooms: make(map[string]*Room),
		Links: make(map[string][]string),
	}

	err := farm.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("file read")

	// Create the adjacency matrix
	adjMatrix, roomNames, startIndex, endIndex := CreateAdjacencyMatrix(farm)
	// lets sort the badrooms
	fmt.Println("matrix created")
	sort.Slice(farm.badrooms, func(i, j int) bool {
		return len(farm.badrooms[i]) < len(farm.badrooms[j])
	})
	maxsize := len(farm.badrooms[len(farm.badrooms)-1])
	all_badrooms := []string{}
	for _, barroom := range farm.badrooms {
		if len(barroom) == maxsize && barroom[0] != farm.EndRoom && barroom[0] != farm.StartRoom {
			all_badrooms = append(all_badrooms, barroom[0])
		}
	}
	fmt.Println("badrooms extracted")
	fmt.Println(all_badrooms)
	// baaaaadroom := all_badrooms[len(all_badrooms)-1]
	// for _, room := range all_badrooms {
	// 	fmt.Println(room)
	// }

	// Find all paths
	_, _ = farm.FindAllPaths(adjMatrix, startIndex, endIndex, roomNames, farm.Ants, all_badrooms)
	// extrapaths := filter(paths.Paths, badpaths)
	// goodpaths = append(goodpaths, extrapaths...)

	// Sort paths by length
	// sort.Slice(goodpaths, func(i, j int) bool {
	// 	return len(goodpaths[i]) < len(goodpaths[j])
	// })
	// for _, path := range goodpaths {
	// 	fmt.Println(path, "good")
	// }
	// for _, badpath := range badpaths {
	// 	fmt.Println(badpath, "bad")
	// }
	// fmt.Println(farm.StartNeighbots)
	// fmt.Println(len(goodpaths))
	// fmt.Println(goodpaths)
	// fmt.Println(baaaaadroom)
	for _, path := range Paths {
		fmt.Println(path, "after")
	}
}
