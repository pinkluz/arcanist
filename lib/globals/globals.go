package globals

var trace bool

func SetTrace(b bool) {
	trace = b
}

func GetTrace() bool {
	return trace
}
