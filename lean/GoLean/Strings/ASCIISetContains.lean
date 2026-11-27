namespace GoLean.Strings

/-- Go 1.27.0 strings/strings.go:964–966. A nonnil pointer provides 256
    Boolean entries. None represents the nil-pointer panic. -/
def asciiSetContains (receiver : Option (BitVec 8 → Bool)) (b : BitVec 8) : Option Bool :=
  receiver.map (fun table => table b)

theorem asciiSetContains_entry (table : BitVec 8 → Bool) (b : BitVec 8) :
    asciiSetContains (some table) b = some (table b) := rfl

theorem asciiSetContains_nil (b : BitVec 8) : asciiSetContains none b = none := rfl

end GoLean.Strings
