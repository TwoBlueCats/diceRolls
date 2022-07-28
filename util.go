package diceRolls

type stackS[T any] struct {
	data []T
}

func (s stackS[T]) get() T {
	if s.size() == 0 {
		var def T
		return def
	}
	return s.data[len(s.data)-1]
}
func (s *stackS[T]) pop() T {
	if s.size() == 0 {
		var def T
		return def
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}
func (s *stackS[T]) add(val T) T {
	s.data = append(s.data, val)
	return val
}
func (s stackS[T]) size() int {
	return len(s.data)
}
