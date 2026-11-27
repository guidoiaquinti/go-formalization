import GoLean.Strings.ASCIISetContains

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:943, 958–961. A nonnil receiver gives
    the selected entry of a 256-element array. None marks a nil panic. -/
def asciiSetContains (receiver : Option (BitVec 8 → Bool)) (b : BitVec 8) : Option Bool :=
  GoLean.Strings.asciiSetContains receiver b

theorem asciiSetContains_entry (table : BitVec 8 → Bool) (b : BitVec 8) :
    asciiSetContains (some table) b = some (table b) :=
  GoLean.Strings.asciiSetContains_entry table b

theorem asciiSetContains_nil (b : BitVec 8) : asciiSetContains none b = none :=
  GoLean.Strings.asciiSetContains_nil b

end GoLean.Bytes
