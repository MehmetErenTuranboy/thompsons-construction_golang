package main

import (
	"fmt"
	"strings"

	"github.com/golang-collections/collections/stack"
)

// presidence list of symbols
func applyPresedence(c rune) int {
	switch c {
	case '*':
		return 3
	case '.':
		return 2
	case '|':
		return 1
	default:
		return -1
	}
}

func operatorLister(c rune) bool {
	return strings.ContainsRune("|*.", c)
}

// converting infix to posfix for the sake of making regex to nfa construction easier
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

// nfa node construction
type state struct {
	label      rune
	firstEdge  *state
	secondEdge *state
}

// nfa's state construction
type nfa struct {
	initialState *state
	endState     *state
}

func compile(postfix string) *nfa {
	nfaStack := stack.New()

	for _, c := range postfix {
		switch c {
		case '*':
			childNFA := nfaStack.Pop().(*nfa)
			initial := &state{label: 0, firstEdge: nil, secondEdge: nil}
			end := &state{label: 0, firstEdge: nil, secondEdge: nil}
			initial.firstEdge = childNFA.initialState
			initial.secondEdge = end
			childNFA.endState.firstEdge = childNFA.initialState
			childNFA.endState.secondEdge = end
			nfaStack.Push(&nfa{initialState: initial, endState: end})
		case '.':
			nfa2 := nfaStack.Pop().(*nfa)
			nfa1 := nfaStack.Pop().(*nfa)

			// Connect nfa1's end state to nfa2's initial state
			nfa1.endState.firstEdge = nfa2.initialState

			// Update the end state of the resulting NFA
			resultNFA := &nfa{
				initialState: nfa1.initialState,
				endState:     nfa2.endState,
			}
			nfaStack.Push(resultNFA)
		default:
			initial := &state{label: c, firstEdge: nil, secondEdge: nil}
			end := &state{label: 0, firstEdge: nil, secondEdge: nil}
			initial.firstEdge = end
			nfaStack.Push(&nfa{initialState: initial, endState: end})
		}
	}

	result := nfaStack.Pop().(*nfa)

	return result
}

// for testing purposes
func printStates(currentState *state, visited map[*state]bool) {
	if currentState == nil || visited[currentState] {
		return
	}

	visited[currentState] = true

	fmt.Printf("State: %p, Label: %c\n", currentState, currentState.label)

	printStates(currentState.firstEdge, visited)
	printStates(currentState.secondEdge, visited)
}

func main() {
	inputInfixValue := "aa.b|b|b"
	infixToPostfix(inputInfixValue)

	postfix := "caacabc*."
	resultNFA := compile(postfix)

	// You can now use the resultNFA for further processing or visualization.
	fmt.Printf("Initial State Label: %c\n", resultNFA.initialState.label)
	fmt.Printf("End State Label: %c\n", resultNFA.endState.label)

	fmt.Println("All States in the NFA:")
	visited := make(map[*state]bool)
	printStates(resultNFA.initialState, visited)

}
