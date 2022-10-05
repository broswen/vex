package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPagination(t *testing.T) {
	tests := []struct {
		Url  string
		Page Pagination
	}{
		{
			Url: "https://vex.broswen.com/accounts",
			Page: Pagination{
				Limit:  100,
				Offset: 0,
			},
		},
		{
			Url: "https://vex.broswen.com/accounts?offset=1&limit=1",
			Page: Pagination{
				Limit:  1,
				Offset: 1,
			},
		},
		{
			Url: "https://vex.broswen.com/accounts?limit=-10",
			Page: Pagination{
				Limit:  100,
				Offset: 0,
			},
		},
		{
			Url: "https://vex.broswen.com/accounts?limit=101",
			Page: Pagination{
				Limit:  100,
				Offset: 0,
			},
		},
		{
			Url: "https://vex.broswen.com/accounts?limit=something",
			Page: Pagination{
				Limit:  100,
				Offset: 0,
			},
		},
		{
			Url: "https://vex.broswen.com/accounts?offset=something&limit=2",
			Page: Pagination{
				Limit:  2,
				Offset: 0,
			},
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest(http.MethodGet, test.Url, nil)
		assert.NoError(t, err)
		p := pagination(req)
		assert.Equal(t, test.Page, p)
	}
}
