package api

import (
	"encoding/json"
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

// Error contains Decimal API error response fields.
type Error struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Err        string `json:"error"`
}

// Error returns error info as string.
func (e Error) Error() string {
	return fmt.Sprintf("statusCode: %d, message: \"%s\", data: \"%s\"", e.StatusCode, e.Message, e.Err)
}

// ResponseError wraps Resty response error and allows to generate error info.
type ResponseError struct {
	*resty.Response
}

// NewResponseError creates new ResponseError object.
func NewResponseError(response *resty.Response) *ResponseError {
	return &ResponseError{Response: response}
}

// Error returns error info as JSON string.
func (res ResponseError) Error() string {
	detailError := map[string]string{
		"statusCode": fmt.Sprintf("%d", res.StatusCode()),
		"status":     res.Status(),
		"time":       fmt.Sprintf("%f seconds", res.Time().Seconds()),
		"receivedAt": fmt.Sprintf("%v", res.ReceivedAt()),
		"headers":    fmt.Sprintf("%v", res.Header()),
		"body":       res.String(),
	}
	marshal, _ := json.Marshal(detailError)
	return string(marshal)
}

func processConnectionError(response *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if response.IsError() {
		return NewResponseError(response)
	}
	return nil
}

// closure function type for universalJSONDecode
type validationFuncType func() (bool, bool)

// Error indicating for universal decoding
var ErrIsRPCError = errors.New("rpc error")
var ErrMissingLogic = errors.New("universal JSON decode missing logic")

// Universal JSON decoding
// valueStruct - MUST BE pointer to struct
// errorStruct - MUST BE pointer to struct or nil
// return error if valueStruct unmarshaling failed or valueStruct not ok
func universalJSONDecode(source []byte, valueStruct interface{}, errorStruct interface{}, validator validationFuncType) error {
	var err1, err2 error
	err1 = json.Unmarshal(source, valueStruct)
	if errorStruct != nil {
		err2 = json.Unmarshal(source, errorStruct)
	}
	okValue, okError := validator()
	// all ok
	if okValue && err1 == nil {
		return nil
	}
	// error ok
	if okError {
		return ErrIsRPCError // indicate that error in errorStruct
	}

	// error during Unmarshaling (wrong JSON)
	if !okError && err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}

	return fmt.Errorf("%w: input value=%#v, input error=%#v, okValue=%#v, okError=%#v,  err1=%v, err2=%v",
		ErrMissingLogic, valueStruct, errorStruct, okValue, okError, err1, err2)
}

func joinErrors(err1, err2 error) error {
	if errors.Is(err1, ErrIsRPCError) {
		return err2
	}
	return err1
}
