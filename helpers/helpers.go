// helpers/helpers.go
package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// lets represent our room as struct with it 's properties
type Room struct {
	// Name string
	X, Y string
}

// lets reprasent our farm of ants as struct with it's properties
type Farm struct {
	Rooms              map[string]*Room
	Links              map[string][]string
	StartRoom, EndRoom string
	Ants               int
	FileSize           int64
}

// / lets represent our graph

/*
this function will read the file data
and checking the foramt of the data
if the data is in the correct formatF.StartRoom = data[i+1]
by checking the number of ants and the
and rooms representation is correct
and valid links between rooms or any doublacate rooms and links
or any invalid data it will return an error
this function check data validation
1. check the first line is number for number of ants
2. check the ##start and ##end room is exist and not duplecated
3. check the rooms are not duplicated and never start with a 'L' or '#' and must have valid and unique  cordonates x,y
4. check the links between rooms are valid and check the room is exist or not because we cant link into a non exist room
*/
func (F *Farm) ReadFile(fileName string) error {
	// open the file
	var err error
	// graph := Add_Graph()
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

func (F *Farm) Path_Finder() [][]string {
	jiran_dyal_start := F.Links[F.StartRoom]
	F.Rooms = nil

	result := [][]string{}

	for _, jar := range jiran_dyal_start {
		visited := make(map[string]bool)

		visited[F.StartRoom] = true
		visited[jar] = true
		queue := [][]string{{F.StartRoom, jar}}

		for len(queue) > 0 {

			path := queue[0]
			queue = queue[1:]

			node := path[len(path)-1]

			if node == F.EndRoom {

				result = append(result, path)
			
			}

			for _, neighbor := range F.Links[node] {

				if !contains(path, neighbor) && !visited[neighbor] {
					visited[neighbor]= true

					newPath := append([]string{}, path...)

					newPath = append(newPath, neighbor)
					queue = append(queue, newPath)
					
				}
			}
		}

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]

		// Explore neighbors
		for i := 0; i < len(adjacencyMatrix[vertex]); i++ {
			if adjacencyMatrix[vertex][i] && !visited[i] {
				visited[i] = true
				queue = append(queue, i)

				// Create a new path for this neighbor
				newPath := append([]string(nil), currentPath...) // Copy current path
				newPath = append(newPath, roomNames[i])          // Add the neighbor
				paths = append(paths, newPath)                   // Store the new path
			}
		}
	}
	
	return result
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// func (F *Farm) Unique_Path(result [][]string, paths [][]string) []string {
// 	res := []string{}
// 	for _, newpath := range paths{
// 		for _, oldpath := range result{
// 			if !F.contains(oldpath, newpath){
// 				res = newpath
// 				break

// 			}

// 		}

// 	}
// 	return res

// }
// func (F *Farm)contains(oldpath, newpath []string) bool {
// 	for _, oldroom := range oldpath {
// 		for _, newroom := range newpath {
// 			if oldroom == newroom && oldroom != F.EndRoom &&  oldroom!= F.StartRoom{
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }
