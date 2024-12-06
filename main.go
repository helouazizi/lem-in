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
	farm := &helpers.Farm{}

	// now lets read the file to extract ou data
	fileName := os.Args[1]
	data, err := helpers.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = farm.ValidateData(data)
	if err != nil {
		fmt.Println("ERROR: invalid data format 2")
		return
	}
	fmt.Println(farm.Ants)
	fmt.Println(farm.StartRoom)
	fmt.Println(farm.EndRoom)
	fmt.Println(farm.Rooms)
	fmt.Println(farm.Links)
	fmt.Println("good data")
}
