// compiler/main.go
package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	// Example Go source code
	src := `
		package main
		import "fmt"
		func hello() {
			fmt.Println("Hello, World!")
		}
		func greet(name string) {
			fmt.Printf("Hello, %s!\n", name)
		}
	`

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	// Print the AST
	//var buf bytes.Buffer
	ast.Print(fset, node)
}
