import GoLean.Strings.ShouldUseASCIISet

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:963–971. The threshold is a signed
    comparison without arithmetic or width-dependent overflow. -/
def shouldUseASCIISet (bufLen : Int) : Bool :=
  GoLean.Strings.shouldUseASCIISet bufLen

theorem shouldUseASCIISet_iff (bufLen : Int) :
    shouldUseASCIISet bufLen = true ↔ 9 ≤ bufLen :=
  GoLean.Strings.shouldUseASCIISet_iff bufLen

end GoLean.Bytes
