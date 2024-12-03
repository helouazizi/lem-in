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
	// line by line using the  function
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || (len(line) > 1 && line[0] == '#' && line[1] != '#') {
			// here we skip the empty lines and the comments
			continue
		} else {
			splitedLine := strings.Split(line, " ")
			if len(splitedLine) == 1 || len(splitedLine) == 3 {
				data = append(data, line)
			}else{
				return nil, errors.New("invalid data format")
			}

		}

	}
	// data, err := os.ReadFile(fileName)
	// if err != nil {
	// 	return nil, err
	// }

	return data, nil
}

func ParseData(data []string) (int, string, string, error) {
	var sratrRoom, endRoom string
	// var startindex , endindex int
	numOfants, err := strconv.Atoi(data[0])
	if err != nil {
		return 0, "", "", err
	}
	for i, target := range data {
		if target == "##start" {
			sratrRoom = data[i+1]
		}
		if target == "##end" {
			endRoom = data[i+1]
		}
	}
	if sratrRoom == "" || endRoom == "" {
		return 0, "", "", errors.New("err")
	}
	return numOfants, sratrRoom, endRoom, nil
}
