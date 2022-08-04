# Dice rolling

Simple dice rolling library.

## Usage

The main entrypoint to the library is the diceRolls.Parse function:

```go
func Parser(expression string) (Result, error)
```

`Result` has to methods:

- `Result.Value() int` - evaluate parsed expression and roll dices. Dice are rolled on each call.
- `Result.Description(detailed bool) string` - return string representation of expression. If `detailed == true` it will include
  dice results in square brackets.
