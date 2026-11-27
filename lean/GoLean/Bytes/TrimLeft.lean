import GoLean.Strings.Trim

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:1006–1060. The shared forward scanner
    removes leading cutset members. Slicing drops the same number of bytes
    from length and capacity, and an empty result is nil. -/
def trimLeft (input cutset : List (BitVec 8)) (spare : Nat) :
    Option (List (BitVec 8) × Nat) :=
  let bytes := GoLean.Strings.trimLeft input cutset
  if bytes.isEmpty then none else some (bytes, bytes.length + spare)

theorem trimLeft_empty (cutset : List (BitVec 8)) (spare : Nat) :
    trimLeft [] cutset spare = none := by
  simp [trimLeft, GoLean.Strings.trimLeft]

theorem trimLeft_emptyCutset (head : BitVec 8) (tail : List (BitVec 8)) (spare : Nat) :
    trimLeft (head :: tail) [] spare = some (head :: tail, (head :: tail).length + spare) := by
  simp [trimLeft, GoLean.Strings.trimLeft_emptyCutset]

theorem trimLeft_singleASCII (input : List (BitVec 8)) (byte : BitVec 8) (spare : Nat)
    (ascii : byte.toNat < 128) :
    trimLeft input [byte] spare =
      let suffix := input.dropWhile (· == byte)
      if suffix.isEmpty then none else some (suffix, suffix.length + spare) := by
  simp [trimLeft, GoLean.Strings.trimLeft_singleASCII input byte ascii]

end GoLean.Bytes
