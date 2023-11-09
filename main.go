package main

import (
	"fmt"

	"github.com/MehmetErenTuranboy/thompsons-construction_golang/tools"
)

func main() {
	input := "WOR"
	fmt.Println("Before addConcatOperators:", input)
	input = tools.AddConcatOperators(input)
	postfixVal := tools.InfixToPostfix(input)

	fmt.Println("Postfix: ", postfixVal) // Changed from 'input' to 'postfixVal'

	automataRes := tools.Compile(postfixVal)
	visited := make(map[*tools.State]bool)

	tools.PrintTransition(automataRes.InitialState, visited)
}
