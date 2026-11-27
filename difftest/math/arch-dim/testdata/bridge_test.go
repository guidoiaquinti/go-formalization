package math

func GoLeanArchMax(x, y float64) (available bool, output float64, failure any) {
	defer func() { failure = recover() }()
	available = haveArchMax
	output = archMax(x, y)
	return
}

func GoLeanArchMin(x, y float64) (available bool, output float64, failure any) {
	defer func() { failure = recover() }()
	available = haveArchMin
	output = archMin(x, y)
	return
}
