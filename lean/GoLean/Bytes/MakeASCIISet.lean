import GoLean.Strings.MakeASCIISet

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:947–956. The return values are a
    256-entry table and a flag. A non-ASCII byte stops insertion and
    returns the partially populated table by value. -/
def makeASCIISet (input : List (BitVec 8)) : (BitVec 8 → Bool) × Bool :=
  GoLean.Strings.makeASCIISet input

theorem makeASCIISet_empty : makeASCIISet [] = (fun _ => false, true) :=
  GoLean.Strings.makeASCIISet_empty

theorem makeASCIISet_stop (front suffix : List (BitVec 8)) (b : BitVec 8)
    (ascii : ∀ c ∈ front, c.toNat < 128) (high : 128 ≤ b.toNat) :
    makeASCIISet (front ++ b :: suffix) = ((makeASCIISet front).1, false) :=
  GoLean.Strings.makeASCIISet_stop front suffix b ascii high

theorem makeASCIISet_high (input : List (BitVec 8)) (b : BitVec 8)
    (high : 128 ≤ b.toNat) : (makeASCIISet input).1 b = false :=
  GoLean.Strings.makeASCIISet_high input b high

end GoLean.Bytes
