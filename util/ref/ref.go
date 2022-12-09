package ref

func T[T any](any T) *T {
	return &any
}
