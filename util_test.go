package diceRolls

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack(t *testing.T) {
	s := stackS[int]{}
	assert.Equal(t, 0, s.get())
	assert.Equal(t, 0, s.pop())
	s.add(1)
	s.add(2)
	s.add(3)
	assert.Equal(t, 3, s.pop())
	assert.Equal(t, 2, s.size())

}
