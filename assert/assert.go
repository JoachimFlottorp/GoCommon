package assert

import "go.uber.org/zap"

// Error: is a helper function to assert that an error is nil
//
// Panics if the error is not nil
//
// ALlows setting a custom error message
func Error(e error, message ...string) {
	if e != nil {
		if len(message) > 0 {
			panic(message[0] + ": " + e.Error())
		}
		panic(e.Error())
	}	
}

// IsTrue: Validates that the statement is true
func IsTrue(statement bool, message ...string) {
	if !statement {
		if len(message) > 0 {
			zap.S().Warn(message)
		}
		panic("assertion failed")
	}
}