package main

import "fmt"

// TODO: parse types and try to remove hardcodedreplacements

func main() {
	fmt.Println("GoTGProto Files Generator (c)2023")
	fmt.Println("Running Generator...")
	fmt.Println("Generating generic helpers for context.go")
	generateCUHelpers()
	fmt.Println("Generated Successfully!")
}
