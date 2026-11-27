package math

func GoLeanArchLog(x float64) (available bool, output float64, failure any) {
	defer func() { failure = recover() }()
	available = haveArchLog
	output = archLog(x)
	return
}
