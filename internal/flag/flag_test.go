package flag

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		flag Flag
		err  error
	}{
		{
			flag: Flag{
				ProjectID: "1",
				Key:       "test",
				Type:      "STRING",
				Value:     "test",
			},
			err: nil,
		},
		{
			flag: Flag{
				ProjectID: "1",
				Key:       "",
				Type:      "STRING",
				Value:     "test",
			},
			err: ErrInvalidData{"flag key must not be empty"},
		},
		{
			flag: Flag{
				ProjectID: "1",
				Key:       "test",
				Type:      "",
				Value:     "test",
			},
			err: ErrInvalidData{"flag type must not be empty"},
		},
		{
			flag: Flag{
				ProjectID: "1",
				Key:       "test",
				Type:      "UNKNOWN",
				Value:     "test",
			},
			err: ErrInvalidData{"invalid flag type"},
		},
		{
			flag: Flag{
				ProjectID: "1",
				Key:       "test",
				Type:      "NUMBER",
				Value:     "a",
			},
			err: ErrInvalidData{"invalid value for number flag"},
		},
		{
			flag: Flag{
				ProjectID: "1",
				Key:       "test",
				Type:      "BOOLEAN",
				Value:     "a",
			},
			err: ErrInvalidData{"invalid value for boolean flag"},
		},
	}

	for _, test := range tests {
		err := Validate(test.flag)
		assert.ErrorIs(t, err, test.err)
	}
}

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
				{
					ID:         "2",
					ProjectID:  "2",
					AccountID:  "3",
					CreatedOn:  time.Now(),
					ModifiedOn: time.Now(),
					Key:        "feature2",
					Type:       "BOOLEAN",
					Value:      "true",
				},
				{
					ID:         "3",
					ProjectID:  "2",
					AccountID:  "3",
					CreatedOn:  time.Now(),
					ModifiedOn: time.Now(),
					Key:        "feature3",
					Type:       "NUMBER",
					Value:      "123",
				},
			},
			json: []byte("{\"feature1\":{\"value\":\"test\",\"type\":\"STRING\"},\"feature2\":{\"value\":\"true\",\"type\":\"BOOLEAN\"},\"feature3\":{\"value\":\"123\",\"type\":\"NUMBER\"}}\n"),
		},
	}
	for _, tc := range tests {
		j, err := RenderConfig(tc.flags)
		assert.Nil(t, err)
		assert.Equalf(t, tc.json, j, "expected %s but got %s", tc.json, j)
	}
}
