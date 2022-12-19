package http_err

type ErrorIs interface {
	Is(error) bool
}

type ErrorUnwrap interface {
	Unwrap() error
}
