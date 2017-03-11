package jsonrpc

// http://www.jsonrpc.org/specification

import (
	"encoding/json"
	"net/http"
)

// Version is JSON-RPC version
const Version = "2.0"

// MethodHandler is
type MethodHandler func(req *Request) (interface{}, error)

// Server is
type Server struct {
	MethodMap map[string]MethodHandler
}

// NewServer return a new Jsonrpc server
func NewServer() *Server {
	return &Server{
		MethodMap: map[string]MethodHandler{},
	}
}

// Request represents a JSON-RPC request sent by a client.
type Request struct {
	// A String specifying the version of the JSON-RPC protocol.
	// MUST be exactly "2.0".
	Version string `json:"jsonrpc"`

	// A String containing the name of the method to be invoked
	Method string `json:"method"`

	// A Structured value that holds the parameter values to be used during the
	// invocation of the method. This member MAY be omitted.
	Params *json.RawMessage `json:"params"`

	// An identifier established by the Client that MUST contain a String, Number,
	// or NULL value if included.
	// If it is not included it is assumed to be a notification.
	// The value SHOULD normally not be Null
	// and Numbers SHOULD NOT contain fractional parts
	ID *json.RawMessage `json:"id"`
}

// rpcResponse represents a JSON-RPC response returned by the server.
type rpcResponse struct {
	// JSON-RPC protocol.
	Version string `json:"jsonrpc"`

	// The Object that was returned by the invoked method. This must be null
	// in case there was an error invoking the method.
	// As per spec the member will be omitted if there was an error.
	Result interface{} `json:"result,omitempty"`

	// An Error object if there was an error invoking the method. It must be
	// null if there was no error.
	// As per spec the member will be omitted if there was no error.
	Error *Error `json:"error,omitempty"`

	// This must be the same id as the request it is responding to.
	ID *json.RawMessage `json:"id"`
}

// Handler is ...
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body and check if RPC method is valid.

	defer r.Body.Close()

	req := new(Request)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeErrorExt(w, ErrorCodeParse, err.Error(), req)
		return
	}
	if req.Version != Version {
		writeErrorExt(w, ErrorCodeInvalidReq, "jsonrpc must be "+Version, req)
		return
	}

	// find method handler
	fn, ok := s.MethodMap[req.Method]
	if !ok {
		writeErrorExt(w, ErrorCodeNoMethod, "Invalid method: "+req.Method, req)
		return
	}

	reply, err := fn(req)
	if err != nil {
		writeError(w, req, err)
		return
	}

	writeResponse(w, req, reply)

}

func writeResponse(w http.ResponseWriter, req *Request, reply interface{}) {
	res := &rpcResponse{
		Version: Version,
		Result:  reply,
		ID:      req.ID,
	}
	doResponse(w, res)
}

func writeErrorExt(w http.ResponseWriter, code ErrorCode, msg string, req *Request) {
	err := &Error{
		Code:    code,
		Message: msg,
		Data:    req,
	}
	writeError(w, req, err)
}

func writeError(w http.ResponseWriter, req *Request, err error) {
	jsonErr, ok := err.(*Error)
	if !ok {
		jsonErr = &Error{
			Code:    ErrorCodeServer,
			Message: err.Error(),
		}
	}
	res := &rpcResponse{
		Version: Version,
		Error:   jsonErr,
		ID:      req.ID,
	}
	doResponse(w, res)
}

func doResponse(w http.ResponseWriter, res *rpcResponse) {
	// ID is null for notifications and they don't have a response.
	if res.ID != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		encoder := json.NewEncoder(w)
		err := encoder.Encode(res)

		// Not sure in which case will this happen. But seems harmless.
		if err != nil {
			writeErrorExt(w, 400, err.Error(), nil)
		}
	}
}
