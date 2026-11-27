package math

func GoLeanArchExp2(x float64) (available bool, output float64, failure any) {
	defer func() { failure = recover() }()
	available = haveArchExp2
	output = archExp2(x)
	return
}
