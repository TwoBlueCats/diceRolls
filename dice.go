package dice_rolls

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type SimpleResult struct {
	rolls int
	sides int
	sum   int
}

var _ Result = (*SimpleResult)(nil)

func (r SimpleResult) Description() string {
	return fmt.Sprintf("%d (%dd%d)", r.sum, r.rolls, r.sides)
}

func (r SimpleResult) Value() int {
	return r.sum
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

	r := SimpleResult{rolls: n, sides: k, sum: 0}
	for i := 0; i < n; i++ {
		r.sum += rand.Intn(k) + 1
	}
	return r, nil
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