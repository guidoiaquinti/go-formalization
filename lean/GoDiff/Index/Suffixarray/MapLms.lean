import GoLean.Index.Suffixarray.MapLms
import GoDiff.Math.Bits.Tuple

open GoLean.Index.Suffixarray

private def wire (v : Option (Array (BitVec width))) (original : Array (BitVec width)) : String :=
  let status := if v.isSome then "ok" else "bounds"
  status ++ ":" ++ String.intercalate "," ((v.getD original).toList.map (toString ∘ BitVec.toInt))

def main : IO Unit := runTuple fun fields => do
  let length :: count :: values := fields | throw (IO.userError "expected length, LMS count and backing")
  let b32 := (values.map (BitVec.ofInt 32)).toArray
  let b64 := (values.map (BitVec.ofInt 64)).toArray
  return wire (mapLms32 length.toNat b32 count) b32 ++ "|" ++ wire (mapLms64 length.toNat b64 count) b64
