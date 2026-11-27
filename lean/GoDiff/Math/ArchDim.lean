import GoLean.Math.ArchDim
import GoDiff.Math.Bits.Tuple

open GoLean.Math

def main : IO Unit := runTuple fun fields => do
  let [mode, first, second] := fields | throw (IO.userError "expected profile and two binary64 words")
  let profile := if mode == 0 then ArchDimProfile.amd64 else
    if mode == 1 then ArchDimProfile.arm64 else ArchDimProfile.fallback
  let x ← wordInput 64 first
  let y ← wordInput 64 second
  match archMaxWord profile x y, archMinWord profile x y with
  | .ok maximum, .ok minimum => return s!"{maximum},{minimum}"
  | .error maxMessage, .error minMessage =>
    if maxMessage == minMessage then return "panic:" ++ maxMessage
    else throw (IO.userError "mismatched panic semantics")
  | _, _ => throw (IO.userError "mismatched architecture profiles")
