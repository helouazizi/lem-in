// helpers/helpers.go
package helpers

import (
	"bufio"
	"errors"
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
*/
func ReadFile(fileName string) ([]string, error) {
	// open the file
	data := []string{}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// read the file by using the buffio pkg
	// that can give us convenient way to read input from a file
	// line by line using the  function scan()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" || (line[0] == '#' && line != "##start" && line != "##end") {
			continue
		}
		data = append(data, line)
	}
	if len(data) < 6 {
		return nil, errors.New("invalid data") // not enough data
	}
	return data, nil
}

/*
this function check data validation
1. check the first line is number for number of ants
2. check the ##start and ##end room is exist and not duplecated
3. check the rooms are not duplicated and never start with a 'L' or '#' and must have valid and unique  cordonates x,y
4. check the links between rooms are valid and check the room is exist or not because we cant link into a non exist room
*/

func (F *Farm) ValidateData(data []string) error {
	var err error

	F.Ants, err = strconv.Atoi(data[0])
	if err != nil {
		return err
	}
	if F.Ants <= 0 {
		return errors.New("invalid number of ants")
	}

	if !F.CheckDoubles(data) {
		return errors.New("duplicates found")
	}

	err = F.RoomsTraitment(data)
	if err != nil {
		return err
	}
	err = F.LinksTraitement(data)
	if err != nil {
		return err
	}

	return nil
}

func (F *Farm) CheckDoubles(data []string) bool {
	index := 0
	for i := 1; i < len(data); i++ {
		check := strings.Split(data[i], " ")
		if len(check) != 1 && len(check) != 3 {
			return false
		}
		for j := i + 1; j < len(data); j++ {
			if data[i] == data[j] {
				return false
			}
		}
		if data[i] == "##start" || data[i] == "##end" {
			index++
		}
	}

	if index != 2 {
		return index == 2
	}

	return true
}

func (F *Farm) LinksTraitement(data []string) error {
	if F.Links == nil {
		F.Links = make(map[string][]string)
	}
	for _, link := range data {
		check := strings.Split(link, "-")
		if len(check) != 2 {
			continue
		}
		_, exist := F.Rooms[check[0]]
		if !exist {
			return errors.New("room not found")
		}
		_, exist1 := F.Rooms[check[1]]
		if !exist1 {
			return errors.New("room not found")
		}
		F.Links[check[0]] = append(F.Links[check[0]], check[1])
		F.Links[check[1]] = append(F.Links[check[1]], check[0])

	}
	return nil
}

func (F *Farm) RoomsTraitment(data []string) error {
	if F.Rooms == nil {
		F.Rooms = make(map[string]Room)
	}
	for i, line := range data {
		if i < len(data)-1 && line == "##start" || line == "##end" {
			check := strings.Split(data[i+1], " ")
			if len(check) != 3 {
				return errors.New("invalid start room")
			}
			if line == "##start" {
				F.StartRoom = check[0]
			} else {
				F.EndRoom = check[0]
			}
			// add the the room to the map
			F.Rooms[check[0]] = Room{Name: check[0], X: check[1], Y: check[2]}
		}
		check := strings.Split(line, " ")
		if len(check) == 3 {
			F.Rooms[check[0]] = Room{Name: check[0], X: check[1], Y: check[2]}
		}
	}
	return nil
}

/*
lets use bfs algorithm to find the shortest path between stert and end rooms
*/

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
			//result = append(result, path)
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
