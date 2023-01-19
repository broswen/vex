package flag

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"
)

type Type string

const (
	STRING  Type = "STRING"
	BOOLEAN Type = "BOOLEAN"
	NUMBER  Type = "NUMBER"
)

type Flag struct {
	ID         string    `json:"id"`
	ProjectID  string    `json:"project_id" db:"project_id"`
	AccountID  string    `json:"account_id" db:"account_id"`
	CreatedOn  time.Time `json:"created_on" db:"created_on"`
	ModifiedOn time.Time `json:"modified_on" db:"modified_on"`
	Key        string    `json:"key" db:"flag_key"`
	Type       Type      `json:"type" db:"flag_type"`
	Value      string    `json:"value" db:"flag_value"`
}

func (f Flag) ToJSON() ([]byte, error) {
	return json.Marshal(&f)
}

func Validate(f Flag) error {
	if f.ProjectID == "" {
		return ErrInvalidData{"project id must not be empty"}
	}

	if f.Key == "" {
		return ErrInvalidData{"flag key must not be empty"}
	}

	switch f.Type {
	case BOOLEAN:
		if f.Value != "false" && f.Value != "true" {
			return ErrInvalidData{"invalid value for boolean flag"}
		}
	case NUMBER:
		_, err := strconv.ParseFloat(f.Value, 64)
		if err != nil {
			return ErrInvalidData{"invalid value for number flag"}
		}
	case STRING:
	case "":
		return ErrInvalidData{"flag type must not be empty"}
	default:
		return ErrInvalidData{"invalid flag type"}
	}
	return nil
}

func (f *Flag) Render() ([]byte, error) {
	jFlag := JsonFlag{
		Key:   f.Key,
		Value: f.Value,
		Type:  f.Type,
	}
	b := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(b).Encode(jFlag)
	return b.Bytes(), err
}

type JsonFlag struct {
	AccountId string `json:"account_id"`
	ProjectId string `json:"project_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Type      Type   `json:"type"`
}

func RenderConfig(flags []*Flag) ([]byte, error) {
	config := make(map[string]JsonFlag)
	for _, f := range flags {
		config[f.Key] = JsonFlag{
			Key:   f.Key,
			Value: f.Value,
			Type:  f.Type,
		}
	}
	b := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(b).Encode(config)
	return b.Bytes(), err
}
