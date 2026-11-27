package math

// GoLeanCallArchStub calls the inventory symbol through its original package.
// A panic must hold the exact string value, and a normal return is distinct.
func GoLeanCallArchStub(name string, a, b uint64) (available bool, status string) {
	x := Float64frombits(a)
	status = "returned"
	defer func() {
		if failure := recover(); failure != nil {
			if message, ok := failure.(string); ok {
				status = "panic:" + message
			} else {
				status = "unexpected non-string panic"
			}
		}
	}()
	switch name {
	case "archAcos":
		available = haveArchAcos
		archAcos(x)
	case "archAcosh":
		available = haveArchAcosh
		archAcosh(x)
	case "archAsin":
		available = haveArchAsin
		archAsin(x)
	case "archAsinh":
		available = haveArchAsinh
		archAsinh(x)
	case "archAtan":
		available = haveArchAtan
		archAtan(x)
	case "archAtan2":
		available = haveArchAtan2
		archAtan2(x, Float64frombits(b))
	case "archAtanh":
		available = haveArchAtanh
		archAtanh(x)
	case "archCbrt":
		available = haveArchCbrt
		archCbrt(x)
	case "archCos":
		available = haveArchCos
		archCos(x)
	case "archCosh":
		available = haveArchCosh
		archCosh(x)
	case "archErf":
		available = haveArchErf
		archErf(x)
	case "archErfc":
		available = haveArchErfc
		archErfc(x)
	case "archExpm1":
		available = haveArchExpm1
		archExpm1(x)
	case "archFrexp":
		available = haveArchFrexp
		archFrexp(x)
	case "archLdexp":
		available = haveArchLdexp
		archLdexp(x, int(b))
	case "archLog10":
		available = haveArchLog10
		archLog10(x)
	case "archLog2":
		available = haveArchLog2
		archLog2(x)
	case "archLog1p":
		available = haveArchLog1p
		archLog1p(x)
	case "archMod":
		available = haveArchMod
		archMod(x, Float64frombits(b))
	case "archPow":
		available = haveArchPow
		archPow(x, Float64frombits(b))
	case "archRemainder":
		available = haveArchRemainder
		archRemainder(x, Float64frombits(b))
	case "archSin":
		available = haveArchSin
		archSin(x)
	case "archSinh":
		available = haveArchSinh
		archSinh(x)
	case "archTan":
		available = haveArchTan
		archTan(x)
	case "archTanh":
		available = haveArchTanh
		archTanh(x)
	default:
		return false, "unknown stub"
	}
	return available, status
}
