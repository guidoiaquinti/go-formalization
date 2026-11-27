import GoLean.Encoding.Binary.Ensure
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def wire (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun row => do
  let isNil :: length :: count :: raw := row | throw (IO.userError "need nil, length, count and backing")
  let bytes ← raw.mapM (wordInput 8)
  match ensure ⟨bytes, length.toNat, isNil != 0⟩ count with
  | .error err => return s!"panic:{err}"
  | .ok out => return s!"{out.nilBuffer}:{out.allocated}:{wire out.buffer}:{wire out.position}"
