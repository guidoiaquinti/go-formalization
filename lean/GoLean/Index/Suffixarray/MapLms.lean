import Init.Data.BitVec.Basic

namespace GoLean.Index.Suffixarray

private def packLms (length : Nat) (backing : Array (BitVec width)) :
    Option (Array (BitVec width)) :=
  if length = 0 then none
  else
    let indices := (List.range (length / 2 + 1)).reverse
    let (_, output) := indices.foldl (fun (write, state) i =>
      let id := state[i]?.getD 0
      if id.toInt > 0 then
        (write - 1, state.set! (write - 1) (id - 1))
      else (write, state)) (length, backing)
    some output

/-- Go 1.27.0 src/index/suffixarray/sais.go:638–646. Scan the low half
    through its midpoint, compacting positive IDs into the high end. -/
def mapLms32 (length : Nat) (backing : Array (BitVec 32)) (_numLMS : Int) :
    Option (Array (BitVec 32)) :=
  packLms length backing

/-- Go 1.27.0 src/index/suffixarray/sais2.go:1176–1184. The count parameter is unused. -/
def mapLms64 (length : Nat) (backing : Array (BitVec 64)) (_numLMS : Int) :
    Option (Array (BitVec 64)) :=
  packLms length backing

theorem mapLms32_compacts :
    mapLms32 6 #[4, 0, 5, 8, 19, 21, 99] (-8) =
      some #[4, 0, 5, 3, 4, 7, 99] := rfl

theorem mapLms64_ignoresCount (length : Nat) (backing : Array (BitVec 64)) (x y : Int) :
    mapLms64 length backing x = mapLms64 length backing y := rfl

theorem mapLms64_compacts :
    mapLms64 6 #[4, 0, 5, 8, 19, 21, 99] 3 =
      some #[4, 0, 5, 3, 4, 7, 99] := rfl

theorem mapLms32_empty (backing : Array (BitVec 32)) (count : Int) :
    mapLms32 0 backing count = none := rfl

end GoLean.Index.Suffixarray
