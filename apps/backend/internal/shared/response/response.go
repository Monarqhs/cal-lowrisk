// Package response provides the standard API response envelope used by all
// controllers, so both clients (Flutter, Next.js) parse one shape.
// See docs/02-system/architecture.md §5.5.
package response

import "github.com/gin-gonic/gin"

// Envelope is the single response shape for the whole API.
type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

// ErrorBody is the error payload inside the envelope.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// OK writes a success envelope with the given status and data.
func OK(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Envelope{Success: true, Data: data})
}

// Fail writes an error envelope with the given status, machine code, and message.
func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, Envelope{Success: false, Error: &ErrorBody{Code: code, Message: message}})
}
