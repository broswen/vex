package flag

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRenderConfig(t *testing.T) {
	tests := []struct {
		flags []*Flag
		json  []byte
	}{
		{
			flags: []*Flag{
				{
					ID:         "1",
					ProjectID:  "2",
					AccountID:  "3",
					CreatedOn:  time.Now(),
					ModifiedOn: time.Now(),
					Key:        "feature1",
					Type:       "STRING",
					Value:      "test",
				},
			},
			json: []byte("{\"feature1\":{\"value\":\"test\",\"type\":\"STRING\"}}\n"),
		},
	}
	for _, tc := range tests {
		j, err := RenderConfig(tc.flags)
		assert.Nil(t, err)
		assert.Equalf(t, tc.json, j, "expected %s but got %s", tc.json, j)
	}
}
