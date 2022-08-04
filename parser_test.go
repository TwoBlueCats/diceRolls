package diceRolls

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserCalc(t *testing.T) {
	testcases := []struct {
		name       string
		expression string
		result     int
	}{
		{"const", "1", 1},
		{"simple", "1 + 2312", 2313},
		{"simple-2", "2 + 2 * 2", 6},
		{"simple-brackets", "(2 + 2) * 2", 8},
		{"simple-brackets-2", "(2*2+3)*2", 14},
		{"simple-brackets-3", "24*( 2*2+ 3) / 2", 84},
		{"dice", "1d10", 5},
		{"dice-sum", "1d10+1000", 1005},
		{"dice-sum-2", "1d10+1d6", 8},
		{"complex", "24*(2*2-3+1d6)^2/2^2 + 1d10", 55},
		{"negative", "-2", -2},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			rand.Seed(395)
			result, err := Parser(testcase.expression)
			require.NoError(t, err)
			assert.Equal(t, testcase.result, result.Value())
		})
	}

	t.Run("description", func(t *testing.T) {
		rand.Seed(395)
		result, err := Parser("24*(2*2-3)/2^2 + 1d10")
		require.NoError(t, err)
		assert.Equal(t, "(((24 * ((2 * 2) - 3)) / (2 ^ 2)) + 1d10)", result.Description())

		result, err = Parser("-2")
		require.NoError(t, err)
		assert.Equal(t, "-2", result.Description())
	})

	t.Run("dice-roll", func(t *testing.T) {
		rand.Seed(395)
		result, err := Parser("d100^3")
		require.NoError(t, err)
		value1 := result.Value()

		rand.Seed(395)
		result, err = Parser("3d100")
		require.NoError(t, err)
		value2 := result.Value()

		assert.Equal(t, value1, value2)
	})
}

func TestParserFail(t *testing.T) {
	testcases := []struct {
		name       string
		expression string
		text       string
	}{
		{"dice", "1dd10", "too much 'd'"},
		{"letters", "1d10 + 2s6", "no letters are allowed"},
		{"omit-1", "1d", "dice notation must contain number of sides"},
		{"dots", "1.1d6", "strconv.Atoi: parsing \".\": invalid syntax"},
		{"dots-2", "1.1", "strconv.Atoi: parsing \".\": invalid syntax"},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			rand.Seed(395)
			_, err := Parser(test.expression)
			assert.ErrorContains(t, err, test.text)

		})
	}
}

func BenchmarkParser(b *testing.B) {
	b.Run("const", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = Parser("1")
		}
	})

	b.Run("simple", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = Parser("1 + 2312")
		}
	})

	b.Run("complex", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = Parser("d6*4/(56*3d4)+789+123")
		}
	})

	b.Run("more-complex", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = Parser("(d6*4/(56*3d4)+789+123)*(d6*4/(56*3d4)+789+123)/((d6*4/(56*3d4)+789+123)/(d6*4/(56*3d4)+789+123))")
		}
	})
}

func BenchmarkResult(b *testing.B) {
	b.Run("const", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		res, _ := Parser("1")
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			res.Value()
		}
	})

	b.Run("simple", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		res, _ := Parser("1 + 2312")
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			res.Value()
		}
	})

	b.Run("complex", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		res, _ := Parser("d6*4/(56*3d4)+789+123")
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			res.Value()
		}
	})

	b.Run("more-complex", func(b *testing.B) {
		b.ReportAllocs()
		b.StopTimer()
		res, _ := Parser("(d6*4/(56*3d4)+789+123)*(d6*4/(56*3d4)+789+123)/((d6*4/(56*3d4)+789+123)/(d6*4/(56*3d4)+789+123))")
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			res.Value()
		}
	})
}
