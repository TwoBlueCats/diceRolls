package diceRolls

import "testing"

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
			_, _ = Parser("d6^4/(56*3d4)+789+123")
		}
	})

	b.Run("more-complex", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = Parser("(d6^4/(56*3d4)+789+123)*(d6^4/(56*3d4)+789+123)/((d6^4/(56*3d4)+789+123)/(d6^4/(56*3d4)+789+123))")
		}
	})
}
