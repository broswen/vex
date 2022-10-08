package api

import (
	"errors"
	"github.com/broswen/vex/internal/account"
	"github.com/broswen/vex/internal/flag"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRenderError(t *testing.T) {
	tests := []struct {
		Err             error
		ExpectedErr     error
		ExpectedMessage string
	}{
		{
			Err:             errors.New("unknown"),
			ExpectedErr:     ErrUnknown,
			ExpectedMessage: "unknown error",
		},
		{
			Err:             flag.ErrInvalidData{"invalid value for number flag"},
			ExpectedErr:     ErrBadRequest.WithError(flag.ErrInvalidData{"invalid value for number flag"}),
			ExpectedMessage: "bad request: invalid value for number flag",
		},
		{
			Err:             account.ErrAccountNotFound{Err: nil},
			ExpectedErr:     ErrNotFound,
			ExpectedMessage: "not found",
		},
	}

	for _, tc := range tests {
		e := translateError(tc.Err)
		assert.Equal(t, tc.ExpectedErr, e)
		assert.Equal(t, tc.ExpectedMessage, e.Message)
		assert.Equal(t, tc.ExpectedMessage, e.Error())
	}
}
