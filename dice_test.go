package diceRolls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleDice(t *testing.T) {
	t.Run("ok-multiple-times", func(t *testing.T) {
		sR, err := RollSimpleDice(10, 1)
		assert.NoError(t, err)
		assert.Equal(t, 10, sR.Value())
	})
	t.Run("ok-zero-times", func(t *testing.T) {
		sR, err := RollSimpleDice(0, 1)
		assert.NoError(t, err)
		assert.Equal(t, 0, sR.Value())
	})
	t.Run("ok-multiple-sides", func(t *testing.T) {
		sR, err := RollSimpleDice(1, 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, sR.Value(), 1)
		assert.Less(t, sR.Value(), 10)
	})

	t.Run("fail-zero-sides", func(t *testing.T) {
		sR, err := RollSimpleDice(10, 0)
		assert.ErrorContains(t, err, "number of sides must be positive")
		assert.Nil(t, sR)
	})
	t.Run("fail-negative-sides", func(t *testing.T) {
		sR, err := RollSimpleDice(10, -10)
		assert.ErrorContains(t, err, "number of sides must be positive")
		assert.Nil(t, sR)
	})
	t.Run("fail-negative-times", func(t *testing.T) {
		sR, err := RollSimpleDice(-10, 10)
		assert.ErrorContains(t, err, "number of rolls must be non-negative")
		assert.Nil(t, sR)
	})
}

func TestDiceNotation(t *testing.T) {
	t.Run("ok-multiple-times", func(t *testing.T) {
		sR, err := RollDiceNotation("10d1")
		assert.NoError(t, err)
		assert.Equal(t, 10, sR.Value())
	})
	t.Run("ok-multiple-times", func(t *testing.T) {
		sR, err := RollDiceNotation("0d1")
		assert.NoError(t, err)
		assert.Equal(t, 0, sR.Value())
	})
	t.Run("ok-multiple-sides", func(t *testing.T) {
		sR, err := RollDiceNotation("1d10")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, sR.Value(), 1)
		assert.Less(t, sR.Value(), 10)
	})
	t.Run("ok-default-times", func(t *testing.T) {
		sR, err := RollDiceNotation("d1")
		assert.NoError(t, err)
		assert.Equal(t, sR.Value(), 1)
	})

	t.Run("fail-zero-sides", func(t *testing.T) {
		sR, err := RollDiceNotation("10d0")
		assert.ErrorContains(t, err, "number of sides must be positive")
		assert.Nil(t, sR)
	})
	t.Run("fail-negative-sides", func(t *testing.T) {
		sR, err := RollDiceNotation("10d-10")
		assert.ErrorContains(t, err, "number of sides must be positive")
		assert.Nil(t, sR)
	})
	t.Run("fail-negative-times", func(t *testing.T) {
		sR, err := RollDiceNotation("-10d10")
		assert.ErrorContains(t, err, "number of rolls must be non-negative")
		assert.Nil(t, sR)
	})
	t.Run("fail-not-dice", func(t *testing.T) {
		sR, err := RollDiceNotation("11000")
		assert.ErrorContains(t, err, "dice notation must contain 'd' symbol")
		assert.Nil(t, sR)
	})
	t.Run("fail-no-sides", func(t *testing.T) {
		sR, err := RollDiceNotation("10d")
		assert.ErrorContains(t, err, "dice notation must contain number of sides")
		assert.Nil(t, sR)
	})
	t.Run("fail-not-times", func(t *testing.T) {
		sR, err := RollDiceNotation("10.1d12")
		assert.Error(t, err)
		assert.Nil(t, sR)
	})
	t.Run("fail-not-sides", func(t *testing.T) {
		sR, err := RollDiceNotation("10d12.123")
		assert.Error(t, err)
		assert.Nil(t, sR)
	})
}
