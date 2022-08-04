package diceRolls

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Result interface {
	Value() int
	Description(bool) string
}

type rollResult struct {
	rolls int
	sides int

	details []int
}

var _ Result = (*rollResult)(nil)

func (r *rollResult) Description(detailed bool) string {
	if detailed {
		return fmt.Sprintf("%dd%d %v", r.rolls, r.sides, r.details)
	}
	return fmt.Sprintf("%dd%d", r.rolls, r.sides)
}

func (r *rollResult) Value() int {
	sum := 0
	r.details = r.details[:0]
	for i := 0; i < r.rolls; i++ {
		val := rand.Intn(r.sides)
		sum += val + 1
		r.details = append(r.details, val)
	}
	return sum
}

// RollSimpleDice rolls k-side dice n times and return sum of result
// If times is negative returns error
// If sides is negative returns error
func RollSimpleDice(n, k int) (Result, error) {
	if n < 0 {
		return nil, errors.New("number of rolls must be non-negative")
	}
	if k <= 0 {
		return nil, errors.New("number of sides must be positive")
	}

	r := rollResult{rolls: n, sides: k}

	return &r, nil
}

func RollDiceNotation(notation string) (Result, error) {
	var err error
	n := 1
	pos := strings.Index(notation, "d")

	if pos == -1 {
		return nil, errors.New("dice notation must contain 'd' symbol")
	}
	if pos == len(notation)-1 {
		return nil, errors.New("dice notation must contain number of sides")
	}

	if pos != 0 {
		n, err = strconv.Atoi(notation[:pos])
		if err != nil {
			return nil, err
		}
	}
	k, err := strconv.Atoi(notation[pos+1:])
	if err != nil {
		return nil, err
	}
	return RollSimpleDice(n, k)
}
