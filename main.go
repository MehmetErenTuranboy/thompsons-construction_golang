package main

import (
	"fmt"
	"strings"
)

func operatorLister(c rune) bool {
	return strings.ContainsRune("|*.", c)
}

func infixToPostfix(infixRegex string) string {
	var resultInPostfix strings.Builder
	var stack []rune
	// var cursor rune
	// var cacheStack []rune

	fmt.Printf("Input of infixToPostfix: %s\n", infixRegex)

	for _, c := range infixRegex {
		if c >= 'a' && c <= 'z' {
			resultInPostfix.WriteRune(c)
		} else if c == '(' {
			stack = append(stack, c)
		} else if c == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				resultInPostfix.WriteRune(stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
		} else if operatorLister(c) {
			for len(stack) > 0 && operatorLister(stack[len(stack)-1]) {
				resultInPostfix.WriteRune(stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, c)
		}
	}

	for len(stack) > 0 {
		resultInPostfix.WriteRune(stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	fmt.Printf("Output of infixToPostfix: %s\n", resultInPostfix.String())
	return resultInPostfix.String()

}

func main() {
	inputInfixValue := "aa.b|b|b"
	infixToPostfix(inputInfixValue)

}
