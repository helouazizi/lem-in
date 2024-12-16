// main.go
package main

import (
	"fmt"
	"os"

	"lem-in/helpers" // add somthig to test in my life
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
	err := farm.ReadFile(fileName)
	if err != nil {
		// fmt.Println("ERROR: invalid data format")
		fmt.Println(err)
		return
	}

	//fmt.Println(farm.StartRoom)
	//fmt.Println(farm.EndRoom)
	//graph := farm.Crate_Graph()
	Matrix, roomNames := farm.Crate_Matricx()
	startIndex := -1
	for i, name := range roomNames {
		if name == farm.StartRoom {
			startIndex = i
			break
		}
	}
	fmt.Println(Matrix)
	paths := farm.Path_Finder(Matrix, startIndex, roomNames)
	fmt.Println(paths)
	fmt.Println("good data")
}
