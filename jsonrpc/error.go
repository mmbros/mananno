package jsonrpc

import (
	"errors"
)

// ErrorCode represents the error type that occurred.
type ErrorCode int

// Error Code constants
const (
	ErrorCodeParse      ErrorCode = -32700
	ErrorCodeInvalidReq ErrorCode = -32600
	ErrorCodeNoMethod   ErrorCode = -32601
	ErrorCodeBadParams  ErrorCode = -32602
	ErrorCodeInternal   ErrorCode = -32603
	ErrorCodeServer     ErrorCode = -32000 /* -32000 to -32099 */
)

// ErrNullResult is ...
var ErrNullResult = errors.New("result is null")

// Error represents an JSON-RPC error object
type Error struct {
	// A Number that indicates the error type that occurred.
	// This MUST be an integer.
	Code ErrorCode `json:"code"` /* required */

	// A String providing a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	Message string `json:"message"` /* required */

	// A Primitive or Structured value that contains additional information about the error.
	Data interface{} `json:"data"` /* optional */
}

func (e *Error) Error() string {
	return e.Message
}
