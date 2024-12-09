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
	Name string
	X, Y string
}

// lets reprasent our farm of ants as struct with it's properties
type Farm struct {
	Rooms              map[string]Room
	Links              map[string][]string
	StartRoom, EndRoom string
	Ants               int
	AntsPositions      map[int]string // Maps ant number to current room name
	AntMoves           [][]string
}

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
		F.Rooms = make(map[string]Room)
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
				F.Rooms[check[0]] = Room{Name: check[0], X: check[1], Y: check[2]}
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
	queue := [][]string{{F.StartRoom}}
	result := [][]string{}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		currentRoom := path[len(path)-1]
		// lets append the path if the room is the end room
		if currentRoom == F.EndRoom {
			// lets check if the pats have seme room at the index because this can ce a collesion
			if notcollesion(result, path) {
				result = append(result, path)
			}
			// result = append(result, path)
		}
		// lets get all the rooms that are connected to the current room
		for _, connection := range F.Links[currentRoom] {
			if !contains(path, connection) {
				newPath := append([]string{}, path...)
				newPath = append(newPath, connection)
				queue = append(queue, newPath)
			}
		}

	}
	return result
}

func contains(path []string, connection string) bool {
	for _, connected := range path {
		if connected == connection {
			return true
		}
	}
	return false
}

func notcollesion(result [][]string, path []string) bool {
	for _, oldpath := range result {
		minlen := len(oldpath)
		// this condition avoiding out of range
		if len(path) < len(oldpath) {
			minlen = len(path)
		}
		// be sure to ignore strat room and end room
		for i := 1; i < minlen-1; i++ {
			if oldpath[i] == path[i] {
				return false
			}
		}

	}
	return true
}
