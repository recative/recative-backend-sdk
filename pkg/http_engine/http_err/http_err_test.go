package http_err

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponseErrorType_Is(t *testing.T) {
	err := InternalServerError.New("test")
	assert.Equal(t, true, InternalServerError.Is(err))
}

func TestResponseError_Is(t *testing.T) {
	err := InternalServerError.New("test")
	assert.Equal(t, true, errors.Is(err, InternalServerError.New()))
}

func TestResponseError_Unwrap(t *testing.T) {
	var err = errors.New("this is a test error")
	var responseError = InternalServerError.Wrap(err)
	assert.Equal(t, true, errors.Is(responseError, err))
}
