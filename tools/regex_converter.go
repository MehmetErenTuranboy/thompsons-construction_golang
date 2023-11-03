package tools

import (
	"fmt"
	"strings"

	"github.com/golang-collections/collections/stack"
)

const EPSILON = rune(0)

func ApplyPrecedence(c rune) int {
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

func OperatorLister(c rune) bool {
	return strings.ContainsRune("|*.", c)
}

func InfixToPostfix(infixRegex string) string {
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
		} else if ApplyPrecedence(c) > -1 {
			for stack.Len() > 0 && OperatorLister(stack.Peek().(rune)) && ApplyPrecedence(c) <= ApplyPrecedence(stack.Peek().(rune)) {
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

type State struct {
	label      rune
	firstEdge  *State
	secondEdge *State
	id         int
	isAccept   bool
}

type nfa struct {
	InitialState *State
	endState     *State
}

var StateID = 0

func NewState(label rune, firstEdge, secondEdge *State, isAccept bool) *State {
	StateID++
	return &State{
		label:      label,
		firstEdge:  firstEdge,
		secondEdge: secondEdge,
		id:         StateID,
		isAccept:   isAccept,
	}
}

func AddConcatOperators(infix string) string {
	var b strings.Builder
	for i, r := range infix {
		if i > 0 && (infix[i-1] >= 'a' && infix[i-1] <= 'z' || infix[i-1] == '*' || infix[i-1] == ')') && (r >= 'a' && r <= 'z' || r == '(') {
			b.WriteRune('.')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func Compile(postfix string) *nfa {
	nfaStack := stack.New()

	for _, c := range postfix {
		switch c {
		case '*':
			peekNFA := nfaStack.Pop().(*nfa)
			initial := NewState(EPSILON, peekNFA.InitialState, nil, false)
			end := NewState(EPSILON, nil, nil, true)
			initial.secondEdge = end
			peekNFA.endState.firstEdge = initial
			peekNFA.endState.secondEdge = end
			peekNFA.InitialState.isAccept = false
			nfaStack.Push(&nfa{InitialState: initial, endState: end})

		case '.':
			nfa2 := nfaStack.Pop().(*nfa)
			nfa1 := nfaStack.Pop().(*nfa)
			nfa1.endState.firstEdge = nfa2.InitialState
			nfa1.endState.isAccept = false
			nfaStack.Push(&nfa{InitialState: nfa1.InitialState, endState: nfa2.endState})

		case '|':
			fmt.Printf("| State")
			nfa2 := nfaStack.Pop().(*nfa)
			nfa1 := nfaStack.Pop().(*nfa)
			initial := NewState(EPSILON, nfa1.InitialState, nfa2.InitialState, false)
			end := NewState(EPSILON, nil, nil, true) // end State should be an accept State
			nfa1.endState.firstEdge = end
			nfa1.endState.isAccept = false // NFA1's end State is no longer an accept State
			nfa2.endState.firstEdge = end
			nfa2.endState.isAccept = false // NFA2's end State is no longer an accept State
			nfaStack.Push(&nfa{InitialState: initial, endState: end})
		default:
			// Create two new States: One initial and one end State.
			end := NewState(EPSILON, nil, nil, true) // This is the accept State.
			initial := NewState(c, end, nil, false)  // The initial State transitions to the accept State on character c.

			// Push the new NFA fragment onto the stack.
			nfaStack.Push(&nfa{InitialState: initial, endState: end})
		}
	}

	finalNFA := nfaStack.Pop().(*nfa)
	return finalNFA
}

func PrintTransition(currentState *State, visited map[*State]bool) {
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
			// Handling of the SPLIT State for Kleene star and union.
			fmt.Printf("node %d splits to node %d and node %d\n", currentState.id, currentState.firstEdge.id, currentState.secondEdge.id)
		} else if currentState.firstEdge != nil {
			// Handling of EPSILON transitions.
			fmt.Printf("node %d goes to node %d on EPSILON\n", currentState.id, currentState.firstEdge.id)
		}
	}

	if currentState.isAccept {

		fmt.Printf("node %d is an accept State\n", currentState.id)
	}

	PrintTransition(currentState.firstEdge, visited)
	PrintTransition(currentState.secondEdge, visited)
}
