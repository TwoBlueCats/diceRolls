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

func (r simpleValue) Description(bool) string {
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

func (o operation) Description(detailed bool) string {
	switch o.op {
	case '+':
		switch {
		case o.left == nil:
			return o.right.Description(detailed)
		case o.right == nil:
			return o.left.Description(detailed)
		}
		return "(" + o.left.Description(detailed) + " + " + o.right.Description(detailed) + ")"
	case '-':
		switch {
		case o.left == nil:
			return "-" + o.right.Description(detailed)
		case o.right == nil:
			return o.left.Description(detailed)
		}
		return "(" + o.left.Description(detailed) + " - " + o.right.Description(detailed) + ")"
	case '*':
		switch {
		case o.left == nil:
			return o.right.Description(detailed)
		case o.right == nil:
			return o.left.Description(detailed)
		}
		return "(" + o.left.Description(detailed) + " * " + o.right.Description(detailed) + ")"
	case '/':
		switch {
		case o.left == nil:
			return "0"
		case o.right == nil:
			return o.left.Description(detailed)
		}
		return "(" + o.left.Description(detailed) + " / " + o.right.Description(detailed) + ")"
	case '^':
		switch {
		case o.left == nil:
			return "0"
		case o.right == nil:
			return o.left.Description(detailed)
		}
		return "(" + o.left.Description(detailed) + " ^ " + o.right.Description(detailed) + ")"

	default:
		panic("unknown operation")
	}
}

func (o operation) Value() int {
	switch o.op {
	case '+':
		switch {
		case o.left == nil:
			return o.right.Value()
		case o.right == nil:
			return o.left.Value()
		}
		return o.left.Value() + o.right.Value()
	case '-':
		switch {
		case o.left == nil:
			return -o.right.Value()
		case o.right == nil:
			return o.left.Value()
		}
		return o.left.Value() - o.right.Value()
	case '*':
		switch {
		case o.left == nil:
			return o.right.Value()
		case o.right == nil:
			return o.left.Value()
		}
		return o.left.Value() * o.right.Value()
	case '/':
		switch {
		case o.left == nil:
			return 0
		case o.right == nil:
			return o.left.Value()
		}
		return o.left.Value() / o.right.Value()
	case '^':
		switch {
		case o.left == nil:
			return 0
		case o.right == nil:
			return o.left.Value()
		}
		sum := 0
		right := o.right.Value()
		for i := 0; i < right; i++ {
			sum += o.left.Value()
		}
		return sum

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

func updateTree(tree stackS[Result], tokenPtr *string) (stackS[Result], error) {
	if tokenPtr == nil {
		return tree, nil
	}
	token := *tokenPtr
	*tokenPtr = ""
	if len(token) == 0 {
		return tree, nil
	}
	switch {
	case strings.Contains(token, "d"):
		dice, err := RollDiceNotation(token)
		if err != nil {
			return tree, err
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
			return tree, err
		}
		tree.add(simpleValue{value: value})
	}
	return tree, nil
}

func Parser(expression string) (Result, error) {
	top := '$'

	ops := stackS[rune]{}
	ops.add(top)

	tree := stackS[Result]{}

	token := ""
	var err error
	for _, char := range expression {
		if unicode.IsSpace(char) {
			continue
		}
		if unicode.IsDigit(char) || unicode.IsLetter(char) {
			token += string(char)
			if strings.Count(token, "d") > 1 {
				return nil, errors.New("too much 'd'")
			}
			if unicode.IsLetter(char) && char != 'd' {
				return nil, errors.New("no letters are allowed")
			}
		} else if char == '(' {
			_ = len(token) == 0
			ops.add(char)
		} else if char == ')' {
			for ops.get() != top && ops.get() != '(' {
				tree, err = updateTree(tree, &token)
				if err != nil {
					return nil, err
				}
				token = string(ops.pop())
			}
			tree, err = updateTree(tree, &token)
			if err != nil {
				return nil, err
			}
			ops.pop()
		} else {
			if prior(char) > prior(ops.get()) {
				tree, err = updateTree(tree, &token)
				if err != nil {
					return nil, err
				}
				ops.add(char)
			} else {
				for ops.get() != top && prior(char) <= prior(ops.get()) {
					tree, err = updateTree(tree, &token)
					if err != nil {
						return nil, err
					}
					token = string(ops.pop())
				}
				tree, err = updateTree(tree, &token)
				if err != nil {
					return nil, err
				}
				ops.add(char)
			}
		}
	}
	for ops.get() != top {
		tree, err = updateTree(tree, &token)
		if err != nil {
			return nil, err
		}
		token = string(ops.pop())
	}
	tree, err = updateTree(tree, &token)
	if err != nil {
		return nil, err
	}

	_ = ops.size() == 1

	if tree.size() != 1 {
		return nil, errors.New("some operands left")
	}

	return tree.pop(), nil
}
