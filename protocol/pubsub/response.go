package pubsub

// Response to RPC request
type Response struct {
	Version string      `json:"tomatorpc"`
	ID      string      `json:"id,omitempty"`
	Type    string      `json:"message_type"`
	Entity  string      `json:"entity"`
	Method  string      `json:"method"`
	Status  string      `json:"status,omitempty"`
	Err     string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponseTo generates a success response to
// a request
func SuccessResponseTo(req Request, data interface{}) Response {
	return Response{
		Version: "0.2",
		ID:      req.ID,
		Type:    "response",
		Entity:  req.Entity,
		Method:  req.Method,
		Status:  "success",
		Data:    data,
	}
}

// ErrorResponseTo generates an error response to
// a request
func ErrorResponseTo(req Request, err error) Response {
	return Response{
		Version: "0.2",
		ID:      req.ID,
		Type:    "response",
		Entity:  req.Entity,
		Method:  req.Method,
		Status:  "error",
		Err:     err.Error(),
	}
}
