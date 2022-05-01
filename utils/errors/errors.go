package errors

import (
	"encoding/json"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Err struct {
	Description string            `json:"description"`
	Codespace   string            `json:"codespace"`
	Params      map[string]string `json:"params,omitempty"`
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

func Encode(codespace string, errorcode uint32, description string, params ...Param) *sdkerrors.Error {
	err := Err{
		Description: description,
		Codespace:   codespace,
	}

	if params != nil {
		err.Params = make(map[string]string)
		for _, param := range params {
			err.Params[param.Key] = param.Value
		}
	}

	result, _ := json.Marshal(err)

	return sdkerrors.New(
		codespace,
		errorcode,
		string(result),
	)
}
