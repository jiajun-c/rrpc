package message

type Header struct {
	Method   string   // the method name
	Version  uint64   // the version number
	SeqNum   uint64   // the sequence number chosen by the client used to mark the call
	Error    string   // the error message
	AuthType AuthType // the AuthType
}

type AuthType uint8

const (
	NullAuth  AuthType = 0 // null authentication
	UnixAuth  AuthType = 1 // Unix authentication
	JwtAuth   AuthType = 2 // JWT authentication
	BasicAuth AuthType = 3 // Basic authentication
)
