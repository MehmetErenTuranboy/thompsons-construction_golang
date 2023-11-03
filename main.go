package main

import (
	"fmt"
	"regex_to_dfa/tools"
)

func main() {
	input := "abcd"
	fmt.Println("Before addConcatOperators:", input)
	input = tools.AddConcatOperators(input)
	postfixVal := tools.InfixToPostfix(input)

	fmt.Println("Postfix: ", postfixVal) // Changed from 'input' to 'postfixVal'

	automataRes := tools.Compile(postfixVal)
	visited := make(map[*tools.State]bool)

	tools.PrintTransition(automataRes.InitialState, visited)
}
