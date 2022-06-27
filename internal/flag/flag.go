package flag

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

type FlagType string

const (
	STRING  FlagType = "STRING"
	BOOLEAN FlagType = "BOOLEAN"
	NUMBER  FlagType = "NUMBER"
)

type Flag struct {
	ID        string   `json:"id"`
	ProjectID string   `json:"-" db:"project_id"`
	Key       string   `json:"key" db:"flag_key"`
	Type      FlagType `json:"type" db:"flag_type"`
	Value     string   `json:"value" db:"flag_value"`
}

func (f Flag) ToJSON() ([]byte, error) {
	return json.Marshal(&f)
}

func Validate(f Flag) error {
	if f.ProjectID == "" {
		return errors.New("project id must not be empty")
	}

	if f.Key == "" {
		return errors.New("flag key must not be empty")
	}

	switch f.Type {
	case BOOLEAN:
		if f.Value != "false" && f.Value != "true" {
			return errors.New("invalid value for boolean flag")
		}
	case NUMBER:
		_, err := strconv.ParseFloat(f.Value, 64)
		if err != nil {
			return errors.New("invalid value for number flag")
		}
	case STRING:
	default:
		return errors.New("invalid flag type")
	}
	return nil
}

func RenderConfig(flags []*Flag) ([]byte, error) {
	config := make(map[string]struct {
		Value string
		Type  FlagType
	}, len(flags))
	for _, f := range flags {
		config[f.Key] = struct {
			Value string
			Type  FlagType
		}{Value: f.Value, Type: f.Type}
	}
	b := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(b).Encode(config)
	return b.Bytes(), err
}
