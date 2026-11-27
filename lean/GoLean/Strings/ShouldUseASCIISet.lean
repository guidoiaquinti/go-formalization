namespace GoLean.Strings

/-- Go 1.27.0 src/strings/strings.go:968–976. The lookup-table threshold
    depends on the buffer length alone. No arithmetic can overflow a Go int. -/
def shouldUseASCIISet (bufLen : Int) : Bool :=
  decide (bufLen > 8)

theorem shouldUseASCIISet_iff (bufLen : Int) :
    shouldUseASCIISet bufLen = true ↔ 9 ≤ bufLen := by
  simp [shouldUseASCIISet]
  omega

end GoLean.Strings
