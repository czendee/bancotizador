package model

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// ErrorType is the list of allowed values for the error's type.
type ErrorType string

// ErrorTypeParam is the list of allowed values for the parameter error's type.
type ParamErrorType string

const (
	InvalidRequestErrorType ErrorType = "invalid_request"
	NotFoundErrorType       ErrorType = "not_found"
	BadGatewayErrorType     ErrorType = "bad_gateway"
	UnauthorizedErrorType   ErrorType = "unauthorized"

	InvalidValueErrorType ParamErrorType = "invalid_value"
	IsRequiredErrorType   ParamErrorType = "is_required"
)

// ErrorGroup groups the errors appeneded in a list.
type ErrorGroup struct {
	Errors []ErrorStruct

	Status int `json:"-"`
}

// NewErrorGroup initializes the error group object.
func NewErrorGroup() *ErrorGroup {
	return &ErrorGroup{Status: 400}
}

// HTTPStatus returns the HTTP status code
func (e *ErrorGroup) HTTPStatus() int {
	return e.Status
}

// MarshalJSON return the JSON message
func (e *ErrorGroup) MarshalJSON() ([]byte, error) {
	return []byte(e.Error()), nil
}

// Error serializes the error objects list to JSON and returns it as a string.
func (e *ErrorGroup) Error() string {
	ret, _ := json.Marshal(e.Errors)
	return string(ret)
}

// HasError checks if ErrorGroup has a error appended.
func (g *ErrorGroup) HasError() bool {
	return len(g.Errors) > 0
}

// HasError append a error in ErrorGroup.
func (g *ErrorGroup) Push(e *ErrorStruct) {
	g.Errors = append(g.Errors, *e)
}

// ErrorStruct is part of the response when a call is unsucessful.
type ErrorStruct struct {
	Type ErrorType `json:"type"`

	Code    string `json:"code,omitempty"`
	Message string `json:"message"`

	ParamType ParamErrorType `json:"param_type,omitempty"`
	Param     string         `json:"param,omitempty"`

	Status int `json:"-"`
}

// HTTPStatus returns the HTTP status code
func (e *ErrorStruct) HTTPStatus() int {
	return e.Status
}

// Error serializes the error object to JSON and returns it as a string.
func (e *ErrorStruct) Error() string {
	ret, _ := json.Marshal(e)
	return string(ret)
}

// InvalidRequestError is a error that occurs when the request contains invalid
// parameters or the body request is invalid.
func InvalidRequestError(message string) *ErrorStruct {
	return &ErrorStruct{Type: InvalidRequestErrorType, Message: message, Status: 400}
}

// UnauthorizedError is a error that occurs when the gateway returns unauthorized.
func UnauthorizedError(message string) *ErrorStruct {
	return &ErrorStruct{Type: UnauthorizedErrorType, Message: message, Status: 401}
}

// BadGatewayError is a error that occurs when the gateway returns bad response.
func BadGatewayError(message string) *ErrorStruct {
	return &ErrorStruct{Type: BadGatewayErrorType, Message: message, Status: 502}
}

// NotFoundError is a error that occurs when the gateway returns not found.
func NotFoundError(message string) *ErrorStruct {
	return &ErrorStruct{Type: NotFoundErrorType, Message: message, Status: 404}
}

// setParamError initialize the ErrorStruct when the error has ParamType and Param values.
func setParamError(param string, paramType ParamErrorType, message []string) *ErrorStruct {
	var _message string
	if m := message; len(m) > 0 && len(m[0]) > 0 {
		_message = m[0]
	}

	err := InvalidRequestError(_message)

	err.Param = param
	err.ParamType = paramType

	return err
}

// InvalidValueError is a error that occurs when the request contains paramater value invalid.
func InvalidValueError(param string, message ...string) *ErrorStruct {
	if len(message) == 0 {
		message = []string{fmt.Sprintf("The parameter '%s' has invalid value.", param)}
	}
	return setParamError(param, InvalidValueErrorType, message)
}

// IsRequiredError is a error that occurs when the request contains paramater value empty or null.
func IsRequiredError(param string, message ...string) *ErrorStruct {
	if len(message) == 0 {
		message = []string{fmt.Sprintf("The parameter '%s' is required.", param)}
	}
	return setParamError(param, IsRequiredErrorType, message)
}

// UnmarshalJSONError is a error that occurs when the request contains a unmarshal JSON error.
func UnmarshalJSONError(e error) error {
	var message string
	if io.EOF == e {
		message = "json: cannot unmarshal, the request is empty"
	} else if _e, ok := e.(*json.UnmarshalTypeError); ok {
		var _t string
		t := _e.Type
		switch t.Kind() {
		case reflect.Struct:
			_t = "object"
			break
		default:
			_t = t.String()
			break
		}
		message = "json: cannot unmarshal " + _e.Value + " into field " + _e.Field + " of type " + _t
	} else {
		switch e.(type) {
		case *ErrorStruct, *ErrorGroup:
			return e
		default:
			message = "json: cannot unmarshal, " + e.Error()
			break
		}
	}

	return InvalidRequestError(message)
}
