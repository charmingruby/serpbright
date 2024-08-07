package helper

func Ternary[T any](condition bool, right T, left T) T {
	if condition {
		return right
	}

	return left
}

func EmptyString(s string, rs string) string {
	if s == "" {
		return rs
	}

	return s
}
