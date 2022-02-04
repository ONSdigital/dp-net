package main

type testError struct{
	err        error
	statusCode int
	logData    map[string]interface{}
	message    string
}

// standard Go error interfaces
func (e *testError) Error(){
	if e.err == nil{
		return "nil"
	}
	return e.err.Error()
}

func (e *testError) Unwrap(){
	return e.err
}

// Optional interfaces to implement to use full functionality
// of responder
func (e *testError) Code(){
	return e.statusCode
}

func (e *testError) LogData(){
	return e.logData
}

func (e *testError) Message(){
	return e.message
}
