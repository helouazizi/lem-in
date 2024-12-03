package main

import (
	"fmt"
	"os"

	"lem-in/helpers"
)

func main() {
	// lets check the argumments to protect the main function
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . exemple.txt")
		return
	}
	// now lets read the file to extract ou data
	fileName := os.Args[1]
	data, err := helpers.ReadFile(fileName)
	if err != nil {
		fmt.Println("ERROR: invalid data format 1")
		return
	}
	numOfants, sratrRoom, endRoom, err := helpers.ParseData(data)
	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}
	fmt.Println(numOfants, "\n", sratrRoom, "\n", endRoom)
}
