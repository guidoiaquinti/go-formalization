import GoLean.Bytes.ConfiguredIter
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def render (r : IterRun) : String :=
  let status := if r.boundsPanic then "panic" else "ok"
  let parts := r.fragments.map fun part =>
    let tag := if part.bytes.isNone then "nil" else "slice"
    s!"{tag}:{part.capacity}:{String.intercalate "," ((part.bytes.getD []).map toString)}"
  status ++ "/" ++ String.intercalate ";" parts

private def runs (action : IterState → Nat → IterRun) (state : IterState) (first second : Nat) : String :=
  let a := action state first
  let b := action a.remaining second
  let c := action b.remaining 4096
  let d := action c.remaining 4096
  String.intercalate "|" [render a, render b, render c, render d]

def main : IO Unit := runTuple fun row => do
  let tag :: size :: spare :: sepSize :: save :: first :: second :: raw := row |
    throw (IO.userError "expected configured iterator fields")
  let data ← raw.mapM (wordInput 8)
  let values := data.map BitVec.toNat
  let storage := values.take (size.toNat + spare.toNat)
  let input : IterState :=
    ⟨(if tag == 0 then none else some (storage.take size.toNat)), storage.drop size.toNat⟩
  let separator := (values.drop (size.toNat + spare.toNat)).take sepSize.toNat
  return runs lines input first.toNat second.toNat ++ "#" ++
    runs (fun state => splitSeqConfigured state separator save) input first.toNat second.toNat
