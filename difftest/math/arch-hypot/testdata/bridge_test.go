package math

func GoLeanArchHypot(x, y float64) (available bool, output float64, failure any) {
	defer func() { failure = recover() }()
	available = haveArchHypot
	output = archHypot(x, y)
	return
}
