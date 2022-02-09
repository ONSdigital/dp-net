package responder

// er is the packages error type
type er struct {
	err     error
	message string
	logData map[string]interface{}
}

// Error satisfies the standard library Go error interface
func (e *er) Error() string {
	if e.err == nil {
		return "nil"
	}
	return e.err.Error()
}

// Unwrap implements the standard library Go unwrapper interface
func (e *er) Unwrap() error {
	return e.err
}

// Message satisfies the messager interface which is used to specify
// a response to be sent to the caller in place of the error text for a
// given error. This is useful when you don't want sensitive information
// or implementation details being exposed to the caller which could be
// used to find exploits in our API
func (e *er) Message() string {
	return e.message
}

// LogData satisfies the dataLogger interface which is used to recover
// log data from an error
func (e *er) LogData() map[string]interface{} {
	return e.logData
}
