package consts

// ----- KEY-VALUE RESULTS -----
const (
	// KeyNameError - error key name
	KeyNameError = "error"
	// KeyNameResult - result key name
	KeyNameResult = "result"
	// ValResultSuccess - success result value
	ValResultSuccess = "success"
	// ValResultFailure - failure result value
	ValResultFailure = "failure"
)

// MakeSuccessResult - make rest service result for success
func MakeSuccessResult() map[string]string {
	return map[string]string{KeyNameResult: ValResultSuccess}
}

// MakeFailureResult - make rest service result for failure
func MakeFailureResult(errMsg string) map[string]string {
	return map[string]string{
		KeyNameResult: ValResultFailure,
		KeyNameError:  errMsg,
	}
}
