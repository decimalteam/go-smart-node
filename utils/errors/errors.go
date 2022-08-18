package errors

import (
	"encoding/json"
)

type Err struct {
	Description string            `json:"description"`
	Codespace   string            `json:"codespace"`
	Code        uint32            `json:"code"`
	Params      map[string]string `json:"params,omitempty"`
}

func (e Err) Error() string {
	bz, _ := json.Marshal(e)
	return string(bz)
}

type Param struct {
	Key   string
	Value string
}

func NewParam(key, value string) Param {
	return Param{
		Key:   key,
		Value: value,
	}
}

func Encode(codespace string, errorcode uint32, description string, params ...Param) error {
	err := Err{
		Description: description,
		Codespace:   codespace,
		Code:        errorcode,
	}

	if params != nil {
		err.Params = make(map[string]string)
		for _, param := range params {
			err.Params[param.Key] = param.Value
		}
	}

	return err
}
