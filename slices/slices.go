package slices

func Detect[A any](input []A, f func(A) bool) bool {
	count := len(input)
	for i := 0; i < count; i++ {
		if f(input[i]) {
			return true
		}
	}
	return false
}

func Select[A any](input []A, f func(A) bool) []A {
	count := len(input)
	res := make([]A, 0, count)
	for i := 0; i < count; i++ {
		if f(input[i]) {
			res = append(res, input[i])
		}
	}
	return res
}
