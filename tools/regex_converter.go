package tools

import (
	"fmt"
	"strings"

	"github.com/golang-collections/collections/stack"
)

const EPSILON = rune(0)

func applyPrecedence(c rune) int {
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
		} else if applyPrecedence(c) > -1 {
			for stack.Len() > 0 && operatorLister(stack.Peek().(rune)) && applyPrecedence(c) <= applyPrecedence(stack.Peek().(rune)) {
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
	id         int
	isAccept   bool
}

type nfa struct {
	initialState *state
	endState     *state
}

var stateID = 0

func newState(label rune, firstEdge, secondEdge *state, isAccept bool) *state {
	stateID++
	return &state{
		label:      label,
		firstEdge:  firstEdge,
		secondEdge: secondEdge,
		id:         stateID,
		isAccept:   isAccept,
	}
}

func addConcatOperators(infix string) string {
	var b strings.Builder
	for i, r := range infix {
		if i > 0 && (infix[i-1] >= 'a' && infix[i-1] <= 'z' || infix[i-1] == '*' || infix[i-1] == ')') && (r >= 'a' && r <= 'z' || r == '(') {
			b.WriteRune('.')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func compile(postfix string) *nfa {
	nfaStack := stack.New()

	for _, c := range postfix {
		switch c {
		case '*':
			peekNFA := nfaStack.Pop().(*nfa)
			initial := newState(EPSILON, peekNFA.initialState, nil, false)
			end := newState(EPSILON, nil, nil, true)
			initial.secondEdge = end
			peekNFA.endState.firstEdge = initial
			peekNFA.endState.secondEdge = end
			peekNFA.initialState.isAccept = false
			nfaStack.Push(&nfa{initialState: initial, endState: end})

		case '.':
			nfa2 := nfaStack.Pop().(*nfa)
			nfa1 := nfaStack.Pop().(*nfa)
			nfa1.endState.firstEdge = nfa2.initialState
			nfa1.endState.isAccept = false
			nfaStack.Push(&nfa{initialState: nfa1.initialState, endState: nfa2.endState})

		case '|':
			fmt.Printf("| state")
			nfa2 := nfaStack.Pop().(*nfa)
			nfa1 := nfaStack.Pop().(*nfa)
			initial := newState(EPSILON, nfa1.initialState, nfa2.initialState, false)
			end := newState(EPSILON, nil, nil, true) // end state should be an accept state
			nfa1.endState.firstEdge = end
			nfa1.endState.isAccept = false // NFA1's end state is no longer an accept state
			nfa2.endState.firstEdge = end
			nfa2.endState.isAccept = false // NFA2's end state is no longer an accept state
			nfaStack.Push(&nfa{initialState: initial, endState: end})
		default:
			// Create two new states: One initial and one end state.
			end := newState(EPSILON, nil, nil, true) // This is the accept state.
			initial := newState(c, end, nil, false)  // The initial state transitions to the accept state on character c.

			// Push the new NFA fragment onto the stack.
			nfaStack.Push(&nfa{initialState: initial, endState: end})
		}
	}

	finalNFA := nfaStack.Pop().(*nfa)
	return finalNFA
}

func printTransition(currentState *state, visited map[*state]bool) {
	if currentState == nil || visited[currentState] {
		return
	}

	visited[currentState] = true

	if currentState.label != EPSILON {
		if currentState.firstEdge != nil {
			fmt.Printf("node %d takes %s goes to node %d\n", currentState.id, string(currentState.label), currentState.firstEdge.id)
		}
	} else {
		if currentState.firstEdge != nil && currentState.secondEdge != nil {
			// Handling of the SPLIT state for Kleene star and union.
			fmt.Printf("node %d splits to node %d and node %d\n", currentState.id, currentState.firstEdge.id, currentState.secondEdge.id)
		} else if currentState.firstEdge != nil {
			// Handling of EPSILON transitions.
			fmt.Printf("node %d goes to node %d on EPSILON\n", currentState.id, currentState.firstEdge.id)
		}
	}

	if currentState.isAccept {

		fmt.Printf("node %d is an accept state\n", currentState.id)
	}

	printTransition(currentState.firstEdge, visited)
	printTransition(currentState.secondEdge, visited)
}
