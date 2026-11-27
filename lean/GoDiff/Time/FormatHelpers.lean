import GoLean.Time.FormatHelpers
import GoDiff.Math.Bits.Tuple

open GoLean.Time

private def encode (r : FormatResult) : String :=
  let status := match r.result with
    | none => "bounds"
    | some (start, value) => s!"{start}:{value}"
  status ++ ":" ++ String.intercalate "," (r.backing.map (toString ·.toNat))

def main : IO Unit := runTuple fun fields => do
  let v :: length :: precision :: raw := fields |
    throw (IO.userError "expected value, visible length, precision and backing")
  let v ← wordInput 64 v
  let backing ← raw.mapM fun b => return UInt8.ofNat (← wordInput 8 b).toNat
  return s!"{encode (fmtInt v.toNat length.toNat backing)}|{encode (fmtFrac v.toNat length.toNat precision backing)}"
