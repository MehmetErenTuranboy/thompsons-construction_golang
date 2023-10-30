package main

import (
	"fmt"
	"strings"

	"github.com/golang-collections/collections/stack"
)

// applyPresedence returns the precedence of a given operator
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

// operatorLister checks if a character is an operator
func operatorLister(c rune) bool {
	return strings.ContainsRune("|*.", c)
}

// infixToPostfix converts an infix regular expression to postfix notation
func infixToPostfix(infixRegex string) string {
	var resultInPostfix strings.Builder
	stack := stack.New()

	for _, c := range infixRegex {
		if c >= 'a' && c <= 'z' {
			resultInPostfix.WriteRune(c)
		} else if c == '(' {
			stack.Push(c)
		} else if c == ')' {
			for stack.Len() > 0 && stack.Peek() != '(' {
				resultInPostfix.WriteRune(stack.Pop().(rune))
			}
			stack.Pop() // Pop '('
		} else if operatorLister(c) {
			for stack.Len() > 0 && operatorLister(stack.Peek().(rune)) && applyPresedence(c) <= applyPresedence(stack.Peek().(rune)) {
				resultInPostfix.WriteRune(stack.Pop().(rune))
			}
			stack.Push(c)
		}
	}

	for stack.Len() > 0 {
		resultInPostfix.WriteRune(stack.Pop().(rune))
	}

	fmt.Printf("Output of infixToPostfix: %s\n", resultInPostfix.String())
	return resultInPostfix.String()
}

type state struct {
	label      rune
	firstEdge  *state
	secondEdge *state
}

type nfa struct {
	initialState *state
	endState     *state
}

func addConcatOperators(infix string) string {
	var b strings.Builder
	for i, r := range infix {
		if i > 0 && (infix[i-1] >= 'a' && infix[i-1] <= 'z' || infix[i-1] == ')') && (r >= 'a' && r <= 'z' || r == '(') {
			b.WriteRune('.')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// compile creates an NFA from a postfix regular expression
func compile(postfix string) *nfa {
	nfaStack := stack.New()

	for _, c := range postfix {
		switch c {
		case '*':
			childNFA := nfaStack.Pop().(*nfa)
			initial := &state{label: 0, firstEdge: childNFA.initialState, secondEdge: nil}
			end := &state{label: 0, firstEdge: nil, secondEdge: nil}
			initial.secondEdge = end
			childNFA.endState.firstEdge = childNFA.initialState
			childNFA.endState.secondEdge = end
			nfaStack.Push(&nfa{initialState: initial, endState: end})

		case '.':
			nfa2 := nfaStack.Pop().(*nfa)
			nfa1 := nfaStack.Pop().(*nfa)
			nfa1.endState.firstEdge = nfa2.initialState
			nfaStack.Push(&nfa{initialState: nfa1.initialState, endState: nfa2.endState})

		default:
			initial := &state{label: c, firstEdge: nil, secondEdge: nil}
			nfaStack.Push(&nfa{initialState: initial, endState: initial})
		}
	}

	result := nfaStack.Pop().(*nfa)
	fmt.Printf("Result initial state label: %c\n", result.initialState.label)
	return result
}

// printStates recursively prints the states of the NFA
func printStates(currentState *state, visited map[*state]bool) {
	if currentState == nil || visited[currentState] {
		return
	}

	visited[currentState] = true

	// Using state labels to indicate special types of states
	stateLabel := string(currentState.label)
	if currentState.label == 0 {
		if currentState.firstEdge != nil && currentState.secondEdge != nil {
			stateLabel = "SPLIT"
		} else {
			stateLabel = "EPSILON"
		}
	}

	fmt.Printf("State: %p, Label: %s\n", currentState, stateLabel)

	printStates(currentState.firstEdge, visited)
	printStates(currentState.secondEdge, visited)
}

func main() {
	input := "aaa*"
	input = addConcatOperators(input)
	postfixVal := infixToPostfix(input)

	fmt.Println("Postfix: ", postfixVal)

	automataRes := compile(postfixVal)
	visited := make(map[*state]bool)
	printStates(automataRes.initialState, visited)
}
