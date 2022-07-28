package diceRolls

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type simpleValue struct {
	value int
}

var _ Result = (*simpleValue)(nil)

func (r simpleValue) Description() string {
	return fmt.Sprintf("%d", r.value)
}

func (r simpleValue) Value() int {
	return r.value
}

type operation struct {
	op    rune
	left  Result
	right Result
}

var _ Result = (*operation)(nil)

func (o operation) Description() string {
	switch o.op {
	case '+':
		return "(" + o.left.Description() + " + " + o.right.Description() + ")"
	case '-':
		return "(" + o.left.Description() + " - " + o.right.Description() + ")"
	case '*':
		return "(" + o.left.Description() + " * " + o.right.Description() + ")"
	case '/':
		return "(" + o.left.Description() + " / " + o.right.Description() + ")"

	default:
		panic("unknown operation")
	}
}

func (o operation) Value() int {
	switch o.op {
	case '+':
		return o.left.Value() + o.right.Value()
	case '-':
		return o.left.Value() - o.right.Value()
	case '*':
		return o.left.Value() * o.right.Value()
	case '/':
		return o.left.Value() / o.right.Value()

	default:
		panic("unknown operation")
	}
}

func prior(ch rune) int {
	if ch == '+' || ch == '-' {
		return 1 //Precedence of + or - is 1
	} else if ch == '*' || ch == '/' {
		return 2 //Precedence of * or / is 2
	} else if ch == '^' {
		return 3 //Precedence of ^ is 3
	} else {
		return 0
	}
}

func Parser(expression string) (Result, error) {
	tokens, err := makePostfix(expression, true)
	if err != nil {
		return nil, err
	}
	tree := stackS[Result]{}
	for _, token := range tokens {
		switch {
		case strings.Contains(token, "d"):
			dice, err := RollDiceNotation(token)
			if err != nil {
				return nil, err
			}
			tree.add(dice)
		case prior(rune(token[0])) > 0:
			right := tree.pop()
			var left Result
			if tree.size() > 0 {
				left = tree.pop()
			}
			tree.add(operation{
				op:    rune(token[0]),
				left:  left,
				right: right,
			})
		default:
			value, err := strconv.Atoi(token)
			if err != nil {
				return nil, err
			}
			tree.add(simpleValue{value: value})
		}
	}

	return tree.pop(), err
}

func makePostfix(expression string, letter bool) ([]string, error) {
	top := '$'
	tokens := make([]string, 0)
	stack := stackS[rune]{}
	stack.add(top)
	cur := ""

	for _, char := range expression {
		if unicode.IsSpace(char) {
			continue
		}
		if unicode.IsDigit(char) || unicode.IsLetter(char) {
			cur += string(char)
			if strings.Count(cur, "d") > 1 {
				return nil, errors.New("too much 'd'")
			}
			if !letter && unicode.IsLetter(char) && char != 'd' {
				return nil, errors.New("no letters are allowed")
			}
		} else if char == '(' {
			if len(cur) > 0 {
				tokens = append(tokens, cur)
				cur = ""
			}
			stack.add(char)
		} else if char == '^' {
			if len(cur) > 0 {
				tokens = append(tokens, cur)
				cur = ""
			}
			stack.add(char)
		} else if char == ')' {
			for stack.get() != top && stack.get() != '(' {
				if len(cur) > 0 {
					tokens = append(tokens, cur)
					cur = ""
				}
				cur = string(stack.pop())
			}
			if len(cur) > 0 {
				tokens = append(tokens, cur)
				cur = ""
			}
			stack.pop()
		} else {
			if prior(char) > prior(stack.get()) {
				if len(cur) > 0 {
					tokens = append(tokens, cur)
					cur = ""
				}
				stack.add(char)
			} else {
				for stack.get() != top && prior(char) <= prior(stack.get()) {
					if len(cur) > 0 {
						tokens = append(tokens, cur)
						cur = ""
					}
					cur = string(stack.pop())
				}
				if len(cur) > 0 {
					tokens = append(tokens, cur)
					cur = ""
				}
				stack.add(char)
			}
		}
	}
	for stack.get() != top {
		if len(cur) > 0 {
			tokens = append(tokens, cur)
			cur = ""
		}
		cur = string(stack.pop())
	}
	if len(cur) > 0 {
		tokens = append(tokens, cur)
		cur = ""
	}
	return tokens, nil
}
