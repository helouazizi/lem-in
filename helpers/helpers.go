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
	X, Y int
}

// lets reprasent our farm of ants as struct with it's properties
type Farm struct {
	Rooms              []string
	Links              []string
	StartRoom, EndRoom string
	Ants               int
}

/*
this function will read the file data
and checking the foramt of the data
if the data is in the correct format
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

	return nil
}

func (F *Farm) CheckDoubles(data []string) bool {
	index := 0
	for i := 1; i < len(data)-1; i++ {
		for j := i + 1; j < len(data); j++ {
			if data[i] == data[j] {
				return false
			}
		}
	}
	for i := 1; i < len(data); i++ {
		check := strings.Split(data[i], " ")
		if len(check) != 1 && len(check) != 3 {
			return false
		}
		if i < len(data)-1 && data[i] == "##start" {
			F.StartRoom = data[i+1]
			check := strings.Split(F.StartRoom, " ")
			if len(check) != 3 || strings.Contains(check[0], "L") || strings.Contains(check[0], "#") {
				return false
			}
			F.Rooms = append(F.Rooms, F.StartRoom)
			index++
		}
		if i < len(data)-1 && data[i] == "##end" {
			F.EndRoom = data[i+1]
			check := strings.Split(F.EndRoom, " ")
			if len(check) != 3 || strings.Contains(check[0], "L") || strings.Contains(check[0], "#") {
				return false
			}
			F.Rooms = append(F.Rooms, F.EndRoom)
			index++
		}

		if len(check) != 3 && data[i] != "##start" && data[i] != "##end" {
			F.Links = append(F.Links, data[i])
		} else {
			F.Rooms = append(F.Rooms, data[i])
		}

	}
	if index != 2 {
		return index == 2
	}

	return true
}
