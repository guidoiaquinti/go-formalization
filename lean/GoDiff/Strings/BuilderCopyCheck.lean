import GoLean.Strings.BuilderCopyCheck
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun row => do
  let mode :: capacity :: raw := row | throw (IO.userError "expected identity, capacity and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let id := if mode = 0 then BuilderIdentity.unset else if mode = 1 then .self else .foreign
  let buffer := if capacity < 0 then none else some (bytes, capacity.toNat)
  let receiver := if mode = 3 then none else some (id, ⟨mode != 0, buffer⟩)
  match builderCopyCheck receiver with
  | none => return "nil"
  | some (panicked, identity, b) =>
    let status := if panicked then "copied" else "ok"
    let relation := match identity with
      | .unset => "unset"
      | .self => "self"
      | .foreign => "foreign"
    let cap : Int := (b.buffer.map fun x => Int.ofNat x.2).getD (-1)
    let contents := b.buffer.map Prod.fst |>.getD []
    return s!"{status}:{relation}:{cap}:" ++ String.intercalate "," (contents.map fun byte => toString byte.toNat)
